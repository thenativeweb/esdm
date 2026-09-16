package rules_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const ambiguousNameParents = `apiVersion: schema.esdm.io/core/v1
kind: domain
name: d
---
apiVersion: schema.esdm.io/core/v1
kind: bounded-context
name: lending
scope:
  domain: d
---
`

const aggregateLoan = `apiVersion: schema.esdm.io/core/v1
kind: aggregate
name: loan
scope:
  domain: d
  boundedContext: lending
identifiedBy:
  source: generated
  generator: uuid
state:
  type: object
---
`

const dcbLoan = `apiVersion: schema.esdm.io/core/v1
kind: dynamic-consistency-boundary
name: loan
scope:
  domain: d
  boundedContext: lending
identifiedBy:
  - name: id
    source: static
    value: solo
consults:
  - boundedContext: lending
    aggregate: loan
    event: granted
    criteria: relevant
---
apiVersion: schema.esdm.io/core/v1
kind: event
name: granted
scope:
  domain: d
  boundedContext: lending
  aggregate: loan
data:
  type: object
---
`

const readModelLoan = `apiVersion: schema.esdm.io/core/v1
kind: read-model
name: loan
scope:
  domain: d
  boundedContext: lending
projections: []
schema:
  type: object
---
`

const entityMoney = `apiVersion: schema.esdm.io/core/v1
kind: entity
name: money
scope:
  domain: d
  boundedContext: lending
schema:
  type: object
identifiedBy:
  source: static
  value: single
---
`

const valueObjectMoney = `apiVersion: schema.esdm.io/core/v1
kind: value-object
name: money
scope:
  domain: d
  boundedContext: lending
schema:
  type: object
---
`

const domainServiceMoney = `apiVersion: schema.esdm.io/core/v1
kind: domain-service
name: money
scope:
  domain: d
  boundedContext: lending
functions:
  - name: convert
    arguments:
      type: object
    returns:
      type: object
---
`

const actorMoney = `apiVersion: schema.esdm.io/core/v1
kind: actor
name: money
scope:
  domain: d
  boundedContext: lending
type: system
---
`

func TestAmbiguousName(t *testing.T) {
	rule := findCatalogRule(t, "esdm/structure/ambiguous-name")

	t.Run("does not throw when an aggregate and a read model share a name", func(t *testing.T) {
		assert.Empty(t, runRule(t, rule, buildModel(t, ambiguousNameParents+aggregateLoan+readModelLoan)))
	})

	t.Run("does not throw when an entity and an actor share a name", func(t *testing.T) {
		assert.Empty(t, runRule(t, rule, buildModel(t, ambiguousNameParents+entityMoney+actorMoney)))
	})

	t.Run("throws when an aggregate and a dynamic consistency boundary share a name", func(t *testing.T) {
		diags := runRule(t, rule, buildModel(t, ambiguousNameParents+aggregateLoan+dcbLoan))
		require.Len(t, diags, 1)
		assert.Equal(t, `"loan" names both an aggregate and a dynamic-consistency-boundary in bounded-context "lending"`, diags[0].Message)
		require.Len(t, diags[0].Related, 1)
		assert.Equal(t, `aggregate "loan" defined here`, diags[0].Related[0].Message)
		// The diagnostic sits at the later element; the note at the first.
		assert.Greater(t, diags[0].Location.Line, diags[0].Related[0].Location.Line)
	})

	t.Run("throws once per additional element when three domain types share a name", func(t *testing.T) {
		diags := runRule(t, rule, buildModel(t, ambiguousNameParents+entityMoney+valueObjectMoney+domainServiceMoney))
		require.Len(t, diags, 2)
		assert.Equal(t, `"money" names both an entity and a value-object in bounded-context "lending"`, diags[0].Message)
		assert.Equal(t, `"money" names both an entity and a domain-service in bounded-context "lending"`, diags[1].Message)
	})

	t.Run("does not throw for the same name in different bounded contexts", func(t *testing.T) {
		other := `apiVersion: schema.esdm.io/core/v1
kind: bounded-context
name: billing
scope:
  domain: d
---
apiVersion: schema.esdm.io/core/v1
kind: value-object
name: money
scope:
  domain: d
  boundedContext: billing
schema:
  type: object
---
`
		assert.Empty(t, runRule(t, rule, buildModel(t, ambiguousNameParents+entityMoney+other)))
	})
}
