package rules_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOrphanActor(t *testing.T) {
	rule := findCatalogRule(t, "esdm/modeling/orphan-actor")

	t.Run("does not throw when an actor is named in a command's actors list", func(t *testing.T) {
		// minimalParents already ties `user` into `do-it.actors`.
		assert.Empty(t, runRule(t, rule, buildModel(t, minimalParents)))
	})

	t.Run("throws when an actor is declared but never referenced by a command", func(t *testing.T) {
		yaml := minimalParents + `---
apiVersion: schema.esdm.io/core/v1
kind: actor
name: bystander
scope:
  domain: d
  boundedContext: bc
type: human
`
		diags := runRule(t, rule, buildModel(t, yaml))
		require.Len(t, diags, 1)
		assert.Contains(t, diags[0].Message, "bystander")
	})

	t.Run("does not throw when an actor is referenced by a query.actors list", func(t *testing.T) {
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
---
apiVersion: schema.esdm.io/core/v1
kind: query
name: q
scope:
  domain: d
  boundedContext: bc
readModel: rm
result:
  type: object
actors:
  - user
`
		diags := runRule(t, rule, buildModel(t, yaml))
		assert.Empty(t, diags)
	})
}
