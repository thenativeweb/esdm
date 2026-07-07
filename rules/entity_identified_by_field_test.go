package rules_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEntityIdentifiedByField(t *testing.T) {
	rule := findCatalogRule(t, "esdm/structure/entity-identified-by-field")

	t.Run("does not throw when identifiedBy.field matches a schema property", func(t *testing.T) {
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

	t.Run("does not throw when identifiedBy uses source: static", func(t *testing.T) {
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

	t.Run("throws when identifiedBy.field is missing from schema.properties", func(t *testing.T) {
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
    name:
      type: string
identifiedBy:
  source: schema
  field: id
`
		diags := runRule(t, rule, buildModel(t, yaml))
		require.Len(t, diags, 1)
		assert.Contains(t, diags[0].Message, "student")
		assert.Contains(t, diags[0].Message, "id")
	})
}
