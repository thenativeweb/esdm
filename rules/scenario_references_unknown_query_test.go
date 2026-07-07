package rules_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScenarioReferencesUnknownQuery(t *testing.T) {
	rule := findCatalogRule(t, "esdm/gwt/scenario-references-unknown-query")

	t.Run("does not throw when the when.query exists", func(t *testing.T) {
		yaml := readModelParents + `---
apiVersion: schema.esdm.io/given-when-then/v1
kind: feature
name: ok
scope:
  domain: d
  boundedContext: bc
  readModel: rm
scenarios:
  - name: trivial
    when:
      query: q
      parameters: {}
    then:
      result: []
`
		assert.Empty(t, runRule(t, rule, buildModel(t, yaml)))
	})

	t.Run("throws when the when.query does not exist", func(t *testing.T) {
		yaml := readModelParents + `---
apiVersion: schema.esdm.io/given-when-then/v1
kind: feature
name: dangling
scope:
  domain: d
  boundedContext: bc
  readModel: rm
scenarios:
  - name: trivial
    when:
      query: phantom-q
      parameters: {}
    then:
      result: []
`
		diags := runRule(t, rule, buildModel(t, yaml))
		require.Len(t, diags, 1)
		assert.Contains(t, diags[0].Message, "phantom-q")
	})
}
