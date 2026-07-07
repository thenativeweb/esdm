package rules_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScenarioThenMismatchedFeatureScope(t *testing.T) {
	rule := findCatalogRule(t, "esdm/gwt/scenario-then-mismatched-feature-scope")

	t.Run("does not throw when an aggregate feature carries an events-shaped then", func(t *testing.T) {
		yaml := minimalParents + `---
apiVersion: schema.esdm.io/given-when-then/v1
kind: feature
name: ok
scope:
  domain: d
  boundedContext: bc
  aggregate: agg
scenarios:
  - name: trivial
    when:
      command: do-it
      data: {}
    then:
      events: []
`
		assert.Empty(t, runRule(t, rule, buildModel(t, yaml)))
	})

	t.Run("throws when an aggregate feature carries a result-shaped then", func(t *testing.T) {
		yaml := minimalParents + `---
apiVersion: schema.esdm.io/given-when-then/v1
kind: feature
name: mismatched
scope:
  domain: d
  boundedContext: bc
  aggregate: agg
scenarios:
  - name: trivial
    when:
      command: do-it
      data: {}
    then:
      result: []
`
		diags := runRule(t, rule, buildModel(t, yaml))
		require.Len(t, diags, 1)
		assert.Contains(t, diags[0].Message, "result")
	})
}
