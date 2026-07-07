package rules_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEventWithoutConsumer(t *testing.T) {
	rule := findCatalogRule(t, "esdm/modeling/event-without-consumer")

	t.Run("throws when no consumer references the event", func(t *testing.T) {
		diags := runRule(t, rule, buildModel(t, minimalParents))
		require.Len(t, diags, 1)
		assert.Contains(t, diags[0].Message, "agg-done")
	})

	t.Run("does not throw when an event-handler handles the event", func(t *testing.T) {
		yaml := minimalParents + `---
apiVersion: schema.esdm.io/core/v1
kind: event-handler
name: notify
scope:
  domain: d
deliveryGuarantee: at-most-once
handles:
  - boundedContext: bc
    aggregate: agg
    event: agg-done
sideEffects:
  - type: other
    rule: log
`
		assert.Empty(t, runRule(t, rule, buildModel(t, yaml)))
	})
}
