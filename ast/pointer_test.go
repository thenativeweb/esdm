package ast_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thenativeweb/esdm/ast"
)

func TestFollowPointer(t *testing.T) {
	t.Run("empty pointer returns the root", func(t *testing.T) {
		n := parse(t, "name: foo\n")
		result := ast.FollowPointer(n, "")
		assert.Equal(t, ast.KindMapping, result.Kind())
	})

	t.Run("descends into a mapping by key", func(t *testing.T) {
		n := parse(t, "name: foo\n")
		v, ok := ast.FollowPointer(n, "/name").Text()
		require.True(t, ok)
		assert.Equal(t, "foo", v)
	})

	t.Run("descends through nested mappings", func(t *testing.T) {
		n := parse(t, "outer:\n  inner: value\n")
		v, ok := ast.FollowPointer(n, "/outer/inner").Text()
		require.True(t, ok)
		assert.Equal(t, "value", v)
	})

	t.Run("descends into a sequence by numeric index", func(t *testing.T) {
		n := parse(t, "items:\n  - a\n  - b\n  - c\n")
		v, ok := ast.FollowPointer(n, "/items/1").Text()
		require.True(t, ok)
		assert.Equal(t, "b", v)
	})

	t.Run("unescapes ~1 to forward slash", func(t *testing.T) {
		n := parse(t, "\"a/b\": value\n")
		v, ok := ast.FollowPointer(n, "/a~1b").Text()
		require.True(t, ok)
		assert.Equal(t, "value", v)
	})

	t.Run("unescapes ~0 to tilde", func(t *testing.T) {
		n := parse(t, "\"a~b\": value\n")
		v, ok := ast.FollowPointer(n, "/a~0b").Text()
		require.True(t, ok)
		assert.Equal(t, "value", v)
	})

	t.Run("unescapes ~01 as tilde-one, not slash", func(t *testing.T) {
		n := parse(t, "\"a~1b\": value\n")
		v, ok := ast.FollowPointer(n, "/a~01b").Text()
		require.True(t, ok)
		assert.Equal(t, "value", v)
	})

	t.Run("returns the last resolved node on missing key", func(t *testing.T) {
		n := parse(t, "outer:\n  inner: value\n")
		result := ast.FollowPointer(n, "/outer/missing")
		assert.Equal(t, ast.KindMapping, result.Kind())

		inner := result.Field("inner")
		v, ok := inner.Text()
		require.True(t, ok)
		assert.Equal(t, "value", v)
	})

	t.Run("returns the last resolved node on out-of-range sequence index", func(t *testing.T) {
		n := parse(t, "items:\n  - a\n")
		result := ast.FollowPointer(n, "/items/5")
		assert.Equal(t, ast.KindSequence, result.Kind())
	})

	t.Run("returns the last resolved node on invalid sequence index", func(t *testing.T) {
		n := parse(t, "items:\n  - a\n")
		result := ast.FollowPointer(n, "/items/not-a-number")
		assert.Equal(t, ast.KindSequence, result.Kind())
	})

	t.Run("returns the current node when descent would require a scalar to be traversed", func(t *testing.T) {
		n := parse(t, "name: foo\n")
		result := ast.FollowPointer(n, "/name/deeper")
		v, ok := result.Text()
		require.True(t, ok)
		assert.Equal(t, "foo", v)
	})

	t.Run("rejects pointer without leading slash", func(t *testing.T) {
		n := parse(t, "name: foo\n")
		result := ast.FollowPointer(n, "name")
		assert.Equal(t, ast.KindMapping, result.Kind())
	})
}
