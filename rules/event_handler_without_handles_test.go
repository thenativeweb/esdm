package rules_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEventHandlerWithoutHandles(t *testing.T) {
	rule := findCatalogRule(t, "esdm/modeling/event-handler-without-handles")

	t.Run("does not throw when the event-handler handles at least one event", func(t *testing.T) {
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
    rule: send mail
`
		assert.Empty(t, runRule(t, rule, buildModel(t, yaml)))
	})

	t.Run("throws when the event-handler has no handles field", func(t *testing.T) {
		yaml := minimalParents + `---
apiVersion: schema.esdm.io/core/v1
kind: event-handler
name: idle
scope:
  domain: d
deliveryGuarantee: at-most-once
sideEffects:
  - type: other
    rule: do nothing
`
		diags := runRule(t, rule, buildModelTolerantOfSchemaErrors(t, yaml))
		require.Len(t, diags, 1)
		assert.Contains(t, diags[0].Message, "idle")
	})
}
