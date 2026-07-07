package rules_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const minimalParentsTwoActors = `apiVersion: schema.esdm.io/core/v1
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
kind: actor
name: admin
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
  - agg-done
actors:
  - user
---
apiVersion: schema.esdm.io/core/v1
kind: event
name: agg-done
scope:
  domain: d
  boundedContext: bc
  aggregate: agg
data:
  type: object
`

func TestScenarioActorNotPermitted(t *testing.T) {
	rule := findCatalogRule(t, "esdm/gwt/scenario-actor-not-permitted")

	t.Run("does not throw when the actor is in the command's actors list", func(t *testing.T) {
		yaml := minimalParentsTwoActors + `---
apiVersion: schema.esdm.io/given-when-then/v1
kind: feature
name: ok
scope:
  domain: d
  boundedContext: bc
  aggregate: agg
scenarios:
  - name: trivial
    when:
      command: do-it
      data: {}
      actor: user
    then:
      events: []
`
		assert.Empty(t, runRule(t, rule, buildModel(t, yaml)))
	})

	t.Run("throws when the actor is not in the command's actors list", func(t *testing.T) {
		yaml := minimalParentsTwoActors + `---
apiVersion: schema.esdm.io/given-when-then/v1
kind: feature
name: forbidden
scope:
  domain: d
  boundedContext: bc
  aggregate: agg
scenarios:
  - name: trivial
    when:
      command: do-it
      data: {}
      actor: admin
    then:
      events: []
`
		diags := runRule(t, rule, buildModel(t, yaml))
		require.Len(t, diags, 1)
		assert.Contains(t, diags[0].Message, "admin")
	})
}
