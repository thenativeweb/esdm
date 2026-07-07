package rules_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDomainWithoutBoundedContext(t *testing.T) {
	rule := findCatalogRule(t, "esdm/modeling/domain-without-bounded-context")

	t.Run("does not throw when the domain hosts at least one bounded context", func(t *testing.T) {
		assert.Empty(t, runRule(t, rule, buildModel(t, minimalParents)))
	})

	t.Run("throws when a domain has no bounded contexts", func(t *testing.T) {
		yaml := `apiVersion: schema.esdm.io/core/v1
kind: domain
name: empty-domain
`
		diags := runRule(t, rule, buildModel(t, yaml))
		require.Len(t, diags, 1)
		assert.Contains(t, diags[0].Message, "empty-domain")
	})
}
