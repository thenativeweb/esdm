package rules_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStoryDuplicateWorkObjectInSentence(t *testing.T) {
	rule := findCatalogRule(t, "esdm/structure/story-duplicate-work-object-in-sentence")

	t.Run("does not throw when work-object names are unique within each sentence", func(t *testing.T) {
		yaml := storyParentsYAML + `---
apiVersion: schema.esdm.io/domain-storytelling/v1
kind: domain-story
name: unique-objects-story
scope:
  domain: d
sentences:
  - sequenceNumber: 1
    workObjects:
      - name: order
      - name: invoice
    edges:
      - from:
          actor: customer
        to:
          workObject: order
`
		assert.Empty(t, runRule(t, rule, buildModel(t, yaml)))
	})

	t.Run("does not throw when the same work-object name reappears in a different sentence", func(t *testing.T) {
		yaml := storyParentsYAML + `---
apiVersion: schema.esdm.io/domain-storytelling/v1
kind: domain-story
name: redrawn-object-story
scope:
  domain: d
sentences:
  - sequenceNumber: 1
    workObjects:
      - name: order
    edges:
      - from:
          actor: customer
        to:
          workObject: order
  - sequenceNumber: 2
    workObjects:
      - name: order
    edges:
      - from:
          actor: clerk
        to:
          workObject: order
`
		assert.Empty(t, runRule(t, rule, buildModel(t, yaml)))
	})

	t.Run("throws when a sentence's workObjects list declares the same name twice", func(t *testing.T) {
		yaml := storyParentsYAML + `---
apiVersion: schema.esdm.io/domain-storytelling/v1
kind: domain-story
name: doubled-object-in-sentence-story
scope:
  domain: d
sentences:
  - sequenceNumber: 1
    workObjects:
      - name: order
        annotation: First declaration.
      - name: order
        annotation: Accidental second declaration in the same sentence.
    edges:
      - from:
          actor: customer
        to:
          workObject: order
`
		diags := runRule(t, rule, buildModel(t, yaml))
		require.Len(t, diags, 1)
		assert.Contains(t, diags[0].Message, "order")
		assert.Contains(t, diags[0].Message, "doubled-object-in-sentence-story")
	})
}
