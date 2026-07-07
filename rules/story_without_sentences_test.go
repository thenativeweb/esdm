package rules_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// storyParentsYAML is the smallest valid prefix that lets a
// domain-story document resolve - just the domain it
// scopes to. Story-rule tests prepend it to the YAML they
// build per case.
const storyParentsYAML = `apiVersion: schema.esdm.io/core/v1
kind: domain
name: d
`

func TestStoryWithoutSentences(t *testing.T) {
	rule := findCatalogRule(t, "esdm/structure/story-without-sentences")

	t.Run("does not throw when the story has at least one sentence", func(t *testing.T) {
		yaml := storyParentsYAML + `---
apiVersion: schema.esdm.io/domain-storytelling/v1
kind: domain-story
name: simple-story
scope:
  domain: d
sentences:
  - sequenceNumber: 1
    edges:
      - from:
          actor: a
        to:
          workObject: o
`
		assert.Empty(t, runRule(t, rule, buildModel(t, yaml)))
	})

	t.Run("throws when the story has no sentences", func(t *testing.T) {
		yaml := storyParentsYAML + `---
apiVersion: schema.esdm.io/domain-storytelling/v1
kind: domain-story
name: hollow-story
scope:
  domain: d
`
		diags := runRule(t, rule, buildModel(t, yaml))
		require.Len(t, diags, 1)
		assert.Contains(t, diags[0].Message, "hollow-story")
	})
}
