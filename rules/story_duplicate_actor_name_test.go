package rules_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStoryDuplicateActorName(t *testing.T) {
	rule := findCatalogRule(t, "esdm/structure/story-duplicate-actor-name")

	t.Run("does not throw when every actor name in the actors list is unique", func(t *testing.T) {
		yaml := storyParentsYAML + `---
apiVersion: schema.esdm.io/domain-storytelling/v1
kind: domain-story
name: unique-cast-story
scope:
  domain: d
actors:
  - name: customer
  - name: clerk
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

	t.Run("throws when the actors list declares the same name twice", func(t *testing.T) {
		yaml := storyParentsYAML + `---
apiVersion: schema.esdm.io/domain-storytelling/v1
kind: domain-story
name: doubled-cast-story
scope:
  domain: d
actors:
  - name: customer
    annotation: First declaration.
  - name: customer
    annotation: Accidental second declaration.
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
		assert.Contains(t, diags[0].Message, "customer")
		assert.Contains(t, diags[0].Message, "doubled-cast-story")
	})
}
