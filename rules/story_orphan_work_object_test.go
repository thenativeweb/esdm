package rules_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStoryOrphanWorkObject(t *testing.T) {
	rule := findCatalogRule(t, "esdm/modeling/story-orphan-work-object")

	t.Run("does not throw when every declared work-object is referenced from at least one edge in its sentence", func(t *testing.T) {
		yaml := storyParentsYAML + `---
apiVersion: schema.esdm.io/domain-storytelling/v1
kind: domain-story
name: connected-story
scope:
  domain: d
sentences:
  - sequenceNumber: 1
    workObjects:
      - name: order
        annotation: The customer's order.
    edges:
      - from:
          actor: customer
        to:
          workObject: order
`
		assert.Empty(t, runRule(t, rule, buildModel(t, yaml)))
	})

	t.Run("throws when a declared work-object never appears in its sentence's edges", func(t *testing.T) {
		yaml := storyParentsYAML + `---
apiVersion: schema.esdm.io/domain-storytelling/v1
kind: domain-story
name: orphan-object-story
scope:
  domain: d
sentences:
  - sequenceNumber: 1
    workObjects:
      - name: order
        annotation: Drawn in edges.
      - name: invoice
        annotation: Declared but not drawn here.
    edges:
      - from:
          actor: customer
        to:
          workObject: order
`
		diags := runRule(t, rule, buildModel(t, yaml))
		require.Len(t, diags, 1)
		assert.Contains(t, diags[0].Message, "invoice")
		assert.Contains(t, diags[0].Message, "orphan-object-story")
	})
}
