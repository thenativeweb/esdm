package rules_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScenarioReferencesUnknownActor(t *testing.T) {
	rule := findCatalogRule(t, "esdm/gwt/scenario-references-unknown-actor")

	t.Run("does not throw when the when.actor exists", func(t *testing.T) {
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
      actor: user
    then:
      events: []
`
		assert.Empty(t, runRule(t, rule, buildModel(t, yaml)))
	})

	t.Run("throws when the when.actor does not exist", func(t *testing.T) {
		yaml := minimalParents + `---
apiVersion: schema.esdm.io/given-when-then/v1
kind: feature
name: dangling
scope:
  domain: d
  boundedContext: bc
  aggregate: agg
scenarios:
  - name: trivial
    when:
      command: do-it
      data: {}
      actor: phantom-actor
    then:
      events: []
`
		diags := runRule(t, rule, buildModel(t, yaml))
		require.Len(t, diags, 1)
		assert.Contains(t, diags[0].Message, "phantom-actor")
	})
}
