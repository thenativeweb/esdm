package schema_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thenativeweb/esdm/schema"
)

func TestWrite(t *testing.T) {
	t.Run("creates the schemas directory and writes every embedded file", func(t *testing.T) {
		root := filepath.Join(t.TempDir(), "schemas")

		require.NoError(t, schema.Write(root))

		for _, file := range schema.Files() {
			data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(file.Path)))
			require.NoError(t, err)
			assert.Equal(t, file.Bytes, data)
		}
	})

	t.Run("works when the target directory already exists and is empty", func(t *testing.T) {
		root := filepath.Join(t.TempDir(), "schemas")
		require.NoError(t, os.MkdirAll(root, 0o755))

		require.NoError(t, schema.Write(root))

		_, err := os.Stat(filepath.Join(root, "core", "v1.yaml"))
		assert.NoError(t, err)
	})

	t.Run("overwrites existing files at the same path", func(t *testing.T) {
		root := filepath.Join(t.TempDir(), "schemas")
		corePath := filepath.Join(root, "core", "v1.yaml")
		require.NoError(t, os.MkdirAll(filepath.Dir(corePath), 0o755))
		require.NoError(t, os.WriteFile(corePath, []byte("stale"), 0o644))

		require.NoError(t, schema.Write(root))

		data, err := os.ReadFile(corePath)
		require.NoError(t, err)
		assert.Equal(t, schema.Core(), data)
	})
}
