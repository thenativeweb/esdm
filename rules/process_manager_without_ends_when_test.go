package rules_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProcessManagerWithoutEndsWhen(t *testing.T) {
	rule := findCatalogRule(t, "esdm/modeling/process-manager-without-ends-when")

	t.Run("does not throw when the process-manager declares endsWhen", func(t *testing.T) {
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

	t.Run("throws when the process-manager has no endsWhen field", func(t *testing.T) {
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
name: forever
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
reactions:
  - when:
      boundedContext: bc
      aggregate: agg
      event: started
    rule: noop
`
		diags := runRule(t, rule, buildModelTolerantOfSchemaErrors(t, yaml))
		require.Len(t, diags, 1)
		assert.Contains(t, diags[0].Message, "forever")
	})
}
