package rules_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBoundedContextWithoutContent(t *testing.T) {
	rule := findCatalogRule(t, "esdm/modeling/bounded-context-without-content")

	t.Run("does not throw when the bounded context hosts an aggregate", func(t *testing.T) {
		assert.Empty(t, runRule(t, rule, buildModel(t, minimalParents)))
	})

	t.Run("does not throw when the bounded context hosts a dynamic-consistency-boundary", func(t *testing.T) {
		yaml := minimalParents + `---
apiVersion: schema.esdm.io/core/v1
kind: bounded-context
name: capacity-planning
scope:
  domain: d
---
apiVersion: schema.esdm.io/core/v1
kind: dynamic-consistency-boundary
name: capacity
scope:
  domain: d
  boundedContext: capacity-planning
identifiedBy:
  - name: id
    source: static
    value: solo
consults:
  - boundedContext: bc
    aggregate: agg
    event: agg-done
    criteria: relevant
`
		assert.Empty(t, runRule(t, rule, buildModel(t, yaml)))
	})

	t.Run("does not throw when the bounded context hosts only read models", func(t *testing.T) {
		// A projection-only context derives every view from events owned
		// elsewhere. It holds no consistency unit by design, so the rule
		// must accept it rather than push a placeholder aggregate into it.
		yaml := minimalParents + `---
apiVersion: schema.esdm.io/core/v1
kind: bounded-context
name: reporting
scope:
  domain: d
---
apiVersion: schema.esdm.io/core/v1
kind: read-model
name: activity-overview
scope:
  domain: d
  boundedContext: reporting
projections:
  - boundedContext: bc
    aggregate: agg
    event: agg-done
    rule: count completed aggregates
schema:
  type: object
`
		assert.Empty(t, runRule(t, rule, buildModel(t, yaml)))
	})

	t.Run("throws when a bounded context hosts no aggregate, no DCB and no read model", func(t *testing.T) {
		yaml := `apiVersion: schema.esdm.io/core/v1
kind: domain
name: d
---
apiVersion: schema.esdm.io/core/v1
kind: bounded-context
name: empty-bc
scope:
  domain: d
`
		diags := runRule(t, rule, buildModel(t, yaml))
		require.Len(t, diags, 1)
		assert.Contains(t, diags[0].Message, "empty-bc")
	})

	t.Run("throws when a bounded context holds only supporting elements", func(t *testing.T) {
		// Entities, value objects and actors only carry meaning around a
		// consistency unit or a read model; on their own they leave the
		// context a placeholder.
		yaml := minimalParents + `---
apiVersion: schema.esdm.io/core/v1
kind: bounded-context
name: shared-terms
scope:
  domain: d
---
apiVersion: schema.esdm.io/core/v1
kind: value-object
name: money
scope:
  domain: d
  boundedContext: shared-terms
schema:
  type: object
---
apiVersion: schema.esdm.io/core/v1
kind: actor
name: clerk
scope:
  domain: d
  boundedContext: shared-terms
type: human
`
		diags := runRule(t, rule, buildModel(t, yaml))
		require.Len(t, diags, 1)
		assert.Contains(t, diags[0].Message, "shared-terms")
	})
}
