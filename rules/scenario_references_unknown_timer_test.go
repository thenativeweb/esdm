package rules_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScenarioReferencesUnknownTimer(t *testing.T) {
	rule := findCatalogRule(t, "esdm/gwt/scenario-references-unknown-timer")

	t.Run("does not throw when the when.timer exists on the process-manager", func(t *testing.T) {
		yaml := processManagerParents + `---
apiVersion: schema.esdm.io/given-when-then/v1
kind: feature
name: ok
scope:
  domain: d
  processManager: order-pm
scenarios:
  - name: trivial
    when:
      timer: deadline
    then:
      ended: true
`
		assert.Empty(t, runRule(t, rule, buildModel(t, yaml)))
	})

	t.Run("throws when the when.timer does not exist on the process-manager", func(t *testing.T) {
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
      timer: phantom-timer
    then:
      ended: true
`
		diags := runRule(t, rule, buildModel(t, yaml))
		require.Len(t, diags, 1)
		assert.Contains(t, diags[0].Message, "phantom-timer")
	})
}
