package rules_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBoundedContextWithoutConsistencyUnit(t *testing.T) {
	rule := findCatalogRule(t, "esdm/modeling/bounded-context-without-consistency-unit")

	t.Run("does not throw when the bounded context hosts an aggregate", func(t *testing.T) {
		assert.Empty(t, runRule(t, rule, buildModel(t, minimalParents)))
	})

	t.Run("throws when a bounded context hosts no aggregate and no DCB", func(t *testing.T) {
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
}
