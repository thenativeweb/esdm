package rules_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadModelWithoutQuery(t *testing.T) {
	rule := findCatalogRule(t, "esdm/modeling/read-model-without-query")

	t.Run("does not throw when a query reads from the read-model", func(t *testing.T) {
		yaml := minimalParents + `---
apiVersion: schema.esdm.io/core/v1
kind: read-model
name: orders-list
scope:
  domain: d
  boundedContext: bc
projections:
  - boundedContext: bc
    aggregate: agg
    event: agg-done
    rule: append to list
schema:
  type: array
---
apiVersion: schema.esdm.io/core/v1
kind: query
name: list-orders
scope:
  domain: d
  boundedContext: bc
readModel: orders-list
result:
  type: array
parameters:
  type: object
actors:
  - user
`
		assert.Empty(t, runRule(t, rule, buildModel(t, yaml)))
	})

	t.Run("throws when the read-model has no query accessing it", func(t *testing.T) {
		yaml := minimalParents + `---
apiVersion: schema.esdm.io/core/v1
kind: read-model
name: orphan-rm
scope:
  domain: d
  boundedContext: bc
projections:
  - boundedContext: bc
    aggregate: agg
    event: agg-done
    rule: append to list
schema:
  type: array
`
		diags := runRule(t, rule, buildModel(t, yaml))
		require.Len(t, diags, 1)
		assert.Contains(t, diags[0].Message, "orphan-rm")
	})
}
