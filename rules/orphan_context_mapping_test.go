package rules_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOrphanContextMapping(t *testing.T) {
	rule := findCatalogRule(t, "esdm/modeling/orphan-context-mapping")

	t.Run("does not throw when the mapping is used by a cross-BC read-model projection", func(t *testing.T) {
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
		diags := runRule(t, rule, buildModel(t, yaml))
		assert.Empty(t, diags)
	})

	t.Run("throws when the mapping connects two BCs that no consumer ever spans", func(t *testing.T) {
		yaml := twoBCParents + `---
apiVersion: schema.esdm.io/core/v1
kind: context-mapping
name: dormant
type: customer-supplier
customer:
  domain: d
  boundedContext: shipping
supplier:
  domain: d
  boundedContext: ordering
`
		diags := runRule(t, rule, buildModel(t, yaml))
		require.Len(t, diags, 1)
		assert.Contains(t, diags[0].Message, "dormant")
		assert.Contains(t, diags[0].Message, "ordering")
		assert.Contains(t, diags[0].Message, "shipping")
	})

	t.Run("does not throw for mappings whose endpoints are external systems only", func(t *testing.T) {
		yaml := twoBCParents + `---
apiVersion: schema.esdm.io/core/v1
kind: external-system
name: payments
scope:
  domain: d
type: third-party
---
apiVersion: schema.esdm.io/core/v1
kind: external-system
name: emailer
scope:
  domain: d
type: third-party
---
apiVersion: schema.esdm.io/core/v1
kind: context-mapping
name: payments-to-emailer
type: anti-corruption-layer
downstream:
  domain: d
  externalSystem: payments
upstream:
  domain: d
  externalSystem: emailer
`
		// The mapping doesn't connect any BC pair, so the
		// orphan-mapping rule has nothing to say about it
		// (the orphan-external-system rule covers
		// the externals' usefulness separately).
		diags := runRule(t, rule, buildModel(t, yaml))
		assert.Empty(t, diags)
	})
}
