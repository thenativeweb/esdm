package rules_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStoryDuplicateGroupName(t *testing.T) {
	rule := findCatalogRule(t, "esdm/structure/story-duplicate-group-name")

	t.Run("does not throw when every group name in the registry is unique", func(t *testing.T) {
		yaml := storyParentsYAML + `---
apiVersion: schema.esdm.io/domain-storytelling/v1
kind: domain-story
name: unique-groups-story
scope:
  domain: d
groups:
  - name: front-office
  - name: back-office
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

	t.Run("throws when the groups registry declares the same name twice", func(t *testing.T) {
		yaml := storyParentsYAML + `---
apiVersion: schema.esdm.io/domain-storytelling/v1
kind: domain-story
name: doubled-groups-story
scope:
  domain: d
groups:
  - name: front-office
    description: First entry.
  - name: front-office
    description: Accidental second entry.
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
		assert.Contains(t, diags[0].Message, "front-office")
		assert.Contains(t, diags[0].Message, "doubled-groups-story")
	})
}
