package rules_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDomainServiceWithoutFunctions(t *testing.T) {
	rule := findCatalogRule(t, "esdm/modeling/domain-service-without-functions")

	t.Run("does not throw when the domain-service declares at least one function", func(t *testing.T) {
		yaml := minimalParents + `---
apiVersion: schema.esdm.io/core/v1
kind: domain-service
name: pricing
scope:
  domain: d
  boundedContext: bc
functions:
  - name: compute-total
    arguments:
      type: object
    returns:
      type: object
`
		assert.Empty(t, runRule(t, rule, buildModel(t, yaml)))
	})

	t.Run("throws when the domain-service has no functions field", func(t *testing.T) {
		yaml := minimalParents + `---
apiVersion: schema.esdm.io/core/v1
kind: domain-service
name: empty
scope:
  domain: d
  boundedContext: bc
`
		diags := runRule(t, rule, buildModelTolerantOfSchemaErrors(t, yaml))
		require.Len(t, diags, 1)
		assert.Contains(t, diags[0].Message, "empty")
	})
}
