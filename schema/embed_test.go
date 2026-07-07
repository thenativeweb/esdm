package schema_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thenativeweb/esdm/schema"
)

func TestCore(t *testing.T) {
	t.Run("returns non-empty YAML bytes that look like the ESDM core schema", func(t *testing.T) {
		data := schema.Core()
		require.NotEmpty(t, data)

		assert.Contains(t, string(data), "schema.esdm.io/core/v1")
	})

	t.Run("returns a defensive copy so callers cannot corrupt the embedded bytes", func(t *testing.T) {
		first := schema.Core()
		first[0] = 0

		second := schema.Core()
		assert.NotEqual(t, byte(0), second[0])
	})
}

func TestExtensions(t *testing.T) {
	t.Run("returns all embedded extension schemas", func(t *testing.T) {
		extensions, err := schema.Extensions()
		require.NoError(t, err)
		require.NotEmpty(t, extensions)

		var names []string
		for _, extension := range extensions {
			names = append(names, extension.Name)
			assert.NotEmpty(t, extension.Bytes)
		}

		assert.Contains(t, names, "domain-storytelling")
	})

	t.Run("returns extensions in deterministic order", func(t *testing.T) {
		first, err := schema.Extensions()
		require.NoError(t, err)

		second, err := schema.Extensions()
		require.NoError(t, err)

		require.Equal(t, len(first), len(second))
		for i := range first {
			assert.Equal(t, first[i].Name, second[i].Name)
		}
	})
}

func TestFiles(t *testing.T) {
	t.Run("returns every embedded schema file with project-relative path", func(t *testing.T) {
		files := schema.Files()
		require.NotEmpty(t, files)

		paths := make([]string, 0, len(files))
		for _, file := range files {
			paths = append(paths, file.Path)
			assert.NotEmpty(t, file.Bytes)
		}

		assert.Contains(t, paths, "core/v1.yaml")
		assert.Contains(t, paths, "domain-storytelling/v1.yaml")
	})

	t.Run("returns files in deterministic, sorted order", func(t *testing.T) {
		first := schema.Files()
		second := schema.Files()

		require.Equal(t, len(first), len(second))
		for i := range first {
			assert.Equal(t, first[i].Path, second[i].Path)
		}
		for i := 1; i < len(first); i++ {
			assert.Less(t, first[i-1].Path, first[i].Path)
		}
	})

	t.Run("returns defensive copies of bytes", func(t *testing.T) {
		first := schema.Files()
		require.NotEmpty(t, first)
		first[0].Bytes[0] = 0

		second := schema.Files()
		assert.NotEqual(t, byte(0), second[0].Bytes[0])
	})
}

func TestRevision(t *testing.T) {
	t.Run("parses the x-esdm-schema-revision out of a YAML schema", func(t *testing.T) {
		rev, err := schema.Revision(schema.Core())
		require.NoError(t, err)

		assert.Equal(t, "1.0.0", rev)
	})

	t.Run("returns an error when the field is missing", func(t *testing.T) {
		_, err := schema.Revision([]byte("$id: foo\n"))
		assert.Error(t, err)
	})

	t.Run("returns an error when the value is not a SemVer", func(t *testing.T) {
		_, err := schema.Revision([]byte("x-esdm-schema-revision: not-semver\n"))
		assert.Error(t, err)
	})
}

func TestAPIVersions(t *testing.T) {
	t.Run("returns the apiVersion of every embedded schema", func(t *testing.T) {
		versions, err := schema.APIVersions()
		require.NoError(t, err)

		assert.Contains(t, versions, "schema.esdm.io/core/v1")
		assert.Contains(t, versions, "schema.esdm.io/domain-storytelling/v1")
	})
}
