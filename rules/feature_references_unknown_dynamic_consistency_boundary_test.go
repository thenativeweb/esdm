package rules_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const featureDcbParents = `apiVersion: schema.esdm.io/core/v1
kind: domain
name: d
---
apiVersion: schema.esdm.io/core/v1
kind: bounded-context
name: bc
scope:
  domain: d
---
apiVersion: schema.esdm.io/core/v1
kind: dynamic-consistency-boundary
name: dcb-unit
scope:
  domain: d
  boundedContext: bc
identifiedBy:
  - name: key
    source: command-payload
    field: key
consults:
  - boundedContext: bc
    event: some-event
    criteria: all
---
apiVersion: schema.esdm.io/core/v1
kind: command
name: dcb-cmd
scope:
  domain: d
  boundedContext: bc
  dynamicConsistencyBoundary: dcb-unit
data:
  type: object
publishes:
  - some-event
`

func TestFeatureReferencesUnknownDynamicConsistencyBoundary(t *testing.T) {
	rule := findCatalogRule(t, "esdm/gwt/feature-references-unknown-dynamic-consistency-boundary")

	t.Run("does not throw when the feature scope points at an existing dynamic-consistency-boundary", func(t *testing.T) {
		yaml := featureDcbParents + `---
apiVersion: schema.esdm.io/given-when-then/v1
kind: feature
name: ok-feature
scope:
  domain: d
  boundedContext: bc
  dynamicConsistencyBoundary: dcb-unit
scenarios:
  - name: trivial
    when:
      command: dcb-cmd
      data: {}
    then:
      events: []
`
		assert.Empty(t, runRule(t, rule, buildModel(t, yaml)))
	})

	t.Run("throws when the feature scope points at a dynamic-consistency-boundary that does not exist", func(t *testing.T) {
		yaml := featureDcbParents + `---
apiVersion: schema.esdm.io/given-when-then/v1
kind: feature
name: dangling-feature
scope:
  domain: d
  boundedContext: bc
  dynamicConsistencyBoundary: phantom-dcb
scenarios:
  - name: trivial
    when:
      command: dcb-cmd
      data: {}
    then:
      events: []
`
		diags := runRule(t, rule, buildModel(t, yaml))
		require.Len(t, diags, 1)
		assert.Contains(t, diags[0].Message, "phantom-dcb")
	})
}
