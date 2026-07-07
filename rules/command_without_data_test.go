package rules_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCommandWithoutData(t *testing.T) {
	rule := findCatalogRule(t, "esdm/modeling/command-without-data")

	t.Run("does not throw when the command declares data", func(t *testing.T) {
		assert.Empty(t, runRule(t, rule, buildModel(t, minimalParents)))
	})

	t.Run("throws when the command has no data field", func(t *testing.T) {
		yaml := minimalParents + `---
apiVersion: schema.esdm.io/core/v1
kind: command
name: dataless
scope:
  domain: d
  boundedContext: bc
  aggregate: agg
publishes:
  - agg-done
actors:
  - user
`
		diags := runRule(t, rule, buildModelTolerantOfSchemaErrors(t, yaml))
		require.Len(t, diags, 1)
		assert.Contains(t, diags[0].Message, "dataless")
	})
}
