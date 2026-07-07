package rules_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAggregateIdentifiedByField(t *testing.T) {
	rule := findCatalogRule(t, "esdm/structure/aggregate-identified-by-field")

	t.Run("does not throw when identifiedBy.source is not 'state'", func(t *testing.T) {
		// minimalParents uses source: generated, generator: uuid -
		// no field reference into state, so the rule has nothing
		// to check.
		assert.Empty(t, runRule(t, rule, buildModel(t, minimalParents)))
	})

	t.Run("does not throw when the referenced field exists in state.properties", func(t *testing.T) {
		yaml := `apiVersion: schema.esdm.io/core/v1
kind: domain
name: d
---
apiVersion: schema.esdm.io/core/v1
kind: bounded-context
name: bc
scope:
  domain: d
---
apiVersion: schema.esdm.io/core/v1
kind: aggregate
name: agg
scope:
  domain: d
  boundedContext: bc
identifiedBy:
  source: state
  field: order-id
state:
  type: object
  properties:
    order-id:
      type: string
`
		assert.Empty(t, runRule(t, rule, buildModel(t, yaml)))
	})

	t.Run("throws when the referenced field is missing from state.properties", func(t *testing.T) {
		yaml := `apiVersion: schema.esdm.io/core/v1
kind: domain
name: d
---
apiVersion: schema.esdm.io/core/v1
kind: bounded-context
name: bc
scope:
  domain: d
---
apiVersion: schema.esdm.io/core/v1
kind: aggregate
name: agg
scope:
  domain: d
  boundedContext: bc
identifiedBy:
  source: state
  field: order-id
state:
  type: object
  properties:
    something-else:
      type: string
`
		diags := runRule(t, rule, buildModel(t, yaml))
		require.Len(t, diags, 1)
		assert.Contains(t, diags[0].Message, "agg")
		assert.Contains(t, diags[0].Message, "order-id")
	})

	t.Run("throws when state has no properties block at all", func(t *testing.T) {
		yaml := `apiVersion: schema.esdm.io/core/v1
kind: domain
name: d
---
apiVersion: schema.esdm.io/core/v1
kind: bounded-context
name: bc
scope:
  domain: d
---
apiVersion: schema.esdm.io/core/v1
kind: aggregate
name: agg
scope:
  domain: d
  boundedContext: bc
identifiedBy:
  source: state
  field: order-id
state:
  type: object
`
		diags := runRule(t, rule, buildModel(t, yaml))
		assert.Len(t, diags, 1)
	})
}
