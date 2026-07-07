package updateschema_test

import (
	"bytes"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thenativeweb/esdm/cmd/esdm/commands/updateschema"
	"github.com/thenativeweb/esdm/schema"
)

func runUpdateSchemaIn(t *testing.T, dir string) (string, error) {
	t.Helper()

	originalWd, err := os.Getwd()
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = os.Chdir(originalWd)
	})

	require.NoError(t, os.Chdir(dir))

	var buf bytes.Buffer
	updateschema.Command.SetOut(&buf)
	updateschema.Command.SetErr(&buf)
	updateschema.Command.SetArgs([]string{})
	err = updateschema.Command.Execute()
	return buf.String(), err
}

// rewriteRevision replaces the value of the
// `x-esdm-schema-revision` field inside YAML bytes for
// test purposes. It exists only here because no
// production code ever rewrites a revision.
func rewriteRevision(t *testing.T, in []byte, newRevision string) []byte {
	t.Helper()
	revisionPattern := regexp.MustCompile(`(?m)^x-esdm-schema-revision:\s*"[^"]*"`)
	out := revisionPattern.ReplaceAll(in, []byte("x-esdm-schema-revision: \""+newRevision+"\""))
	require.NotEqual(t, string(in), string(out), "test fixture did not change; revision field not found")
	return out
}

// writeSchemas writes a freshly-embedded `schemas/` tree
// at dir/schemas, like add-schema would.
func writeSchemas(t *testing.T, dir string) {
	t.Helper()
	require.NoError(t, schema.Write(filepath.Join(dir, "schemas")))
}

func TestUpdateSchema(t *testing.T) {
	t.Run("refuses to run when schemas/ does not exist", func(t *testing.T) {
		dir := t.TempDir()

		out, err := runUpdateSchemaIn(t, dir)
		require.Error(t, err)
		assert.Contains(t, out+err.Error(), "add-schema")
	})

	t.Run("refuses to run when schemas exists as a file", func(t *testing.T) {
		dir := t.TempDir()
		require.NoError(t, os.WriteFile(filepath.Join(dir, "schemas"), []byte(""), 0o644))

		_, err := runUpdateSchemaIn(t, dir)
		require.Error(t, err)
	})

	t.Run("reports nothing to update when local matches embedded byte-for-byte", func(t *testing.T) {
		dir := t.TempDir()
		writeSchemas(t, dir)

		out, err := runUpdateSchemaIn(t, dir)
		require.NoError(t, err)
		assert.Contains(t, out, "nothing")
	})

	t.Run("wipes and rewrites when an embedded revision is newer than local", func(t *testing.T) {
		dir := t.TempDir()
		writeSchemas(t, dir)

		// Pretend the local copy is on an older revision.
		localCore := filepath.Join(dir, "schemas", "core", "v1.yaml")
		original, err := os.ReadFile(localCore)
		require.NoError(t, err)
		stale := rewriteRevision(t, original, "0.9.0")
		require.NoError(t, os.WriteFile(localCore, stale, 0o644))

		_, err = runUpdateSchemaIn(t, dir)
		require.NoError(t, err)

		// Local should now match embedded again.
		updated, err := os.ReadFile(localCore)
		require.NoError(t, err)
		assert.Equal(t, schema.Core(), updated)
	})

	t.Run("removes stray files inside schemas/ during a rewrite", func(t *testing.T) {
		dir := t.TempDir()
		writeSchemas(t, dir)

		stray := filepath.Join(dir, "schemas", "stray.yaml")
		require.NoError(t, os.WriteFile(stray, []byte("hand-edited"), 0o644))

		// Force a rewrite by lowering one local revision.
		localCore := filepath.Join(dir, "schemas", "core", "v1.yaml")
		original, err := os.ReadFile(localCore)
		require.NoError(t, err)
		stale := rewriteRevision(t, original, "0.9.0")
		require.NoError(t, os.WriteFile(localCore, stale, 0o644))

		_, err = runUpdateSchemaIn(t, dir)
		require.NoError(t, err)

		_, err = os.Stat(stray)
		assert.True(t, os.IsNotExist(err), "stray file should be gone after rewrite")
	})

	t.Run("rewrites when local has stray files even if every known revision matches", func(t *testing.T) {
		dir := t.TempDir()
		writeSchemas(t, dir)

		stray := filepath.Join(dir, "schemas", "stray.yaml")
		require.NoError(t, os.WriteFile(stray, []byte("foreign"), 0o644))

		_, err := runUpdateSchemaIn(t, dir)
		require.NoError(t, err)

		_, err = os.Stat(stray)
		assert.True(t, os.IsNotExist(err))
	})

	t.Run("rejects local revisions strictly higher than embedded", func(t *testing.T) {
		dir := t.TempDir()
		writeSchemas(t, dir)

		localCore := filepath.Join(dir, "schemas", "core", "v1.yaml")
		original, err := os.ReadFile(localCore)
		require.NoError(t, err)
		future := rewriteRevision(t, original, "99.0.0")
		require.NoError(t, os.WriteFile(localCore, future, 0o644))

		out, err := runUpdateSchemaIn(t, dir)
		require.Error(t, err)
		assert.Contains(t, out+err.Error(), "newer")
	})

	t.Run("leaves files outside schemas/ alone", func(t *testing.T) {
		dir := t.TempDir()
		writeSchemas(t, dir)
		require.NoError(t, os.WriteFile(filepath.Join(dir, "model.esdm.yaml"), []byte("kind: domain\n"), 0o644))

		// Force a rewrite.
		localCore := filepath.Join(dir, "schemas", "core", "v1.yaml")
		original, err := os.ReadFile(localCore)
		require.NoError(t, err)
		stale := rewriteRevision(t, original, "0.9.0")
		require.NoError(t, os.WriteFile(localCore, stale, 0o644))

		_, err = runUpdateSchemaIn(t, dir)
		require.NoError(t, err)

		data, err := os.ReadFile(filepath.Join(dir, "model.esdm.yaml"))
		require.NoError(t, err)
		assert.Equal(t, "kind: domain\n", string(data))
	})
}
