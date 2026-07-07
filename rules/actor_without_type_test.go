package rules_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestActorWithoutType(t *testing.T) {
	rule := findCatalogRule(t, "esdm/modeling/actor-without-type")

	t.Run("does not throw when the actor declares a type", func(t *testing.T) {
		assert.Empty(t, runRule(t, rule, buildModel(t, minimalParents)))
	})

	t.Run("throws when the actor has no type field", func(t *testing.T) {
		yaml := `apiVersion: schema.esdm.io/core/v1
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
kind: actor
name: anonymous
scope:
  domain: d
  boundedContext: bc
`
		diags := runRule(t, rule, buildModelTolerantOfSchemaErrors(t, yaml))
		require.Len(t, diags, 1)
		assert.Contains(t, diags[0].Message, "anonymous")
	})
}
