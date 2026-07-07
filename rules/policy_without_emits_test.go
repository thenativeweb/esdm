package rules_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPolicyWithoutEmits(t *testing.T) {
	rule := findCatalogRule(t, "esdm/modeling/policy-without-emits")

	t.Run("does not throw when the policy emits at least one command", func(t *testing.T) {
		yaml := minimalParents + `---
apiVersion: schema.esdm.io/core/v1
kind: policy
name: react
scope:
  domain: d
deliveryGuarantee: at-most-once
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

	t.Run("throws when the policy has no emits field", func(t *testing.T) {
		yaml := minimalParents + `---
apiVersion: schema.esdm.io/core/v1
kind: policy
name: passive
scope:
  domain: d
deliveryGuarantee: at-most-once
handles:
  - boundedContext: bc
    aggregate: agg
    event: agg-done
`
		diags := runRule(t, rule, buildModelTolerantOfSchemaErrors(t, yaml))
		require.Len(t, diags, 1)
		assert.Contains(t, diags[0].Message, "passive")
	})
}
