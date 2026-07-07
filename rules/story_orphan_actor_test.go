package rules_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStoryOrphanActor(t *testing.T) {
	rule := findCatalogRule(t, "esdm/modeling/story-orphan-actor")

	t.Run("does not throw when every declared actor is referenced from at least one edge", func(t *testing.T) {
		yaml := storyParentsYAML + `---
apiVersion: schema.esdm.io/domain-storytelling/v1
kind: domain-story
name: connected-story
scope:
  domain: d
actors:
  - name: customer
    annotation: The retail customer.
sentences:
  - sequenceNumber: 1
    edges:
      - from:
          actor: customer
        to:
          workObject: order
`
		assert.Empty(t, runRule(t, rule, buildModel(t, yaml)))
	})

	t.Run("throws when a declared actor never appears in any edge", func(t *testing.T) {
		yaml := storyParentsYAML + `---
apiVersion: schema.esdm.io/domain-storytelling/v1
kind: domain-story
name: lonely-cast-story
scope:
  domain: d
actors:
  - name: customer
    annotation: Used in edges.
  - name: courier
    annotation: Declared but never drawn.
sentences:
  - sequenceNumber: 1
    edges:
      - from:
          actor: customer
        to:
          workObject: order
`
		diags := runRule(t, rule, buildModel(t, yaml))
		require.Len(t, diags, 1)
		assert.Contains(t, diags[0].Message, "courier")
		assert.Contains(t, diags[0].Message, "lonely-cast-story")
	})

	t.Run("does not throw for actors that are referenced from edges but not in the actors list", func(t *testing.T) {
		yaml := storyParentsYAML + `---
apiVersion: schema.esdm.io/domain-storytelling/v1
kind: domain-story
name: implicit-actor-story
scope:
  domain: d
sentences:
  - sequenceNumber: 1
    edges:
      - from:
          actor: customer
        to:
          workObject: order
`
		assert.Empty(t, runRule(t, rule, buildModel(t, yaml)))
	})
}
