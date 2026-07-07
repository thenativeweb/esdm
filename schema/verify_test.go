package schema_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thenativeweb/esdm/schema"
)

func TestVerify(t *testing.T) {
	t.Run("passes for a freshly written schemas directory", func(t *testing.T) {
		dir := t.TempDir()
		root := filepath.Join(dir, "schemas")
		require.NoError(t, schema.Write(root))

		assert.NoError(t, schema.Verify(root))
	})

	t.Run("rejects an extra file inside schemas", func(t *testing.T) {
		dir := t.TempDir()
		root := filepath.Join(dir, "schemas")
		require.NoError(t, schema.Write(root))
		require.NoError(t, os.WriteFile(filepath.Join(root, "stray.yaml"), []byte("x"), 0o644))

		err := schema.Verify(root)
		require.Error(t, err)

		var ve *schema.VerifyError
		require.True(t, errors.As(err, &ve))
		assert.Equal(t, "stray.yaml", ve.Path)
	})

	t.Run("rejects a missing embedded file", func(t *testing.T) {
		dir := t.TempDir()
		root := filepath.Join(dir, "schemas")
		require.NoError(t, schema.Write(root))
		require.NoError(t, os.Remove(filepath.Join(root, "core", "v1.yaml")))

		err := schema.Verify(root)
		require.Error(t, err)

		var ve *schema.VerifyError
		require.True(t, errors.As(err, &ve))
		assert.Equal(t, "core/v1.yaml", ve.Path)
	})

	t.Run("rejects a hand-edited file", func(t *testing.T) {
		dir := t.TempDir()
		root := filepath.Join(dir, "schemas")
		require.NoError(t, schema.Write(root))
		require.NoError(t, os.WriteFile(filepath.Join(root, "core", "v1.yaml"), []byte("hand-edited"), 0o644))

		err := schema.Verify(root)
		require.Error(t, err)

		var ve *schema.VerifyError
		require.True(t, errors.As(err, &ve))
		assert.Equal(t, "core/v1.yaml", ve.Path)
	})
}
