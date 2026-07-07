package rules_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHumanActorWithBackedBy(t *testing.T) {
	rule := findCatalogRule(t, "esdm/modeling/human-actor-with-backed-by")

	t.Run("does not throw when a human actor has no backedBy", func(t *testing.T) {
		assert.Empty(t, runRule(t, rule, buildModel(t, minimalParents)))
	})

	t.Run("does not throw when a system actor has backedBy", func(t *testing.T) {
		yaml := minimalParents + `---
apiVersion: schema.esdm.io/core/v1
kind: external-system
name: stripe
scope:
  domain: d
direction: outbound
---
apiVersion: schema.esdm.io/core/v1
kind: actor
name: billing
scope:
  domain: d
  boundedContext: bc
type: system
backedBy:
  - stripe
`
		assert.Empty(t, runRule(t, rule, buildModel(t, yaml)))
	})

	t.Run("throws when a human actor declares backedBy", func(t *testing.T) {
		yaml := minimalParents + `---
apiVersion: schema.esdm.io/core/v1
kind: external-system
name: stripe
scope:
  domain: d
direction: outbound
---
apiVersion: schema.esdm.io/core/v1
kind: actor
name: customer
scope:
  domain: d
  boundedContext: bc
type: human
backedBy:
  - stripe
`
		diags := runRule(t, rule, buildModelTolerantOfSchemaErrors(t, yaml))
		require.Len(t, diags, 1)
		assert.Contains(t, diags[0].Message, "customer")
	})
}
