package modelpath_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thenativeweb/esdm/modelpath"
)

func TestParsePath(t *testing.T) {
	t.Run("returns an empty path for the empty string", func(t *testing.T) {
		path, err := modelpath.ParsePath("")
		require.NoError(t, err)
		assert.Empty(t, path.Segments)
	})

	t.Run("splits a single-segment path", func(t *testing.T) {
		path, err := modelpath.ParsePath("sample")
		require.NoError(t, err)
		assert.Equal(t, []string{"sample"}, path.Segments)
	})

	t.Run("splits a multi-segment path", func(t *testing.T) {
		path, err := modelpath.ParsePath("sample/context-one/widget")
		require.NoError(t, err)
		assert.Equal(t, []string{"sample", "context-one", "widget"}, path.Segments)
	})

	t.Run("tolerates a trailing slash", func(t *testing.T) {
		path, err := modelpath.ParsePath("sample/context-one/")
		require.NoError(t, err)
		assert.Equal(t, []string{"sample", "context-one"}, path.Segments)
	})

	t.Run("rejects a leading slash", func(t *testing.T) {
		_, err := modelpath.ParsePath("/sample")
		assert.Error(t, err)
	})

	t.Run("rejects empty segments", func(t *testing.T) {
		_, err := modelpath.ParsePath("sample//widget")
		assert.Error(t, err)
	})
}
