package rules_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProcessManagerWithoutStartsWhen(t *testing.T) {
	rule := findCatalogRule(t, "esdm/modeling/process-manager-without-starts-when")

	t.Run("does not throw when the process-manager declares startsWhen", func(t *testing.T) {
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
deliveryGuarantee: at-most-once
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
		assert.Empty(t, runRule(t, rule, buildModel(t, yaml)))
	})

	t.Run("throws when the process-manager has no startsWhen field", func(t *testing.T) {
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
deliveryGuarantee: at-most-once
correlatedBy:
  source: event-field
  field: id
state:
  type: object
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
		assert.Contains(t, diags[0].Message, "tracker")
	})
}
