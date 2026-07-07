package rules_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStoryUnknownGroupMembership(t *testing.T) {
	rule := findCatalogRule(t, "esdm/structure/story-unknown-group-membership")

	t.Run("does not throw when every group membership resolves to a declared group", func(t *testing.T) {
		yaml := storyParentsYAML + `---
apiVersion: schema.esdm.io/domain-storytelling/v1
kind: domain-story
name: well-known-groups-story
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
        groups:
          - front-office
`
		assert.Empty(t, runRule(t, rule, buildModel(t, yaml)))
	})

	t.Run("throws when an actor's groups membership references an undeclared group", func(t *testing.T) {
		yaml := storyParentsYAML + `---
apiVersion: schema.esdm.io/domain-storytelling/v1
kind: domain-story
name: dangling-actor-group-story
scope:
  domain: d
groups:
  - name: front-office
actors:
  - name: customer
    groups:
      - back-office
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
		assert.Contains(t, diags[0].Message, "back-office")
		assert.Contains(t, diags[0].Message, "customer")
	})

	t.Run("throws when an edge's groups membership references an undeclared group", func(t *testing.T) {
		yaml := storyParentsYAML + `---
apiVersion: schema.esdm.io/domain-storytelling/v1
kind: domain-story
name: dangling-edge-group-story
scope:
  domain: d
groups:
  - name: front-office
sentences:
  - sequenceNumber: 1
    edges:
      - from:
          actor: customer
        to:
          workObject: order
        groups:
          - back-office
`
		diags := runRule(t, rule, buildModel(t, yaml))
		require.Len(t, diags, 1)
		assert.Contains(t, diags[0].Message, "back-office")
	})

	t.Run("throws when a work-object's groups membership references an undeclared group", func(t *testing.T) {
		yaml := storyParentsYAML + `---
apiVersion: schema.esdm.io/domain-storytelling/v1
kind: domain-story
name: dangling-work-object-group-story
scope:
  domain: d
groups:
  - name: front-office
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
		diags := runRule(t, rule, buildModel(t, yaml))
		require.Len(t, diags, 1)
		assert.Contains(t, diags[0].Message, "back-office")
		assert.Contains(t, diags[0].Message, "ledger")
	})
}
