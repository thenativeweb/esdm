package rules_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const twoBCParents = `apiVersion: schema.esdm.io/core/v1
kind: domain
name: d
---
apiVersion: schema.esdm.io/core/v1
kind: bounded-context
name: ordering
scope:
  domain: d
---
apiVersion: schema.esdm.io/core/v1
kind: bounded-context
name: shipping
scope:
  domain: d
---
apiVersion: schema.esdm.io/core/v1
kind: aggregate
name: order
scope:
  domain: d
  boundedContext: ordering
identifiedBy:
  source: generated
  generator: uuid
state:
  type: object
---
apiVersion: schema.esdm.io/core/v1
kind: aggregate
name: shipment
scope:
  domain: d
  boundedContext: shipping
identifiedBy:
  source: generated
  generator: uuid
state:
  type: object
---
apiVersion: schema.esdm.io/core/v1
kind: actor
name: user
scope:
  domain: d
  boundedContext: ordering
type: human
---
apiVersion: schema.esdm.io/core/v1
kind: command
name: place-order
scope:
  domain: d
  boundedContext: ordering
  aggregate: order
data:
  type: object
publishes:
  - placed
actors:
  - user
---
apiVersion: schema.esdm.io/core/v1
kind: event
name: placed
scope:
  domain: d
  boundedContext: ordering
  aggregate: order
data:
  type: object
`

func TestCrossBoundedContextReferenceWithoutMapping(t *testing.T) {
	rule := findCatalogRule(t, "esdm/modeling/cross-bounded-context-reference-without-mapping")

	t.Run("does not throw when an in-BC read-model projection stays in its own BC", func(t *testing.T) {
		yaml := twoBCParents + `---
apiVersion: schema.esdm.io/core/v1
kind: read-model
name: order-list
scope:
  domain: d
  boundedContext: ordering
projections:
  - boundedContext: ordering
    aggregate: order
    event: placed
    rule: tally orders
schema:
  type: object
`
		assert.Empty(t, runRule(t, rule, buildModel(t, yaml)))
	})

	t.Run("throws when a read-model projects an event from a foreign BC and no mapping exists", func(t *testing.T) {
		yaml := twoBCParents + `---
apiVersion: schema.esdm.io/core/v1
kind: read-model
name: shipping-overview
scope:
  domain: d
  boundedContext: shipping
projections:
  - boundedContext: ordering
    aggregate: order
    event: placed
    rule: track upstream orders
schema:
  type: object
`
		diags := runRule(t, rule, buildModel(t, yaml))
		require.Len(t, diags, 1)
		assert.Contains(t, diags[0].Message, "shipping-overview")
		assert.Contains(t, diags[0].Message, "ordering")
		assert.Contains(t, diags[0].Message, "shipping")
	})

	t.Run("does not throw when a context mapping connects the two BCs", func(t *testing.T) {
		yaml := twoBCParents + `---
apiVersion: schema.esdm.io/core/v1
kind: context-mapping
name: order-to-shipping
type: customer-supplier
customer:
  domain: d
  boundedContext: shipping
supplier:
  domain: d
  boundedContext: ordering
---
apiVersion: schema.esdm.io/core/v1
kind: read-model
name: shipping-overview
scope:
  domain: d
  boundedContext: shipping
projections:
  - boundedContext: ordering
    aggregate: order
    event: placed
    rule: track upstream orders
schema:
  type: object
`
		assert.Empty(t, runRule(t, rule, buildModel(t, yaml)))
	})

	t.Run("throws on a domain-scoped policy that emits into a foreign BC without a mapping", func(t *testing.T) {
		// Policy reacts to ordering.placed and emits into shipping -
		// two BCs touched, no mapping declared.
		yaml := twoBCParents + `---
apiVersion: schema.esdm.io/core/v1
kind: command
name: ship-it
scope:
  domain: d
  boundedContext: shipping
  aggregate: shipment
data:
  type: object
publishes:
  - shipped
actors:
  - user
---
apiVersion: schema.esdm.io/core/v1
kind: event
name: shipped
scope:
  domain: d
  boundedContext: shipping
  aggregate: shipment
data:
  type: object
---
apiVersion: schema.esdm.io/core/v1
kind: policy
name: ship-on-order
scope:
  domain: d
deliveryGuarantee: at-least-once
handles:
  - boundedContext: ordering
    aggregate: order
    event: placed
emits:
  - boundedContext: shipping
    aggregate: shipment
    command: ship-it
`
		diags := runRule(t, rule, buildModel(t, yaml))
		require.Len(t, diags, 1)
		assert.Contains(t, diags[0].Message, "ship-on-order")
	})

	t.Run("does not throw on a domain-scoped policy whose references all stay in one BC", func(t *testing.T) {
		yaml := twoBCParents + `---
apiVersion: schema.esdm.io/core/v1
kind: command
name: re-place
scope:
  domain: d
  boundedContext: ordering
  aggregate: order
data:
  type: object
publishes:
  - placed
actors:
  - user
---
apiVersion: schema.esdm.io/core/v1
kind: policy
name: in-context
scope:
  domain: d
deliveryGuarantee: at-least-once
handles:
  - boundedContext: ordering
    aggregate: order
    event: placed
emits:
  - boundedContext: ordering
    aggregate: order
    command: re-place
`
		assert.Empty(t, runRule(t, rule, buildModel(t, yaml)))
	})
}
