package rules_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const queryParents = minimalParents + `---
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

func TestQueryWithoutResult(t *testing.T) {
	rule := findCatalogRule(t, "esdm/modeling/query-without-result")

	t.Run("does not throw when the query declares result", func(t *testing.T) {
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

	t.Run("throws when the query has no result field", func(t *testing.T) {
		yaml := queryParents + `---
apiVersion: schema.esdm.io/core/v1
kind: query
name: q
scope:
  domain: d
  boundedContext: bc
readModel: rm
`
		diags := runRule(t, rule, buildModelTolerantOfSchemaErrors(t, yaml))
		require.Len(t, diags, 1)
		assert.Contains(t, diags[0].Message, "q")
	})
}
