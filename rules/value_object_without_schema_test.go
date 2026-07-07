package rules_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValueObjectWithoutSchema(t *testing.T) {
	rule := findCatalogRule(t, "esdm/modeling/value-object-without-schema")

	t.Run("does not throw when the value-object declares a schema", func(t *testing.T) {
		yaml := minimalParents + `---
apiVersion: schema.esdm.io/core/v1
kind: value-object
name: money
scope:
  domain: d
  boundedContext: bc
schema:
  type: object
`
		assert.Empty(t, runRule(t, rule, buildModel(t, yaml)))
	})

	t.Run("throws when the value-object has no schema field", func(t *testing.T) {
		yaml := minimalParents + `---
apiVersion: schema.esdm.io/core/v1
kind: value-object
name: money
scope:
  domain: d
  boundedContext: bc
`
		diags := runRule(t, rule, buildModelTolerantOfSchemaErrors(t, yaml))
		require.Len(t, diags, 1)
		assert.Contains(t, diags[0].Message, "money")
	})
}
