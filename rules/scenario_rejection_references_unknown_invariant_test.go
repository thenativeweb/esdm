package rules_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const aggregateWithInvariantParents = `apiVersion: schema.esdm.io/core/v1
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
kind: aggregate
name: agg
scope:
  domain: d
  boundedContext: bc
identifiedBy:
  source: generated
  generator: uuid
state:
  type: object
invariants:
  - name: shipping-blocks-cancellation
    rule: a shipped order may not be canceled
---
apiVersion: schema.esdm.io/core/v1
kind: command
name: do-it
scope:
  domain: d
  boundedContext: bc
  aggregate: agg
data:
  type: object
publishes:
  - agg-done
---
apiVersion: schema.esdm.io/core/v1
kind: event
name: agg-done
scope:
  domain: d
  boundedContext: bc
  aggregate: agg
data:
  type: object
`

func TestScenarioRejectionReferencesUnknownInvariant(t *testing.T) {
	rule := findCatalogRule(t, "esdm/gwt/scenario-rejection-references-unknown-invariant")

	t.Run("does not throw when the rejection invariant is declared on the aggregate", func(t *testing.T) {
		yaml := aggregateWithInvariantParents + `---
apiVersion: schema.esdm.io/given-when-then/v1
kind: feature
name: ok
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

	t.Run("throws when the rejection points at an invariant the aggregate does not declare", func(t *testing.T) {
		yaml := aggregateWithInvariantParents + `---
apiVersion: schema.esdm.io/given-when-then/v1
kind: feature
name: dangling
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
        invariant: phantom-invariant
`
		diags := runRule(t, rule, buildModel(t, yaml))
		require.Len(t, diags, 1)
		assert.Contains(t, diags[0].Message, "phantom-invariant")
	})

	t.Run("does not throw when the rejection uses free-form prose", func(t *testing.T) {
		yaml := aggregateWithInvariantParents + `---
apiVersion: schema.esdm.io/given-when-then/v1
kind: feature
name: prose
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
        reason: too many reasons to list
`
		assert.Empty(t, runRule(t, rule, buildModel(t, yaml)))
	})
}
