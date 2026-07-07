package rules_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEventNameWithAggregatePrefix(t *testing.T) {
	rule := findCatalogRule(t, "esdm/modeling/event-name-with-aggregate-prefix")

	t.Run("throws when the event name starts with its aggregate's name", func(t *testing.T) {
		// minimalParents uses aggregate "agg" and event "agg-done" -
		// the aggregate name is redundantly repeated at the start
		// of the event name.
		diags := runRule(t, rule, buildModel(t, minimalParents))
		require.Len(t, diags, 1)
		assert.Contains(t, diags[0].Message, "agg-done")
		assert.Contains(t, diags[0].Message, "agg")
	})

	t.Run("does not throw when the event name does not start with the aggregate's name", func(t *testing.T) {
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
name: order
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
kind: event
name: placed
scope:
  domain: d
  boundedContext: bc
  aggregate: order
data:
  type: object
`
		assert.Empty(t, runRule(t, rule, buildModel(t, yaml)))
	})

	t.Run("does not throw for BC-scoped events that do not belong to an aggregate", func(t *testing.T) {
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
kind: event
name: bc-wide-signal
scope:
  domain: d
  boundedContext: bc
data:
  type: object
`
		assert.Empty(t, runRule(t, rule, buildModel(t, yaml)))
	})
}
