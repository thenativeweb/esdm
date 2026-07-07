package rules_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScenarioReferencesUnknownEvent(t *testing.T) {
	rule := findCatalogRule(t, "esdm/gwt/scenario-references-unknown-event")

	t.Run("does not throw when every referenced event exists", func(t *testing.T) {
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
    given:
      - event: agg-done
        data: {}
    when:
      command: do-it
      data: {}
    then:
      events:
        - event: agg-done
          data: {}
`
		assert.Empty(t, runRule(t, rule, buildModel(t, yaml)))
	})

	t.Run("throws when a given event does not exist", func(t *testing.T) {
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
    given:
      - event: phantom-event
        data: {}
    when:
      command: do-it
      data: {}
    then:
      events: []
`
		diags := runRule(t, rule, buildModel(t, yaml))
		require.Len(t, diags, 1)
		assert.Contains(t, diags[0].Message, "phantom-event")
	})

	t.Run("throws when a then.events entry references an unknown event", func(t *testing.T) {
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
    then:
      events:
        - event: phantom-event
          data: {}
`
		diags := runRule(t, rule, buildModel(t, yaml))
		require.Len(t, diags, 1)
		assert.Contains(t, diags[0].Message, "phantom-event")
	})
}
