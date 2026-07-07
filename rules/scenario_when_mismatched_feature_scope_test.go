package rules_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScenarioWhenMismatchedFeatureScope(t *testing.T) {
	rule := findCatalogRule(t, "esdm/gwt/scenario-when-mismatched-feature-scope")

	t.Run("does not throw when an aggregate feature carries a command-shaped when", func(t *testing.T) {
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

	t.Run("throws when an aggregate feature carries a query-shaped when", func(t *testing.T) {
		// Build a minimal model with an aggregate; the
		// schema would normally reject this document, but
		// the resolver still indexes it so the linter rule
		// can throw on the modeling concern.
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
      query: bogus
      parameters: {}
    then:
      events: []
`
		diags := runRule(t, rule, buildModel(t, yaml))
		require.Len(t, diags, 1)
		assert.Contains(t, diags[0].Message, "query")
	})
}
