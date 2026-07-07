package rules_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQueryWithoutReadModel(t *testing.T) {
	rule := findCatalogRule(t, "esdm/modeling/query-without-read-model")

	t.Run("does not throw when the query references a readModel", func(t *testing.T) {
		yaml := queryParents + `---
apiVersion: schema.esdm.io/core/v1
kind: query
name: q
scope:
  domain: d
  boundedContext: bc
readModel: rm
result:
  type: object
`
		assert.Empty(t, runRule(t, rule, buildModel(t, yaml)))
	})

	t.Run("throws when the query has no readModel field", func(t *testing.T) {
		yaml := queryParents + `---
apiVersion: schema.esdm.io/core/v1
kind: query
name: q
scope:
  domain: d
  boundedContext: bc
result:
  type: object
`
		diags := runRule(t, rule, buildModelTolerantOfSchemaErrors(t, yaml))
		require.Len(t, diags, 1)
		assert.Contains(t, diags[0].Message, "q")
	})
}
