package rules_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScenarioWithoutThen(t *testing.T) {
	rule := findCatalogRule(t, "esdm/gwt/scenario-without-then")

	t.Run("does not throw when every scenario carries a then", func(t *testing.T) {
		yaml := featureParentsYAML + `---
apiVersion: schema.esdm.io/given-when-then/v1
kind: feature
name: simple-feature
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

	t.Run("throws when a scenario is missing its then field", func(t *testing.T) {
		yaml := featureParentsYAML + `---
apiVersion: schema.esdm.io/given-when-then/v1
kind: feature
name: simple-feature
scope:
  domain: d
  boundedContext: bc
  aggregate: agg
scenarios:
  - name: hollow
    when:
      command: do-it
      data: {}
`
		diags := runRule(t, rule, buildModel(t, yaml))
		require.Len(t, diags, 1)
		assert.Contains(t, diags[0].Message, "hollow")
	})
}
