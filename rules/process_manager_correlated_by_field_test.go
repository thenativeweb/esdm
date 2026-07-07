package rules_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const pmParents = `apiVersion: schema.esdm.io/core/v1
kind: domain
name: d
---
apiVersion: schema.esdm.io/core/v1
kind: bounded-context
name: bc
scope:
  domain: d
---
apiVersion: schema.esdm.io/core/v1
kind: aggregate
name: agg
scope:
  domain: d
  boundedContext: bc
identifiedBy:
  source: generated
  generator: uuid
state:
  type: object
---
apiVersion: schema.esdm.io/core/v1
kind: actor
name: user
scope:
  domain: d
  boundedContext: bc
type: human
---
apiVersion: schema.esdm.io/core/v1
kind: command
name: do-it
scope:
  domain: d
  boundedContext: bc
  aggregate: agg
data:
  type: object
publishes:
  - started
  - finished
actors:
  - user
`

func TestProcessManagerCorrelatedByField(t *testing.T) {
	rule := findCatalogRule(t, "esdm/structure/process-manager-correlated-by-field")

	t.Run("does not throw when every consumed event declares the correlation field", func(t *testing.T) {
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
  properties:
    correlation-id:
      type: string
---
apiVersion: schema.esdm.io/core/v1
kind: event
name: finished
scope:
  domain: d
  boundedContext: bc
  aggregate: agg
data:
  type: object
  properties:
    correlation-id:
      type: string
---
apiVersion: schema.esdm.io/core/v1
kind: process-manager
name: tracker
scope:
  domain: d
correlatedBy:
  source: event-field
  field: correlation-id
state:
  type: object
startsWhen:
  - boundedContext: bc
    aggregate: agg
    event: started
endsWhen:
  - name: done
    condition: state.completed is true
reactions:
  - when:
      boundedContext: bc
      aggregate: agg
      event: finished
    rule: mark complete
`
		assert.Empty(t, runRule(t, rule, buildModel(t, yaml)))
	})

	t.Run("throws when a startsWhen event lacks the correlation field", func(t *testing.T) {
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
kind: event
name: finished
scope:
  domain: d
  boundedContext: bc
  aggregate: agg
data:
  type: object
  properties:
    correlation-id:
      type: string
---
apiVersion: schema.esdm.io/core/v1
kind: process-manager
name: tracker
scope:
  domain: d
correlatedBy:
  source: event-field
  field: correlation-id
state:
  type: object
startsWhen:
  - boundedContext: bc
    aggregate: agg
    event: started
endsWhen:
  - name: done
    condition: state.completed is true
reactions:
  - when:
      boundedContext: bc
      aggregate: agg
      event: finished
    rule: mark complete
`
		diags := runRule(t, rule, buildModel(t, yaml))
		require.Len(t, diags, 1)
		assert.Contains(t, diags[0].Message, "tracker")
		assert.Contains(t, diags[0].Message, "started")
		assert.Contains(t, diags[0].Message, "correlation-id")
	})

	t.Run("ignores reactions whose when is a timer rather than an event", func(t *testing.T) {
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
  properties:
    correlation-id:
      type: string
---
apiVersion: schema.esdm.io/core/v1
kind: event
name: finished
scope:
  domain: d
  boundedContext: bc
  aggregate: agg
data:
  type: object
  properties:
    correlation-id:
      type: string
---
apiVersion: schema.esdm.io/core/v1
kind: process-manager
name: tracker
scope:
  domain: d
correlatedBy:
  source: event-field
  field: correlation-id
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
    condition: state.completed is true
timers:
  - name: deadline
    at: expires-at
reactions:
  - when:
      boundedContext: bc
      aggregate: agg
      event: finished
    rule: mark complete
  - when:
      timer: deadline
    rule: timeout
`
		assert.Empty(t, runRule(t, rule, buildModel(t, yaml)))
	})
}
