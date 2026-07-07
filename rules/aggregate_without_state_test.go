package rules_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAggregateWithoutState(t *testing.T) {
	rule := findCatalogRule(t, "esdm/modeling/aggregate-without-state")

	t.Run("does not throw when the aggregate declares state", func(t *testing.T) {
		assert.Empty(t, runRule(t, rule, buildModel(t, minimalParents)))
	})

	t.Run("throws when the aggregate has no state field", func(t *testing.T) {
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
  source: generated
  generator: uuid
`
		diags := runRule(t, rule, buildModelTolerantOfSchemaErrors(t, yaml))
		require.Len(t, diags, 1)
		assert.Contains(t, diags[0].Message, "agg")
	})
}
