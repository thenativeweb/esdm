package rules_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDynamicConsistencyBoundaryWithoutConsults(t *testing.T) {
	rule := findCatalogRule(t, "esdm/modeling/dynamic-consistency-boundary-without-consults")

	t.Run("does not throw when the DCB consults at least one event", func(t *testing.T) {
		yaml := dcbParents + `---
apiVersion: schema.esdm.io/core/v1
kind: dynamic-consistency-boundary
name: capacity
scope:
  domain: d
  boundedContext: bc
identifiedBy:
  - name: id
    source: static
    value: solo
consults:
  - boundedContext: bc
    aggregate: agg
    event: enrolled
    criteria: relevant
`
		assert.Empty(t, runRule(t, rule, buildModel(t, yaml)))
	})

	t.Run("throws when the DCB has no consults field", func(t *testing.T) {
		yaml := dcbParents + `---
apiVersion: schema.esdm.io/core/v1
kind: dynamic-consistency-boundary
name: capacity
scope:
  domain: d
  boundedContext: bc
identifiedBy:
  - name: id
    source: static
    value: solo
`
		diags := runRule(t, rule, buildModelTolerantOfSchemaErrors(t, yaml))
		require.Len(t, diags, 1)
		assert.Contains(t, diags[0].Message, "capacity")
	})
}
