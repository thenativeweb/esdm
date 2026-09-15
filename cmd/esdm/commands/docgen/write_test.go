package docgen_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thenativeweb/esdm/cmd/esdm/commands/docgen"
)

func samplePages() []docgen.Page {
	return []docgen.Page{
		{Path: "README.md", Content: "root"},
		{Path: "library/README.md", Content: "domain"},
		{Path: "library/cataloging/book/acquired.md", Content: "event"},
	}
}

func TestWrite(t *testing.T) {
	t.Run("creates the tree and writes each page", func(t *testing.T) {
		out := filepath.Join(t.TempDir(), "docs")

		err := docgen.Write(samplePages(), out, false)
		require.NoError(t, err)

		content, err := os.ReadFile(filepath.Join(out, "library", "cataloging", "book", "acquired.md"))
		require.NoError(t, err)
		assert.Equal(t, "event", string(content))

		root, err := os.ReadFile(filepath.Join(out, "README.md"))
		require.NoError(t, err)
		assert.Equal(t, "root", string(root))
	})

	t.Run("refuses to write into a non-empty directory without force", func(t *testing.T) {
		out := t.TempDir()
		stray := filepath.Join(out, "stray.txt")
		require.NoError(t, os.WriteFile(stray, []byte("keep"), 0o644))

		err := docgen.Write(samplePages(), out, false)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "--force")

		_, statErr := os.Stat(stray)
		assert.NoError(t, statErr)
	})

	t.Run("clears the directory before writing when force is set", func(t *testing.T) {
		out := t.TempDir()
		stray := filepath.Join(out, "stray.txt")
		require.NoError(t, os.WriteFile(stray, []byte("remove"), 0o644))

		err := docgen.Write(samplePages(), out, true)
		require.NoError(t, err)

		_, statErr := os.Stat(stray)
		assert.True(t, os.IsNotExist(statErr))

		content, err := os.ReadFile(filepath.Join(out, "README.md"))
		require.NoError(t, err)
		assert.Equal(t, "root", string(content))
	})

	t.Run("writes into an existing empty directory", func(t *testing.T) {
		out := t.TempDir()

		err := docgen.Write(samplePages(), out, false)
		require.NoError(t, err)

		content, err := os.ReadFile(filepath.Join(out, "library", "README.md"))
		require.NoError(t, err)
		assert.Equal(t, "domain", string(content))
	})
}
