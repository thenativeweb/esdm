package rules_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAggregateWithoutEvents(t *testing.T) {
	rule := findCatalogRule(t, "esdm/modeling/aggregate-without-events")

	t.Run("does not throw when an aggregate has a matching event", func(t *testing.T) {
		m := buildModel(t, minimalParents)
		assert.Empty(t, runRule(t, rule, m))
	})

	t.Run("throws when an aggregate has no matching event", func(t *testing.T) {
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
name: lonely
scope:
  domain: d
  boundedContext: bc
identifiedBy:
  source: generated
  generator: uuid
state:
  type: object
`
		diags := runRule(t, rule, buildModel(t, yaml))
		require.Len(t, diags, 1)
		assert.Contains(t, diags[0].Message, "lonely")
	})
}
