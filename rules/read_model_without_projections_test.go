package rules_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadModelWithoutProjections(t *testing.T) {
	rule := findCatalogRule(t, "esdm/modeling/read-model-without-projections")

	t.Run("does not throw when the read-model has at least one projection", func(t *testing.T) {
		yaml := minimalParents + `---
apiVersion: schema.esdm.io/core/v1
kind: read-model
name: rm
scope:
  domain: d
  boundedContext: bc
projections:
  - boundedContext: bc
    aggregate: agg
    event: agg-done
    rule: tally events
schema:
  type: object
`
		assert.Empty(t, runRule(t, rule, buildModel(t, yaml)))
	})

	t.Run("throws when the read-model has no projections field", func(t *testing.T) {
		// Same shape as the subdomain test: the JSON Schema
		// rejects this today (minItems: 1, required), so the
		// rule is the defense-in-depth net for a hypothetical
		// schema relaxation. The tolerant model builder lets
		// us exercise the rule directly.
		yaml := minimalParents + `---
apiVersion: schema.esdm.io/core/v1
kind: read-model
name: rm
scope:
  domain: d
  boundedContext: bc
schema:
  type: object
`
		diags := runRule(t, rule, buildModelTolerantOfSchemaErrors(t, yaml))
		require.Len(t, diags, 1)
		assert.Contains(t, diags[0].Message, "rm")
	})

	t.Run("throws when projections is present but empty", func(t *testing.T) {
		yaml := minimalParents + `---
apiVersion: schema.esdm.io/core/v1
kind: read-model
name: rm
scope:
  domain: d
  boundedContext: bc
projections: []
schema:
  type: object
`
		diags := runRule(t, rule, buildModelTolerantOfSchemaErrors(t, yaml))
		assert.Len(t, diags, 1)
	})
}
