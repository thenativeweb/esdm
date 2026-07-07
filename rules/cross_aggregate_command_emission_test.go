package rules_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCrossAggregateCommandEmission(t *testing.T) {
	rule := findCatalogRule(t, "esdm/modeling/cross-aggregate-command-emission")

	t.Run("does not throw when a command publishes an event in its own aggregate", func(t *testing.T) {
		assert.Empty(t, runRule(t, rule, buildModel(t, minimalParents)))
	})

	t.Run("throws when a command publishes an event owned by a different aggregate", func(t *testing.T) {
		yaml := minimalParents + `---
apiVersion: schema.esdm.io/core/v1
kind: aggregate
name: shipment
scope:
  domain: d
  boundedContext: bc
identifiedBy:
  source: generated
  generator: uuid
state:
  type: object
---
apiVersion: schema.esdm.io/core/v1
kind: command
name: ship-it
scope:
  domain: d
  boundedContext: bc
  aggregate: shipment
data:
  type: object
publishes:
  - agg-done
actors:
  - user
`
		diags := runRule(t, rule, buildModel(t, yaml))
		require.Len(t, diags, 1)
		assert.Contains(t, diags[0].Message, "ship-it")
		assert.Contains(t, diags[0].Message, "agg-done")
		assert.Contains(t, diags[0].Message, "shipment")
		assert.Contains(t, diags[0].Message, "agg")
	})

	t.Run("ignores published names that do not resolve to a known event", func(t *testing.T) {
		// Unresolved references are reported by the resolver,
		// not by this rule, so a typo in `publishes` should
		// produce no diagnostic from us.
		yaml := minimalParents + `---
apiVersion: schema.esdm.io/core/v1
kind: command
name: another
scope:
  domain: d
  boundedContext: bc
  aggregate: agg
data:
  type: object
publishes:
  - does-not-exist
actors:
  - user
`
		assert.Empty(t, runRule(t, rule, buildModel(t, yaml)))
	})

	t.Run("does not throw for DCB-scoped commands publishing aggregate-owned events", func(t *testing.T) {
		// A command scoped to a DCB has no aggregate field; this
		// rule is specifically about aggregate-vs-aggregate
		// emission and stays silent in the DCB case (covered by
		// other rules / cross-context checks if at all).
		yaml := minimalParents + `---
apiVersion: schema.esdm.io/core/v1
kind: dynamic-consistency-boundary
name: dcb
scope:
  domain: d
  boundedContext: bc
identifiedBy:
  - field: id
    source: command
consults:
  - boundedContext: bc
    aggregate: agg
    event: agg-done
    criteria: relevant events
---
apiVersion: schema.esdm.io/core/v1
kind: command
name: dcb-cmd
scope:
  domain: d
  boundedContext: bc
  dynamicConsistencyBoundary: dcb
data:
  type: object
publishes:
  - agg-done
actors:
  - user
`
		assert.Empty(t, runRule(t, rule, buildModel(t, yaml)))
	})
}
