package rules_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOrphanExternalSystem(t *testing.T) {
	rule := findCatalogRule(t, "esdm/modeling/orphan-external-system")

	t.Run("does not throw when an actor is backed by the external system", func(t *testing.T) {
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
name: billing-system
scope:
  domain: d
  boundedContext: bc
type: system
backedBy:
  - stripe
`
		assert.Empty(t, runRule(t, rule, buildModel(t, yaml)))
	})

	t.Run("throws when nothing references the external system", func(t *testing.T) {
		yaml := minimalParents + `---
apiVersion: schema.esdm.io/core/v1
kind: external-system
name: lonely-system
scope:
  domain: d
direction: outbound
`
		diags := runRule(t, rule, buildModel(t, yaml))
		require.Len(t, diags, 1)
		assert.Contains(t, diags[0].Message, "lonely-system")
	})
}
