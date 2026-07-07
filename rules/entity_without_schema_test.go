package rules_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEntityWithoutSchema(t *testing.T) {
	rule := findCatalogRule(t, "esdm/modeling/entity-without-schema")

	t.Run("does not throw when the entity declares a schema", func(t *testing.T) {
		yaml := minimalParents + `---
apiVersion: schema.esdm.io/core/v1
kind: entity
name: student
scope:
  domain: d
  boundedContext: bc
schema:
  type: object
identifiedBy:
  source: static
  value: the-one
`
		assert.Empty(t, runRule(t, rule, buildModel(t, yaml)))
	})

	t.Run("throws when the entity has no schema field", func(t *testing.T) {
		yaml := minimalParents + `---
apiVersion: schema.esdm.io/core/v1
kind: entity
name: student
scope:
  domain: d
  boundedContext: bc
identifiedBy:
  source: static
  value: the-one
`
		diags := runRule(t, rule, buildModelTolerantOfSchemaErrors(t, yaml))
		require.Len(t, diags, 1)
		assert.Contains(t, diags[0].Message, "student")
	})
}
