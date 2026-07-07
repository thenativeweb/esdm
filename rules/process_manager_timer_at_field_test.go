package rules_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProcessManagerTimerAtField(t *testing.T) {
	rule := findCatalogRule(t, "esdm/structure/process-manager-timer-at-field")

	t.Run("does not throw when an absolute timer's `at` field exists in state.properties", func(t *testing.T) {
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
name: tracker
scope:
  domain: d
state:
  type: object
  properties:
    expires-at:
      type: string
startsWhen:
  - boundedContext: bc
    aggregate: agg
    event: started
endsWhen:
  - name: done
    condition: timed out
timers:
  - name: deadline
    at: expires-at
reactions:
  - when:
      timer: deadline
    rule: timeout
`
		assert.Empty(t, runRule(t, rule, buildModel(t, yaml)))
	})

	t.Run("throws when an absolute timer references a field missing from state.properties", func(t *testing.T) {
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
name: tracker
scope:
  domain: d
state:
  type: object
  properties:
    other-field:
      type: string
startsWhen:
  - boundedContext: bc
    aggregate: agg
    event: started
endsWhen:
  - name: done
    condition: timed out
timers:
  - name: deadline
    at: expires-at
reactions:
  - when:
      timer: deadline
    rule: timeout
`
		diags := runRule(t, rule, buildModel(t, yaml))
		require.Len(t, diags, 1)
		assert.Contains(t, diags[0].Message, "tracker")
		assert.Contains(t, diags[0].Message, "deadline")
		assert.Contains(t, diags[0].Message, "expires-at")
	})

	t.Run("does not throw for relative timers (after-shape)", func(t *testing.T) {
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
name: tracker
scope:
  domain: d
state:
  type: object
startsWhen:
  - boundedContext: bc
    aggregate: agg
    event: started
endsWhen:
  - name: done
    condition: timed out
timers:
  - name: cool-off
    after:
      value: 30
      unit: minutes
reactions:
  - when:
      timer: cool-off
    rule: cool off
`
		assert.Empty(t, runRule(t, rule, buildModel(t, yaml)))
	})
}
