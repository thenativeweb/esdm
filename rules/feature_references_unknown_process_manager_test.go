package rules_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const processManagerParents = `apiVersion: schema.esdm.io/core/v1
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
kind: process-manager
name: order-pm
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
    event: agg-done
endsWhen:
  - name: done
    condition: state.done is true
timers:
  - name: deadline
    after:
      value: 30
      unit: minutes
reactions:
  - when:
      timer: deadline
    rule: cancel the order
    emits:
      - boundedContext: bc
        aggregate: agg
        command: do-it
`

func TestFeatureReferencesUnknownProcessManager(t *testing.T) {
	rule := findCatalogRule(t, "esdm/gwt/feature-references-unknown-process-manager")

	t.Run("does not throw when the feature scope points at an existing process-manager", func(t *testing.T) {
		yaml := processManagerParents + `---
apiVersion: schema.esdm.io/given-when-then/v1
kind: feature
name: ok-feature
scope:
  domain: d
  processManager: order-pm
scenarios:
  - name: trivial
    when:
      timer: deadline
    then:
      ended: true
`
		assert.Empty(t, runRule(t, rule, buildModel(t, yaml)))
	})

	t.Run("throws when the feature scope points at a process-manager that does not exist", func(t *testing.T) {
		yaml := processManagerParents + `---
apiVersion: schema.esdm.io/given-when-then/v1
kind: feature
name: dangling-feature
scope:
  domain: d
  processManager: phantom-pm
scenarios:
  - name: trivial
    when:
      timer: deadline
    then:
      ended: true
`
		diags := runRule(t, rule, buildModel(t, yaml))
		require.Len(t, diags, 1)
		assert.Contains(t, diags[0].Message, "phantom-pm")
	})
}
