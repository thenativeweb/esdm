package docgen_test

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thenativeweb/esdm/cmd/esdm/commands/docgen"
)

// shopYAML covers one of every kind the tree places, with the
// relations the pages spell out: a command with publisher and
// actor, a free-standing event, a translated term, and a
// context mapping with a term pair.
const shopYAML = `apiVersion: schema.esdm.io/core/v1
kind: domain
name: shop
description: Everything about selling.
---
apiVersion: schema.esdm.io/core/v1
kind: subdomain
name: core-sub
scope:
  domain: shop
type: core
boundedContexts:
  - ordering
---
apiVersion: schema.esdm.io/core/v1
kind: bounded-context
name: ordering
scope:
  domain: shop
description: Taking and fulfilling orders.
language: en
ubiquitousLanguage:
  - term: Order
    definition: A customer's request to purchase.
    avoid:
      - term: Basket
        reason: Reserved for the cart.
    translations:
      - language: de
        term: Bestellung
        definition: Der Kaufwunsch eines Kunden.
        avoid:
          - term: Auftrag
  - term: Customer
    definition: A person who places orders.
---
apiVersion: schema.esdm.io/core/v1
kind: bounded-context
name: billing
scope:
  domain: shop
language: en
ubiquitousLanguage:
  - term: Account
    definition: The ledger of a party we bill.
---
apiVersion: schema.esdm.io/core/v1
kind: aggregate
name: order
scope:
  domain: shop
  boundedContext: ordering
description: An order placed by a customer.
identifiedBy:
  source: generated
  generator: uuid
state:
  type: object
invariants:
  - name: has-lines
    rule: An order has at least one line.
---
apiVersion: schema.esdm.io/core/v1
kind: actor
name: customer
scope:
  domain: shop
  boundedContext: ordering
type: human
---
apiVersion: schema.esdm.io/core/v1
kind: command
name: place
scope:
  domain: shop
  boundedContext: ordering
  aggregate: order
data:
  type: object
publishes:
  - placed
actors:
  - customer
---
apiVersion: schema.esdm.io/core/v1
kind: event
name: placed
scope:
  domain: shop
  boundedContext: ordering
  aggregate: order
data:
  type: object
---
apiVersion: schema.esdm.io/core/v1
kind: read-model
name: orders
scope:
  domain: shop
  boundedContext: ordering
projections:
  - boundedContext: ordering
    aggregate: order
    event: placed
    rule: append
schema:
  type: object
---
apiVersion: schema.esdm.io/core/v1
kind: query
name: list-orders
scope:
  domain: shop
  boundedContext: ordering
readModel: orders
result:
  type: object
actors:
  - customer
---
apiVersion: schema.esdm.io/core/v1
kind: dynamic-consistency-boundary
name: capacity
scope:
  domain: shop
  boundedContext: ordering
identifiedBy:
  - name: id
    source: static
    value: solo
consults:
  - boundedContext: ordering
    aggregate: order
    event: placed
    criteria: relevant
---
apiVersion: schema.esdm.io/core/v1
kind: command
name: reserve
scope:
  domain: shop
  boundedContext: ordering
  dynamicConsistencyBoundary: capacity
data:
  type: object
publishes:
  - reserved
---
apiVersion: schema.esdm.io/core/v1
kind: event
name: reserved
scope:
  domain: shop
  boundedContext: ordering
data:
  type: object
---
apiVersion: schema.esdm.io/core/v1
kind: entity
name: line-item
scope:
  domain: shop
  boundedContext: ordering
schema:
  type: object
  properties:
    sku:
      type: string
identifiedBy:
  source: schema
  field: sku
---
apiVersion: schema.esdm.io/core/v1
kind: value-object
name: money
scope:
  domain: shop
  boundedContext: ordering
schema:
  type: object
---
apiVersion: schema.esdm.io/core/v1
kind: domain-service
name: pricing
scope:
  domain: shop
  boundedContext: ordering
functions:
  - name: compute-total
    arguments:
      type: object
    returns:
      type: object
---
apiVersion: schema.esdm.io/core/v1
kind: process-manager
name: tracker
scope:
  domain: shop
deliveryGuarantee: at-most-once
correlatedBy:
  source: static
  value: solo
state:
  type: object
startsWhen:
  - boundedContext: ordering
    aggregate: order
    event: placed
endsWhen:
  - name: done
    condition: state.completed is true
reactions:
  - when:
      boundedContext: ordering
      aggregate: order
      event: placed
    rule: mark complete
---
apiVersion: schema.esdm.io/core/v1
kind: event-handler
name: notify
scope:
  domain: shop
deliveryGuarantee: at-most-once
handles:
  - boundedContext: ordering
    aggregate: order
    event: placed
sideEffects:
  - type: other
    rule: send mail
---
apiVersion: schema.esdm.io/core/v1
kind: policy
name: react
scope:
  domain: shop
deliveryGuarantee: at-most-once
handles:
  - boundedContext: ordering
    aggregate: order
    event: placed
emits:
  - boundedContext: ordering
    aggregate: order
    command: place
---
apiVersion: schema.esdm.io/core/v1
kind: external-system
name: stripe
scope:
  domain: shop
direction: outbound
---
apiVersion: schema.esdm.io/core/v1
kind: context-mapping
name: ordering-billing
type: customer-supplier
customer:
  domain: shop
  boundedContext: ordering
supplier:
  domain: shop
  boundedContext: billing
terms:
  - customer: Customer
    supplier: Account
---
apiVersion: schema.esdm.io/domain-storytelling/v1
kind: domain-story
name: place-an-order
scope:
  domain: shop
sentences:
  - sequenceNumber: 1
    workObjects:
      - name: order
        annotation: The customer's order.
    edges:
      - from:
          actor: customer
        to:
          workObject: order
---
apiVersion: schema.esdm.io/given-when-then/v1
kind: feature
name: order-placement
scope:
  domain: shop
  boundedContext: ordering
  aggregate: order
scenarios:
  - name: places-an-order
    when:
      command: place
      data: {}
      actor: customer
    then:
      events:
        - event: placed
          data: {}
`

func runDocumentationCommand(t *testing.T, model string, args []string) (string, error) {
	t.Helper()
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "model.esdm.yaml"), []byte(model), 0o644))

	// docgen.Command is a package-level cobra command, so a
	// flag set by one run would otherwise stay in effect for
	// every run after it.
	docgen.Command.Flags().VisitAll(func(flag *pflag.Flag) {
		_ = flag.Value.Set(flag.DefValue)
		flag.Changed = false
	})

	var buf bytes.Buffer
	docgen.Command.SetOut(&buf)
	docgen.Command.SetErr(&buf)
	docgen.Command.SetArgs(append([]string{"--directory", dir}, args...))
	err := docgen.Command.Execute()
	return buf.String(), err
}

func readPage(t *testing.T, out, relative string) string {
	t.Helper()
	content, err := os.ReadFile(filepath.Join(out, filepath.FromSlash(relative)))
	require.NoError(t, err, relative)
	return string(content)
}

func TestDocumentationCommand(t *testing.T) {
	t.Run("writes one page per element at its containment path", func(t *testing.T) {
		out := filepath.Join(t.TempDir(), "docs")

		_, err := runDocumentationCommand(t, shopYAML, []string{"--output", out})
		require.NoError(t, err)

		for _, page := range []string{
			"README.md",
			"domain_shop/README.md",
			"domain_shop/subdomain_core-sub.md",
			"domain_shop/bounded-context_ordering/README.md",
			"domain_shop/bounded-context_ordering/aggregate_order/README.md",
			"domain_shop/bounded-context_ordering/aggregate_order/command_place.md",
			"domain_shop/bounded-context_ordering/aggregate_order/event_placed.md",
			"domain_shop/bounded-context_ordering/aggregate_order/feature_order-placement.md",
			"domain_shop/bounded-context_ordering/dynamic-consistency-boundary_capacity/README.md",
			"domain_shop/bounded-context_ordering/dynamic-consistency-boundary_capacity/command_reserve.md",
			"domain_shop/bounded-context_ordering/event_reserved.md",
			"domain_shop/bounded-context_ordering/read-model_orders/README.md",
			"domain_shop/bounded-context_ordering/query_list-orders.md",
			"domain_shop/bounded-context_ordering/entity_line-item.md",
			"domain_shop/bounded-context_ordering/value-object_money.md",
			"domain_shop/bounded-context_ordering/domain-service_pricing.md",
			"domain_shop/bounded-context_ordering/actor_customer.md",
			"domain_shop/process-manager_tracker/README.md",
			"domain_shop/event-handler_notify.md",
			"domain_shop/policy_react.md",
			"domain_shop/external-system_stripe.md",
			"domain_shop/domain-story_place-an-order.md",
			"context-mapping_ordering-billing.md",
		} {
			assert.FileExists(t, filepath.Join(out, filepath.FromSlash(page)), page)
		}
	})

	t.Run("renders the root page with domains and context mappings", func(t *testing.T) {
		out := filepath.Join(t.TempDir(), "docs")

		_, err := runDocumentationCommand(t, shopYAML, []string{"--output", out})
		require.NoError(t, err)

		root := readPage(t, out, "README.md")
		assert.Contains(t, root, "# Model\n")
		assert.Contains(t, root, "## Domains\n\n- [shop](domain_shop/README.md)")
		assert.Contains(t, root, "## Context Mappings\n\n- [ordering-billing](context-mapping_ordering-billing.md) (customer-supplier)")
	})

	t.Run("renders an aggregate page with header, description, details, and linked children", func(t *testing.T) {
		out := filepath.Join(t.TempDir(), "docs")

		_, err := runDocumentationCommand(t, shopYAML, []string{"--output", out})
		require.NoError(t, err)

		page := readPage(t, out, "domain_shop/bounded-context_ordering/aggregate_order/README.md")
		assert.Contains(t, page, "# order\n\nAggregate `esdm:domain=shop/bounded-context=ordering/aggregate=order`\n\n1 cmd · 1 evt · 1 inv · 1 feat\n\nAn order placed by a customer.\n")
		assert.Contains(t, page, "## Details\n\n- identifiedBy: generated/uuid\n- invariant \"has-lines\": An order has at least one line.\n")
		assert.Contains(t, page, "## Commands\n\n- [place](command_place.md) – publishes [placed](event_placed.md), issued by [customer](../actor_customer.md)\n")
		assert.Contains(t, page, "## Events\n\n- [placed](event_placed.md) – published by [place](command_place.md)\n")
		assert.Contains(t, page, "## Features\n\n- [order-placement](feature_order-placement.md) · 1 scenario\n")
	})

	t.Run("renders a leaf page with its own header and relations", func(t *testing.T) {
		out := filepath.Join(t.TempDir(), "docs")

		_, err := runDocumentationCommand(t, shopYAML, []string{"--output", out})
		require.NoError(t, err)

		page := readPage(t, out, "domain_shop/bounded-context_ordering/aggregate_order/command_place.md")
		assert.Contains(t, page, "# place\n\nCommand `esdm:domain=shop/bounded-context=ordering/aggregate=order/command=place`\n")
		assert.Contains(t, page, "Publishes [placed](event_placed.md), issued by [customer](../actor_customer.md).\n")

		event := readPage(t, out, "domain_shop/bounded-context_ordering/event_reserved.md")
		assert.Contains(t, event, "# reserved\n\nEvent `esdm:domain=shop/bounded-context=ordering/event=reserved`\n")
		assert.Contains(t, event, "Published by [reserve](dynamic-consistency-boundary_capacity/command_reserve.md).\n")
	})

	t.Run("renders the ubiquitous language with translations indented under the term", func(t *testing.T) {
		out := filepath.Join(t.TempDir(), "docs")

		_, err := runDocumentationCommand(t, shopYAML, []string{"--output", out})
		require.NoError(t, err)

		page := readPage(t, out, "domain_shop/bounded-context_ordering/README.md")
		assert.Contains(t, page, "## Ubiquitous Language\n\nWritten in `en`.\n\n- **Customer** – A person who places orders.\n- **Order** – A customer's request to purchase.\n  - _Avoid the term \"Basket\"._ Reserved for the cart.\n  - `de` **Bestellung** – Der Kaufwunsch eines Kunden. _Avoid the term \"Auftrag\"._\n")
	})

	t.Run("writes a read model's projections out as a linked relation", func(t *testing.T) {
		out := filepath.Join(t.TempDir(), "docs")

		_, err := runDocumentationCommand(t, shopYAML, []string{"--output", out})
		require.NoError(t, err)

		page := readPage(t, out, "domain_shop/bounded-context_ordering/README.md")
		assert.Contains(t, page, "## Read Models\n\n- [orders](read-model_orders/README.md) – projects [placed](aggregate_order/event_placed.md)\n")

		readModel := readPage(t, out, "domain_shop/bounded-context_ordering/read-model_orders/README.md")
		assert.Contains(t, readModel, "# orders\n\nRead Model `esdm:domain=shop/bounded-context=ordering/read-model=orders`\n\nProjects [placed](../aggregate_order/event_placed.md).\n")
	})

	t.Run("renders a context mapping with its endpoints and term pairs", func(t *testing.T) {
		out := filepath.Join(t.TempDir(), "docs")

		_, err := runDocumentationCommand(t, shopYAML, []string{"--output", out})
		require.NoError(t, err)

		page := readPage(t, out, "context-mapping_ordering-billing.md")
		assert.Contains(t, page, "# ordering-billing\n\nContext Mapping `esdm:context-mapping=ordering-billing`\n\ncustomer-supplier: customer [ordering](domain_shop/bounded-context_ordering/README.md), supplier [billing](domain_shop/bounded-context_billing/README.md)\n")
		assert.Contains(t, page, "## Terms\n\n- Customer in [ordering](domain_shop/bounded-context_ordering/README.md) corresponds to Account in [billing](domain_shop/bounded-context_billing/README.md)\n")
	})

	t.Run("narrows to a subtree, keeps full paths, and falls back to references for links outside it", func(t *testing.T) {
		out := filepath.Join(t.TempDir(), "docs")

		_, err := runDocumentationCommand(t, shopYAML, []string{"--output", out, "shop/ordering/order"})
		require.NoError(t, err)

		assert.FileExists(t, filepath.Join(out, "domain_shop", "bounded-context_ordering", "aggregate_order", "README.md"))
		assert.NoFileExists(t, filepath.Join(out, "domain_shop", "README.md"))
		assert.NoFileExists(t, filepath.Join(out, "README.md"))

		page := readPage(t, out, "domain_shop/bounded-context_ordering/aggregate_order/README.md")
		assert.Contains(t, page, "- [place](command_place.md) – publishes [placed](event_placed.md), issued by customer (`esdm:domain=shop/bounded-context=ordering/actor=customer`)\n")
	})

	t.Run("refuses a non-empty output directory unless --force is given", func(t *testing.T) {
		out := t.TempDir()
		require.NoError(t, os.WriteFile(filepath.Join(out, "stale.md"), []byte("old"), 0o644))

		_, err := runDocumentationCommand(t, shopYAML, []string{"--output", out})
		require.Error(t, err)
		// %q mirrors the command's own formatting, which escapes
		// the backslashes of a Windows path.
		assert.Equal(t, fmt.Sprintf("output directory %q is not empty; use --force to clear it first", out), err.Error())

		_, err = runDocumentationCommand(t, shopYAML, []string{"--output", out, "--force"})
		require.NoError(t, err)
		assert.NoFileExists(t, filepath.Join(out, "stale.md"))
		assert.FileExists(t, filepath.Join(out, "README.md"))
	})

	t.Run("requires --output", func(t *testing.T) {
		_, err := runDocumentationCommand(t, shopYAML, nil)
		require.Error(t, err)
	})

	t.Run("returns an error for an unknown path segment", func(t *testing.T) {
		out := filepath.Join(t.TempDir(), "docs")

		_, err := runDocumentationCommand(t, shopYAML, []string{"--output", out, "shop/nowhere"})
		require.Error(t, err)
		assert.Equal(t, `no element "nowhere" under "shop"`, err.Error())
	})
}
