package rules_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProcessManagerWithoutEventReactions(t *testing.T) {
	rule := findCatalogRule(t, "esdm/modeling/process-manager-without-event-reactions")

	t.Run("does not throw when the process manager has at least one event reaction", func(t *testing.T) {
		yaml := minimalParents + `---
apiVersion: schema.esdm.io/core/v1
kind: process-manager
name: saga
scope:
  domain: d
  boundedContext: bc
deliveryGuarantee: at-most-once
correlatedBy:
  name: id
  field: id
state:
  type: object
  properties:
    done:
      type: boolean
startsWhen:
  - boundedContext: bc
    aggregate: agg
    event: agg-done
endsWhen:
  - name: done
    condition: state.done is true
reactions:
  - when:
      boundedContext: bc
      aggregate: agg
      event: agg-done
    rule: mark it done
`
		assert.Empty(t, runRule(t, rule, buildModel(t, yaml)))
	})

	t.Run("throws when all reactions are timer-based", func(t *testing.T) {
		yaml := minimalParents + `---
apiVersion: schema.esdm.io/core/v1
kind: process-manager
name: saga
scope:
  domain: d
  boundedContext: bc
deliveryGuarantee: at-most-once
correlatedBy:
  name: id
  field: id
state:
  type: object
  properties:
    done:
      type: boolean
startsWhen:
  - boundedContext: bc
    aggregate: agg
    event: agg-done
endsWhen:
  - name: done
    condition: state.done is true
timers:
  - name: expiry
    after:
      value: 30
      unit: minutes
reactions:
  - when:
      timer: expiry
    rule: expire the saga
`
		diags := runRule(t, rule, buildModel(t, yaml))
		require.Len(t, diags, 1)
		assert.Contains(t, diags[0].Message, "saga")
	})
}
