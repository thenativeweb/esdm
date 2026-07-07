package rules_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStorySentenceWithoutEdges(t *testing.T) {
	rule := findCatalogRule(t, "esdm/structure/story-sentence-without-edges")

	t.Run("does not throw when every sentence carries at least one edge", func(t *testing.T) {
		yaml := storyParentsYAML + `---
apiVersion: schema.esdm.io/domain-storytelling/v1
kind: domain-story
name: full-story
scope:
  domain: d
sentences:
  - sequenceNumber: 1
    edges:
      - from:
          actor: a
        to:
          workObject: o
  - sequenceNumber: 2
    edges:
      - from:
          actor: a
        to:
          workObject: o
`
		assert.Empty(t, runRule(t, rule, buildModel(t, yaml)))
	})

	t.Run("throws when a sentence carries no edges", func(t *testing.T) {
		yaml := storyParentsYAML + `---
apiVersion: schema.esdm.io/domain-storytelling/v1
kind: domain-story
name: gappy-story
scope:
  domain: d
sentences:
  - sequenceNumber: 1
    edges:
      - from:
          actor: a
        to:
          workObject: o
  - sequenceNumber: 2
    edges: []
`
		diags := runRule(t, rule, buildModel(t, yaml))
		require.Len(t, diags, 1)
		assert.Contains(t, diags[0].Message, "gappy-story")
		assert.Contains(t, diags[0].Message, "2")
	})
}
