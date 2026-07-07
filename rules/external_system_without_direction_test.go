package rules_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExternalSystemWithoutDirection(t *testing.T) {
	rule := findCatalogRule(t, "esdm/modeling/external-system-without-direction")

	t.Run("does not throw when the external-system declares a direction", func(t *testing.T) {
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

	t.Run("throws when the external-system has no direction field", func(t *testing.T) {
		yaml := minimalParents + `---
apiVersion: schema.esdm.io/core/v1
kind: external-system
name: directionless
scope:
  domain: d
`
		diags := runRule(t, rule, buildModelTolerantOfSchemaErrors(t, yaml))
		require.Len(t, diags, 1)
		assert.Contains(t, diags[0].Message, "directionless")
	})
}
