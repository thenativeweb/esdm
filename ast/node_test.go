package ast_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"github.com/thenativeweb/esdm/ast"
)

func parse(t *testing.T, src string) ast.Node {
	t.Helper()

	var doc yaml.Node
	require.NoError(t, yaml.Unmarshal([]byte(src), &doc))
	return ast.NewNode("test.esdm.yaml", &doc)
}

func TestNewNode(t *testing.T) {
	t.Run("returns a zero Node for nil input", func(t *testing.T) {
		n := ast.NewNode("f", nil)
		assert.False(t, n.Exists())
	})

	t.Run("unwraps a document node to its content", func(t *testing.T) {
		n := parse(t, "name: order-placed\n")
		assert.Equal(t, ast.KindMapping, n.Kind())
	})
}

func TestNodeNavigation(t *testing.T) {
	t.Run("Field returns the value for a known key", func(t *testing.T) {
		n := parse(t, "name: order-placed\n")
		name := n.Field("name")
		v, ok := name.Text()
		require.True(t, ok)
		assert.Equal(t, "order-placed", v)
	})

	t.Run("Field returns a zero Node for a missing key", func(t *testing.T) {
		n := parse(t, "name: order-placed\n")
		assert.False(t, n.Field("missing").Exists())
	})

	t.Run("Field on a non-mapping returns a zero Node", func(t *testing.T) {
		n := parse(t, "- a\n- b\n")
		assert.False(t, n.Field("name").Exists())
	})

	t.Run("chained Field calls are safe on missing intermediate paths", func(t *testing.T) {
		n := parse(t, "name: order-placed\n")
		assert.False(t, n.Field("missing").Field("deep").Field("deeper").Exists())
	})

	t.Run("At returns the i-th element of a sequence", func(t *testing.T) {
		n := parse(t, "- a\n- b\n- c\n")
		v, ok := n.At(1).Text()
		require.True(t, ok)
		assert.Equal(t, "b", v)
	})

	t.Run("At returns a zero Node for negative or out-of-range indexes", func(t *testing.T) {
		n := parse(t, "- a\n- b\n")
		assert.False(t, n.At(-1).Exists())
		assert.False(t, n.At(5).Exists())
	})

	t.Run("At on a non-sequence returns a zero Node", func(t *testing.T) {
		n := parse(t, "name: order-placed\n")
		assert.False(t, n.At(0).Exists())
	})

	t.Run("Seq returns all elements of a sequence", func(t *testing.T) {
		n := parse(t, "- a\n- b\n- c\n")
		elements := n.Seq()
		require.Len(t, elements, 3)
		v, _ := elements[2].Text()
		assert.Equal(t, "c", v)
	})

	t.Run("Seq on a non-sequence returns nil", func(t *testing.T) {
		n := parse(t, "name: foo\n")
		assert.Nil(t, n.Seq())
	})

	t.Run("Entries returns key/value pairs in file order", func(t *testing.T) {
		n := parse(t, "name: foo\nkind: bar\n")
		entries := n.Entries()
		require.Len(t, entries, 2)

		k0, _ := entries[0].Key.Text()
		v0, _ := entries[0].Value.Text()
		assert.Equal(t, "name", k0)
		assert.Equal(t, "foo", v0)

		k1, _ := entries[1].Key.Text()
		v1, _ := entries[1].Value.Text()
		assert.Equal(t, "kind", k1)
		assert.Equal(t, "bar", v1)
	})

	t.Run("Entries on a non-mapping returns nil", func(t *testing.T) {
		n := parse(t, "- a\n")
		assert.Nil(t, n.Entries())
	})

	t.Run("HasField reports presence of a key", func(t *testing.T) {
		n := parse(t, "name: foo\n")
		assert.True(t, n.HasField("name"))
		assert.False(t, n.HasField("missing"))
	})
}

func TestNodeLeafValues(t *testing.T) {
	t.Run("Text returns scalar string values", func(t *testing.T) {
		n := parse(t, "name: order-placed\n")
		v, ok := n.Field("name").Text()
		require.True(t, ok)
		assert.Equal(t, "order-placed", v)
	})

	t.Run("Text rejects non-string scalars", func(t *testing.T) {
		n := parse(t, "count: 3\n")
		_, ok := n.Field("count").Text()
		assert.False(t, ok)
	})

	t.Run("Int returns integer scalar values", func(t *testing.T) {
		n := parse(t, "count: 42\n")
		v, ok := n.Field("count").Int()
		require.True(t, ok)
		assert.Equal(t, int64(42), v)
	})

	t.Run("Int rejects non-integer scalars", func(t *testing.T) {
		n := parse(t, "name: hello\n")
		_, ok := n.Field("name").Int()
		assert.False(t, ok)
	})

	t.Run("Bool returns boolean scalar values", func(t *testing.T) {
		n := parse(t, "flag: true\n")
		v, ok := n.Field("flag").Bool()
		require.True(t, ok)
		assert.True(t, v)
	})

	t.Run("Bool rejects non-boolean scalars", func(t *testing.T) {
		n := parse(t, "count: 3\n")
		_, ok := n.Field("count").Bool()
		assert.False(t, ok)
	})
}

func TestNodeLocation(t *testing.T) {
	t.Run("returns file, line, and column for a known position", func(t *testing.T) {
		n := parse(t, "name: order-placed\n")
		loc := n.Field("name").Location()
		assert.Equal(t, "test.esdm.yaml", loc.File)
		assert.Equal(t, 1, loc.Line)
		assert.Positive(t, loc.Column)
	})

	t.Run("returns a zero Location for a missing node", func(t *testing.T) {
		n := parse(t, "name: x\n")
		assert.True(t, n.Field("missing").Location().IsZero())
	})
}
