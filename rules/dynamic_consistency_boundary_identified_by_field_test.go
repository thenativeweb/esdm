package rules_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const dcbParents = `apiVersion: schema.esdm.io/core/v1
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
kind: event
name: enrolled
scope:
  domain: d
  boundedContext: bc
  aggregate: agg
data:
  type: object
`

func TestDynamicConsistencyBoundaryIdentifiedByField(t *testing.T) {
	rule := findCatalogRule(t, "esdm/structure/dynamic-consistency-boundary-identified-by-field")

	t.Run("does not throw when triggering command's data declares the referenced field", func(t *testing.T) {
		yaml := dcbParents + `---
apiVersion: schema.esdm.io/core/v1
kind: dynamic-consistency-boundary
name: capacity
scope:
  domain: d
  boundedContext: bc
identifiedBy:
  - name: studentId
    source: command-payload
    field: student-id
consults:
  - boundedContext: bc
    aggregate: agg
    event: enrolled
    criteria: relevant
---
apiVersion: schema.esdm.io/core/v1
kind: command
name: enroll
scope:
  domain: d
  boundedContext: bc
  dynamicConsistencyBoundary: capacity
data:
  type: object
  properties:
    student-id:
      type: string
publishes:
  - enrolled
actors:
  - user
`
		assert.Empty(t, runRule(t, rule, buildModel(t, yaml)))
	})

	t.Run("throws when triggering command's data is missing the referenced field", func(t *testing.T) {
		yaml := dcbParents + `---
apiVersion: schema.esdm.io/core/v1
kind: dynamic-consistency-boundary
name: capacity
scope:
  domain: d
  boundedContext: bc
identifiedBy:
  - name: studentId
    source: command-payload
    field: student-id
consults:
  - boundedContext: bc
    aggregate: agg
    event: enrolled
    criteria: relevant
---
apiVersion: schema.esdm.io/core/v1
kind: command
name: enroll
scope:
  domain: d
  boundedContext: bc
  dynamicConsistencyBoundary: capacity
data:
  type: object
  properties:
    other-field:
      type: string
publishes:
  - enrolled
actors:
  - user
`
		diags := runRule(t, rule, buildModel(t, yaml))
		require.Len(t, diags, 1)
		assert.Contains(t, diags[0].Message, "capacity")
		assert.Contains(t, diags[0].Message, "enroll")
		assert.Contains(t, diags[0].Message, "student-id")
	})

	t.Run("does not throw for static / generated identifiedBy entries", func(t *testing.T) {
		yaml := dcbParents + `---
apiVersion: schema.esdm.io/core/v1
kind: dynamic-consistency-boundary
name: singleton
scope:
  domain: d
  boundedContext: bc
identifiedBy:
  - name: theOnlyOne
    source: static
    value: global
consults:
  - boundedContext: bc
    aggregate: agg
    event: enrolled
    criteria: relevant
`
		assert.Empty(t, runRule(t, rule, buildModel(t, yaml)))
	})
}
