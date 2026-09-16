package rules_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEventWithoutPublisher(t *testing.T) {
	rule := findCatalogRule(t, "esdm/modeling/event-without-publisher")

	t.Run("does not throw when a command publishes the event", func(t *testing.T) {
		assert.Empty(t, runRule(t, rule, buildModel(t, minimalParents)))
	})

	t.Run("throws when no command publishes the event", func(t *testing.T) {
		yaml := `apiVersion: schema.esdm.io/core/v1
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
kind: event
name: orphan
scope:
  domain: d
  boundedContext: bc
  aggregate: agg
data:
  type: object
`
		diags := runRule(t, rule, buildModel(t, yaml))
		require.Len(t, diags, 1)
		assert.Contains(t, diags[0].Message, "orphan")
	})

	// A DCB-bound command may publish an aggregate-owned
	// event of its bounded context; the resolver accepts
	// that reference, so the rule must count it as the
	// event's publisher.
	t.Run("does not throw when only a DCB-bound command publishes an aggregate-owned event", func(t *testing.T) {
		assert.Empty(t, runRule(t, rule, buildModel(t, dcbPublishesAggregateEventYAML)))
	})
}

const dcbPublishesAggregateEventYAML = `apiVersion: schema.esdm.io/core/v1
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
name: order
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
kind: event
name: checked
scope:
  domain: d
  boundedContext: bc
  aggregate: order
data:
  type: object
---
apiVersion: schema.esdm.io/core/v1
kind: dynamic-consistency-boundary
name: capacity
scope:
  domain: d
  boundedContext: bc
identifiedBy:
  - name: id
    source: static
    value: solo
consults:
  - boundedContext: bc
    aggregate: order
    event: checked
    criteria: relevant
---
apiVersion: schema.esdm.io/core/v1
kind: command
name: reserve
scope:
  domain: d
  boundedContext: bc
  dynamicConsistencyBoundary: capacity
data:
  type: object
publishes:
  - reserved
  - checked
---
apiVersion: schema.esdm.io/core/v1
kind: event
name: reserved
scope:
  domain: d
  boundedContext: bc
data:
  type: object
`
