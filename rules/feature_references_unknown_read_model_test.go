package rules_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const readModelParents = `apiVersion: schema.esdm.io/core/v1
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
name: agg-done
scope:
  domain: d
  boundedContext: bc
  aggregate: agg
data:
  type: object
---
apiVersion: schema.esdm.io/core/v1
kind: read-model
name: rm
scope:
  domain: d
  boundedContext: bc
projections:
  - boundedContext: bc
    aggregate: agg
    event: agg-done
    rule: append entry
schema:
  type: object
---
apiVersion: schema.esdm.io/core/v1
kind: query
name: q
scope:
  domain: d
  boundedContext: bc
readModel: rm
result:
  type: object
`

func TestFeatureReferencesUnknownReadModel(t *testing.T) {
	rule := findCatalogRule(t, "esdm/gwt/feature-references-unknown-read-model")

	t.Run("does not throw when the feature scope points at an existing read-model", func(t *testing.T) {
		yaml := readModelParents + `---
apiVersion: schema.esdm.io/given-when-then/v1
kind: feature
name: ok-feature
scope:
  domain: d
  boundedContext: bc
  readModel: rm
scenarios:
  - name: trivial
    when:
      query: q
      parameters: {}
    then:
      result: []
`
		assert.Empty(t, runRule(t, rule, buildModel(t, yaml)))
	})

	t.Run("throws when the feature scope points at a read-model that does not exist", func(t *testing.T) {
		yaml := readModelParents + `---
apiVersion: schema.esdm.io/given-when-then/v1
kind: feature
name: dangling-feature
scope:
  domain: d
  boundedContext: bc
  readModel: phantom-rm
scenarios:
  - name: trivial
    when:
      query: q
      parameters: {}
    then:
      result: []
`
		diags := runRule(t, rule, buildModel(t, yaml))
		require.Len(t, diags, 1)
		assert.Contains(t, diags[0].Message, "phantom-rm")
	})
}
