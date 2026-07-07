package view_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thenativeweb/esdm/cmd/esdm/commands/view"
)

// extendedCatalogYAML covers every catalog kind in a
// single model so that each `build*` branch in the
// builder is exercised end-to-end. Diagnostics may or
// may not be present - the only requirement is that
// loader and renderer accept the documents.
const extendedCatalogYAML = `apiVersion: schema.esdm.io/core/v1
kind: domain
name: shop
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
---
apiVersion: schema.esdm.io/core/v1
kind: aggregate
name: order
scope:
  domain: shop
  boundedContext: ordering
identifiedBy:
  source: generated
  generator: uuid
state:
  type: object
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
  source: event-field
  field: correlation-id
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
name: shop-to-self
type: customer-supplier
customer:
  domain: shop
  boundedContext: ordering
supplier:
  domain: shop
  boundedContext: ordering
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
description: Place an order on the customer's behalf.
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

// minimalDomainYAML is a single-document model with a
// domain plus a bounded context plus an aggregate plus a
// command/event/actor - the smallest model that
// produces no diagnostics. Tests use it as a baseline.
const minimalDomainYAML = `apiVersion: schema.esdm.io/core/v1
kind: domain
name: shop
---
apiVersion: schema.esdm.io/core/v1
kind: bounded-context
name: ordering
scope:
  domain: shop
---
apiVersion: schema.esdm.io/core/v1
kind: aggregate
name: order
scope:
  domain: shop
  boundedContext: ordering
identifiedBy:
  source: generated
  generator: uuid
state:
  type: object
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
`

// contextMappingEndpointsYAML covers each role-named
// mapping branch (customer/supplier, conformist/upstream,
// downstream/upstream, host/consumer, publisher/consumer)
// plus a participants-based branch (shared-kernel) so
// every endpoint accessor in contextMappingTouchesDomain
// is exercised. All endpoints stay inside the same
// domain, keeping the model topologically simple.
const contextMappingEndpointsYAML = `apiVersion: schema.esdm.io/core/v1
kind: domain
name: shop
---
apiVersion: schema.esdm.io/core/v1
kind: bounded-context
name: a
scope:
  domain: shop
---
apiVersion: schema.esdm.io/core/v1
kind: bounded-context
name: b
scope:
  domain: shop
---
apiVersion: schema.esdm.io/core/v1
kind: context-mapping
name: cm-customer-supplier
type: customer-supplier
customer:
  domain: shop
  boundedContext: a
supplier:
  domain: shop
  boundedContext: b
---
apiVersion: schema.esdm.io/core/v1
kind: context-mapping
name: cm-conformist
type: conformist
conformist:
  domain: shop
  boundedContext: a
upstream:
  domain: shop
  boundedContext: b
---
apiVersion: schema.esdm.io/core/v1
kind: context-mapping
name: cm-acl
type: anti-corruption-layer
downstream:
  domain: shop
  boundedContext: a
upstream:
  domain: shop
  boundedContext: b
---
apiVersion: schema.esdm.io/core/v1
kind: context-mapping
name: cm-open-host
type: open-host-service
host:
  domain: shop
  boundedContext: a
consumer:
  domain: shop
  boundedContext: b
---
apiVersion: schema.esdm.io/core/v1
kind: context-mapping
name: cm-published-language
type: published-language
publisher:
  domain: shop
  boundedContext: a
consumer:
  domain: shop
  boundedContext: b
---
apiVersion: schema.esdm.io/core/v1
kind: context-mapping
name: cm-shared-kernel
type: shared-kernel
participants:
  - domain: shop
    boundedContext: a
  - domain: shop
    boundedContext: b
`

func writeFile(t *testing.T, dir, name, content string) {
	t.Helper()
	require.NoError(t, os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644))
}

func writeMinimalModel(t *testing.T, dir string) {
	t.Helper()
	writeFile(t, dir, "model.esdm.yaml", minimalDomainYAML)
}

func runViewCommand(t *testing.T, args []string) (string, error) {
	t.Helper()
	var buf bytes.Buffer
	view.Command.SetOut(&buf)
	view.Command.SetErr(&buf)
	view.Command.SetArgs(args)
	err := view.Command.Execute()
	return buf.String(), err
}

func TestViewCommand(t *testing.T) {
	t.Run("renders the domain skeleton without arguments", func(t *testing.T) {
		dir := t.TempDir()
		writeMinimalModel(t, dir)

		out, err := runViewCommand(t, []string{"--directory", dir, "--color", "never"})
		require.NoError(t, err)
		assert.Contains(t, out, "domain shop")
		assert.Contains(t, out, "bounded-context ordering")
		assert.Contains(t, out, "aggregate order")
	})

	t.Run("narrows to a sub-tree when given a path", func(t *testing.T) {
		dir := t.TempDir()
		writeMinimalModel(t, dir)

		out, err := runViewCommand(t, []string{"--directory", dir, "--color", "never", "shop/ordering/order"})
		require.NoError(t, err)
		assert.Contains(t, out, "aggregate order")
		assert.Contains(t, out, "command place")
		assert.Contains(t, out, "event placed")
		assert.NotContains(t, out, "domain shop")
	})

	t.Run("returns an error for an unknown path segment", func(t *testing.T) {
		dir := t.TempDir()
		writeMinimalModel(t, dir)

		_, err := runViewCommand(t, []string{"--directory", dir, "--color", "never", "shop/nonexistent"})
		assert.Error(t, err)
	})

	t.Run("emits compact stats by default and detail lines with --with-details", func(t *testing.T) {
		dir := t.TempDir()
		writeMinimalModel(t, dir)

		compact, err := runViewCommand(t, []string{"--directory", dir, "--color", "never", "shop/ordering/order"})
		require.NoError(t, err)
		// Compact does not include the identifiedBy line.
		assert.NotContains(t, compact, "identifiedBy")

		detailed, err := runViewCommand(t, []string{"--directory", dir, "--color", "never", "--with-details", "shop/ordering/order"})
		require.NoError(t, err)
		assert.Contains(t, detailed, "identifiedBy: generated/uuid")
	})

	t.Run("renders every catalog kind in an extended model", func(t *testing.T) {
		dir := t.TempDir()
		writeFile(t, dir, "model.esdm.yaml", extendedCatalogYAML)

		out, err := runViewCommand(t, []string{"--directory", dir, "--color", "never"})
		require.NoError(t, err)

		assert.Contains(t, out, "domain shop")
		assert.Contains(t, out, "subdomain core-sub")
		assert.Contains(t, out, "bounded-context ordering")
		assert.Contains(t, out, "aggregate order")
		assert.Contains(t, out, "command place")
		assert.Contains(t, out, "event placed")
		assert.Contains(t, out, "read-model orders")
		assert.Contains(t, out, "query list-orders")
		assert.Contains(t, out, "actor customer")
		assert.Contains(t, out, "dynamic-consistency-boundary capacity")
		assert.Contains(t, out, "value-object money")
		assert.Contains(t, out, "domain-service pricing")
		assert.Contains(t, out, "process-manager tracker")
		assert.Contains(t, out, "event-handler notify")
		assert.Contains(t, out, "policy react")
		assert.Contains(t, out, "external-system stripe")
		assert.Contains(t, out, "context-mapping shop-to-self")
		assert.Contains(t, out, "domain-story place-an-order")
		assert.Contains(t, out, "feature order-placement")
	})

	t.Run("annotates the domain with per-kind stats covering every direct child", func(t *testing.T) {
		dir := t.TempDir()
		writeFile(t, dir, "model.esdm.yaml", extendedCatalogYAML)

		out, err := runViewCommand(t, []string{"--directory", dir, "--color", "never"})
		require.NoError(t, err)
		// The domain may carry an annotator marker (an
		// error or warning glyph) between name and stats,
		// so assert on the stats fragment alone.
		assert.Contains(t, out, "1 sub · 1 bc · 1 pm · 1 eh · 1 pol · 1 es · 1 cm · 1 story · 1 feat")
	})

	t.Run("omits zero counts from the domain stats line", func(t *testing.T) {
		dir := t.TempDir()
		writeMinimalModel(t, dir)

		out, err := runViewCommand(t, []string{"--directory", dir, "--color", "never"})
		require.NoError(t, err)
		assert.Contains(t, out, "domain shop  1 bc\n")
		assert.NotContains(t, out, "0 sub")
		assert.NotContains(t, out, "0 pm")
	})

	t.Run("renders every context-mapping endpoint shape", func(t *testing.T) {
		dir := t.TempDir()
		writeFile(t, dir, "model.esdm.yaml", contextMappingEndpointsYAML)

		out, err := runViewCommand(t, []string{"--directory", dir, "--color", "never"})
		require.NoError(t, err)

		assert.Contains(t, out, "context-mapping cm-customer-supplier")
		assert.Contains(t, out, "context-mapping cm-conformist")
		assert.Contains(t, out, "context-mapping cm-acl")
		assert.Contains(t, out, "context-mapping cm-open-host")
		assert.Contains(t, out, "context-mapping cm-published-language")
		assert.Contains(t, out, "context-mapping cm-shared-kernel")
	})

	t.Run("annotates events with the commands that publish them", func(t *testing.T) {
		dir := t.TempDir()
		writeMinimalModel(t, dir)

		out, err := runViewCommand(t, []string{"--directory", dir, "--color", "never"})
		require.NoError(t, err)
		assert.Contains(t, out, "event placed  ← place")
	})

	t.Run("lists every command on the event when several publish it", func(t *testing.T) {
		dir := t.TempDir()
		writeMinimalModel(t, dir)
		// A second command in the same aggregate also
		// publishes "placed" - the event's stats line must
		// surface both publishers, sorted alphabetically.
		writeFile(t, dir, "extra-command.esdm.yaml", `apiVersion: schema.esdm.io/core/v1
kind: command
name: replace
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
`)

		out, err := runViewCommand(t, []string{"--directory", dir, "--color", "never"})
		require.NoError(t, err)
		assert.Contains(t, out, "event placed  ← place, replace")
	})

	t.Run("omits the publisher stats line when no command publishes the event", func(t *testing.T) {
		dir := t.TempDir()
		writeMinimalModel(t, dir)
		// Add an orphan event in the same aggregate that
		// no command publishes - its row must not carry a
		// left-arrow publisher suffix.
		writeFile(t, dir, "orphan-event.esdm.yaml", `apiVersion: schema.esdm.io/core/v1
kind: event
name: orphaned
scope:
  domain: shop
  boundedContext: ordering
  aggregate: order
data:
  type: object
`)

		out, err := runViewCommand(t, []string{"--directory", dir, "--color", "never"})
		require.NoError(t, err)
		// The orphan event may carry a diagnostic marker
		// after its name, but its row must not include
		// the left-arrow publisher suffix.
		assert.NotRegexp(t, `event orphaned[^\n]*←`, out)
	})

	t.Run("inline-marks nodes that are subject to a linter diagnostic", func(t *testing.T) {
		dir := t.TempDir()
		writeMinimalModel(t, dir)
		// Add an event whose name reuses its aggregate's
		// prefix - that throws the
		// event-name-with-aggregate-prefix warning.
		writeFile(t, dir, "extra-event.esdm.yaml", `apiVersion: schema.esdm.io/core/v1
kind: event
name: order-shipped
scope:
  domain: shop
  boundedContext: ordering
  aggregate: order
data:
  type: object
`)

		out, err := runViewCommand(t, []string{"--directory", dir, "--color", "never"})
		require.NoError(t, err)
		assert.Contains(t, out, "⚠")
	})
}
