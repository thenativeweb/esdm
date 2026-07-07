package rules_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStoryOrphanGroup(t *testing.T) {
	rule := findCatalogRule(t, "esdm/modeling/story-orphan-group")

	t.Run("does not throw when every declared group is referenced from at least one element", func(t *testing.T) {
		yaml := storyParentsYAML + `---
apiVersion: schema.esdm.io/domain-storytelling/v1
kind: domain-story
name: well-grouped-story
scope:
  domain: d
groups:
  - name: front-office
  - name: back-office
actors:
  - name: customer
    groups:
      - front-office
sentences:
  - sequenceNumber: 1
    workObjects:
      - name: ledger
        groups:
          - back-office
    edges:
      - from:
          actor: customer
        to:
          workObject: ledger
`
		assert.Empty(t, runRule(t, rule, buildModel(t, yaml)))
	})

	t.Run("throws when a declared group has no members", func(t *testing.T) {
		yaml := storyParentsYAML + `---
apiVersion: schema.esdm.io/domain-storytelling/v1
kind: domain-story
name: lonely-group-story
scope:
  domain: d
groups:
  - name: front-office
  - name: tax-office
actors:
  - name: customer
    groups:
      - front-office
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
		assert.Contains(t, diags[0].Message, "tax-office")
		assert.Contains(t, diags[0].Message, "lonely-group-story")
	})
}
