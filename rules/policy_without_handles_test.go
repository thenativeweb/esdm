package rules_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPolicyWithoutHandles(t *testing.T) {
	rule := findCatalogRule(t, "esdm/modeling/policy-without-handles")

	t.Run("does not throw when the policy handles at least one event", func(t *testing.T) {
		yaml := minimalParents + `---
apiVersion: schema.esdm.io/core/v1
kind: policy
name: react
scope:
  domain: d
deliveryGuarantee: at-least-once
idempotency:
  approach: inbox
  storage:
    inbox: in-memory
handles:
  - boundedContext: bc
    aggregate: agg
    event: agg-done
emits:
  - boundedContext: bc
    aggregate: agg
    command: do-it
`
		assert.Empty(t, runRule(t, rule, buildModel(t, yaml)))
	})

	t.Run("throws when the policy has no handles field", func(t *testing.T) {
		yaml := minimalParents + `---
apiVersion: schema.esdm.io/core/v1
kind: policy
name: idle
scope:
  domain: d
deliveryGuarantee: at-most-once
emits:
  - boundedContext: bc
    aggregate: agg
    command: do-it
`
		diags := runRule(t, rule, buildModelTolerantOfSchemaErrors(t, yaml))
		require.Len(t, diags, 1)
		assert.Contains(t, diags[0].Message, "idle")
	})
}
