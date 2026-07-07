package rules_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadModelWithoutSchema(t *testing.T) {
	rule := findCatalogRule(t, "esdm/modeling/read-model-without-schema")

	t.Run("does not throw when the read-model declares a schema", func(t *testing.T) {
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

	t.Run("throws when the read-model has no schema field", func(t *testing.T) {
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
`
		diags := runRule(t, rule, buildModelTolerantOfSchemaErrors(t, yaml))
		require.Len(t, diags, 1)
		assert.Contains(t, diags[0].Message, "rm")
	})
}
