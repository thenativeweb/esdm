package rules_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScenarioReferencesUnknownCommand(t *testing.T) {
	rule := findCatalogRule(t, "esdm/gwt/scenario-references-unknown-command")

	t.Run("does not throw when the when.command exists", func(t *testing.T) {
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

	t.Run("throws when the when.command does not exist", func(t *testing.T) {
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
      command: phantom-command
      data: {}
    then:
      events: []
`
		diags := runRule(t, rule, buildModel(t, yaml))
		require.Len(t, diags, 1)
		assert.Contains(t, diags[0].Message, "phantom-command")
	})

	t.Run("throws when a process-manager then.emits points at an unknown command", func(t *testing.T) {
		yaml := processManagerParents + `---
apiVersion: schema.esdm.io/given-when-then/v1
kind: feature
name: dangling
scope:
  domain: d
  processManager: order-pm
scenarios:
  - name: trivial
    when:
      timer: deadline
    then:
      emits:
        - boundedContext: bc
          aggregate: agg
          command: phantom-command
`
		diags := runRule(t, rule, buildModel(t, yaml))
		require.Len(t, diags, 1)
		assert.Contains(t, diags[0].Message, "phantom-command")
	})
}
