package rules_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConsumerAtLeastOnceWithoutIdempotency(t *testing.T) {
	rule := findCatalogRule(t, "esdm/modeling/consumer-at-least-once-without-idempotency")

	t.Run("does not throw when at-most-once consumer omits idempotency", func(t *testing.T) {
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

	t.Run("does not throw when at-least-once consumer declares idempotency", func(t *testing.T) {
		yaml := minimalParents + `---
apiVersion: schema.esdm.io/core/v1
kind: event-handler
name: notify
scope:
  domain: d
deliveryGuarantee: at-least-once
idempotency:
  approach: inbox
  storage:
    inbox: in-memory
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

	t.Run("throws when an event-handler is at-least-once without idempotency", func(t *testing.T) {
		yaml := minimalParents + `---
apiVersion: schema.esdm.io/core/v1
kind: event-handler
name: dupe-prone
scope:
  domain: d
deliveryGuarantee: at-least-once
handles:
  - boundedContext: bc
    aggregate: agg
    event: agg-done
sideEffects:
  - type: other
    rule: send mail
`
		diags := runRule(t, rule, buildModelTolerantOfSchemaErrors(t, yaml))
		require.Len(t, diags, 1)
		assert.Contains(t, diags[0].Message, "event-handler")
		assert.Contains(t, diags[0].Message, "dupe-prone")
	})

	t.Run("throws when a policy is at-least-once without idempotency", func(t *testing.T) {
		yaml := minimalParents + `---
apiVersion: schema.esdm.io/core/v1
kind: policy
name: dupe-prone
scope:
  domain: d
deliveryGuarantee: at-least-once
handles:
  - boundedContext: bc
    aggregate: agg
    event: agg-done
emits:
  - boundedContext: bc
    aggregate: agg
    command: do-it
`
		diags := runRule(t, rule, buildModelTolerantOfSchemaErrors(t, yaml))
		require.Len(t, diags, 1)
		assert.Contains(t, diags[0].Message, "policy")
	})

	t.Run("throws when a process-manager is at-least-once without idempotency", func(t *testing.T) {
		yaml := pmParents + `---
apiVersion: schema.esdm.io/core/v1
kind: event
name: started
scope:
  domain: d
  boundedContext: bc
  aggregate: agg
data:
  type: object
---
apiVersion: schema.esdm.io/core/v1
kind: process-manager
name: dupe-prone
scope:
  domain: d
deliveryGuarantee: at-least-once
correlatedBy:
  source: event-field
  field: id
state:
  type: object
startsWhen:
  - boundedContext: bc
    aggregate: agg
    event: started
endsWhen:
  - name: done
    condition: completed
reactions:
  - when:
      boundedContext: bc
      aggregate: agg
      event: started
    rule: noop
`
		diags := runRule(t, rule, buildModelTolerantOfSchemaErrors(t, yaml))
		require.Len(t, diags, 1)
		assert.Contains(t, diags[0].Message, "process-manager")
	})
}
