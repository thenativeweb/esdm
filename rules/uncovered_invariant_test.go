package rules_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUncoveredInvariant(t *testing.T) {
	rule := findCatalogRule(t, "esdm/gwt/uncovered-invariant")

	t.Run("does not throw when a scenario rejection covers the invariant", func(t *testing.T) {
		yaml := aggregateWithInvariantParents + `---
apiVersion: schema.esdm.io/given-when-then/v1
kind: feature
name: covered
scope:
  domain: d
  boundedContext: bc
  aggregate: agg
scenarios:
  - name: rejects
    when:
      command: do-it
      data: {}
    then:
      rejection:
        invariant: shipping-blocks-cancellation
`
		assert.Empty(t, runRule(t, rule, buildModel(t, yaml)))
	})

	t.Run("throws when a Given-When-Then-tested unit leaves an invariant uncovered", func(t *testing.T) {
		yaml := aggregateWithInvariantParents + `---
apiVersion: schema.esdm.io/given-when-then/v1
kind: feature
name: happy-path-only
scope:
  domain: d
  boundedContext: bc
  aggregate: agg
scenarios:
  - name: succeeds
    when:
      command: do-it
      data: {}
    then:
      events:
        - event: agg-done
          data: {}
`
		diags := runRule(t, rule, buildModel(t, yaml))
		require.Len(t, diags, 1)
		assert.Contains(t, diags[0].Message, "shipping-blocks-cancellation")
		assert.Contains(t, diags[0].Message, "not covered by any scenario")
	})

	t.Run("does not throw when the unit has no Given-When-Then feature", func(t *testing.T) {
		assert.Empty(t, runRule(t, rule, buildModel(t, aggregateWithInvariantParents)))
	})
}
