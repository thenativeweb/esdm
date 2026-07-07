package rules_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"github.com/thenativeweb/esdm/ast"
	"github.com/thenativeweb/esdm/rules"
)

func parseYAMLNode(t *testing.T, src string) ast.Node {
	t.Helper()
	var raw yaml.Node
	require.NoError(t, yaml.Unmarshal([]byte(src), &raw))
	return ast.NewNode("test.yaml", &raw)
}

func TestSchemaHasProperty(t *testing.T) {
	t.Run("returns true when the schema declares the property at the top level", func(t *testing.T) {
		n := parseYAMLNode(t, `
type: object
properties:
  id:
    type: string
  amount:
    type: number
`)
		assert.True(t, rules.SchemaHasProperty(n, "id"))
		assert.True(t, rules.SchemaHasProperty(n, "amount"))
	})

	t.Run("returns false when the property is absent", func(t *testing.T) {
		n := parseYAMLNode(t, `
type: object
properties:
  id:
    type: string
`)
		assert.False(t, rules.SchemaHasProperty(n, "name"))
	})

	t.Run("returns false when the schema has no properties block", func(t *testing.T) {
		n := parseYAMLNode(t, `
type: object
`)
		assert.False(t, rules.SchemaHasProperty(n, "id"))
	})

	t.Run("returns false on a missing schema node", func(t *testing.T) {
		assert.False(t, rules.SchemaHasProperty(ast.Node{}, "id"))
	})
}
