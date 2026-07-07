package rules_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFeatureReferencesUnknownAggregate(t *testing.T) {
	rule := findCatalogRule(t, "esdm/gwt/feature-references-unknown-aggregate")

	t.Run("does not throw when the feature scope points at an existing aggregate", func(t *testing.T) {
		yaml := minimalParents + `---
apiVersion: schema.esdm.io/given-when-then/v1
kind: feature
name: ok-feature
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

	t.Run("throws when the feature scope points at an aggregate that does not exist", func(t *testing.T) {
		yaml := minimalParents + `---
apiVersion: schema.esdm.io/given-when-then/v1
kind: feature
name: dangling-feature
scope:
  domain: d
  boundedContext: bc
  aggregate: phantom-agg
scenarios:
  - name: trivial
    when:
      command: do-it
      data: {}
    then:
      events: []
`
		diags := runRule(t, rule, buildModel(t, yaml))
		require.Len(t, diags, 1)
		assert.Contains(t, diags[0].Message, "phantom-agg")
	})
}
