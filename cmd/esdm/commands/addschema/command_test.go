package addschema_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thenativeweb/esdm/cmd/esdm/commands/addschema"
	"github.com/thenativeweb/esdm/schema"
)

func runAddSchemaIn(t *testing.T, dir string) (string, error) {
	t.Helper()

	originalWorkingDirectory, err := os.Getwd()
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = os.Chdir(originalWorkingDirectory)
	})

	require.NoError(t, os.Chdir(dir))

	var output bytes.Buffer
	addschema.Command.SetOut(&output)
	addschema.Command.SetErr(&output)
	addschema.Command.SetArgs([]string{})
	err = addschema.Command.Execute()
	return output.String(), err
}

func TestAddSchema(t *testing.T) {
	t.Run("writes the embedded schema set into a fresh schemas/ directory at the working directory", func(t *testing.T) {
		dir := t.TempDir()

		_, err := runAddSchemaIn(t, dir)
		require.NoError(t, err)

		for _, file := range schema.Files() {
			data, readErr := os.ReadFile(filepath.Join(dir, "schemas", filepath.FromSlash(file.Path)))
			require.NoError(t, readErr)
			assert.Equal(t, file.Bytes, data)
		}
	})

	t.Run("leaves unrelated files in the working directory untouched", func(t *testing.T) {
		dir := t.TempDir()
		require.NoError(t, os.WriteFile(filepath.Join(dir, "model.esdm.yaml"), []byte("kind: domain\n"), 0o644))

		_, err := runAddSchemaIn(t, dir)
		require.NoError(t, err)

		data, err := os.ReadFile(filepath.Join(dir, "model.esdm.yaml"))
		require.NoError(t, err)
		assert.Equal(t, "kind: domain\n", string(data))
	})

	t.Run("refuses to run when schemas/ already exists", func(t *testing.T) {
		dir := t.TempDir()
		require.NoError(t, os.Mkdir(filepath.Join(dir, "schemas"), 0o755))

		out, err := runAddSchemaIn(t, dir)
		require.Error(t, err)
		assert.Contains(t, out+err.Error(), "update-schema")
	})

	t.Run("refuses to run when schemas exists as a file", func(t *testing.T) {
		dir := t.TempDir()
		require.NoError(t, os.WriteFile(filepath.Join(dir, "schemas"), []byte(""), 0o644))

		_, err := runAddSchemaIn(t, dir)
		require.Error(t, err)
	})
}
