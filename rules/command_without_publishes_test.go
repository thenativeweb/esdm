package rules_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCommandWithoutPublishes(t *testing.T) {
	rule := findCatalogRule(t, "esdm/modeling/command-without-publishes")

	t.Run("does not throw when the command publishes at least one event", func(t *testing.T) {
		assert.Empty(t, runRule(t, rule, buildModel(t, minimalParents)))
	})

	t.Run("throws when the command has no publishes field", func(t *testing.T) {
		yaml := minimalParents + `---
apiVersion: schema.esdm.io/core/v1
kind: command
name: silent
scope:
  domain: d
  boundedContext: bc
  aggregate: agg
data:
  type: object
actors:
  - user
`
		diags := runRule(t, rule, buildModelTolerantOfSchemaErrors(t, yaml))
		require.Len(t, diags, 1)
		assert.Contains(t, diags[0].Message, "silent")
	})
}
