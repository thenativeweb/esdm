package rules_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// featureParentsYAML is the smallest valid prefix that
// lets a given-when-then feature resolve - just the
// domain it scopes to. Feature-rule tests prepend it to
// the YAML they build per case. Tests that need richer
// surrounding entities (aggregates, commands, events,
// process managers, read-models, queries, actors) extend
// this prefix per case.
const featureParentsYAML = `apiVersion: schema.esdm.io/core/v1
kind: domain
name: d
`

func TestFeatureWithoutScenarios(t *testing.T) {
	rule := findCatalogRule(t, "esdm/gwt/feature-without-scenarios")

	t.Run("does not throw when the feature has at least one scenario", func(t *testing.T) {
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

	t.Run("throws when the feature has no scenarios", func(t *testing.T) {
		yaml := featureParentsYAML + `---
apiVersion: schema.esdm.io/given-when-then/v1
kind: feature
name: hollow-feature
scope:
  domain: d
  boundedContext: bc
  aggregate: agg
`
		diags := runRule(t, rule, buildModel(t, yaml))
		require.Len(t, diags, 1)
		assert.Contains(t, diags[0].Message, "hollow-feature")
	})
}
