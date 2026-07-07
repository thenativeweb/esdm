package rules_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEntityWithoutIdentifiedBy(t *testing.T) {
	rule := findCatalogRule(t, "esdm/modeling/entity-without-identified-by")

	t.Run("does not throw when the entity declares identifiedBy", func(t *testing.T) {
		yaml := minimalParents + `---
apiVersion: schema.esdm.io/core/v1
kind: entity
name: student
scope:
  domain: d
  boundedContext: bc
schema:
  type: object
  properties:
    id:
      type: string
identifiedBy:
  source: schema
  field: id
`
		assert.Empty(t, runRule(t, rule, buildModel(t, yaml)))
	})

	t.Run("throws when the entity has no identifiedBy field", func(t *testing.T) {
		yaml := minimalParents + `---
apiVersion: schema.esdm.io/core/v1
kind: entity
name: student
scope:
  domain: d
  boundedContext: bc
schema:
  type: object
`
		diags := runRule(t, rule, buildModelTolerantOfSchemaErrors(t, yaml))
		require.Len(t, diags, 1)
		assert.Contains(t, diags[0].Message, "student")
	})
}
