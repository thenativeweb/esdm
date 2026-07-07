package rules_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSubdomainWithoutBoundedContext(t *testing.T) {
	rule := findCatalogRule(t, "esdm/modeling/subdomain-without-bounded-context")

	t.Run("does not throw when the subdomain lists at least one bounded context", func(t *testing.T) {
		yaml := minimalParents + `---
apiVersion: schema.esdm.io/core/v1
kind: subdomain
name: sales
scope:
  domain: d
type: core
boundedContexts:
  - bc
`
		assert.Empty(t, runRule(t, rule, buildModel(t, yaml)))
	})

	t.Run("throws when the subdomain has no boundedContexts field at all", func(t *testing.T) {
		// The JSON Schema rejects this today (minItems: 1), but
		// this rule is the safety net for an accidental schema
		// relaxation - we deliberately build the model
		// tolerantly so the rule can be exercised in isolation.
		yaml := minimalParents + `---
apiVersion: schema.esdm.io/core/v1
kind: subdomain
name: empty-sub
scope:
  domain: d
type: core
`
		diags := runRule(t, rule, buildModelTolerantOfSchemaErrors(t, yaml))
		require.Len(t, diags, 1)
		assert.Contains(t, diags[0].Message, "empty-sub")
	})

	t.Run("throws when boundedContexts is present but empty", func(t *testing.T) {
		yaml := minimalParents + `---
apiVersion: schema.esdm.io/core/v1
kind: subdomain
name: empty-sub
scope:
  domain: d
type: core
boundedContexts: []
`
		diags := runRule(t, rule, buildModelTolerantOfSchemaErrors(t, yaml))
		assert.Len(t, diags, 1)
	})
}
