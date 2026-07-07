package loader_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thenativeweb/esdm/loader"
)

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	require.NoError(t, os.MkdirAll(filepath.Dir(path), 0o755))
	require.NoError(t, os.WriteFile(path, []byte(content), 0o644))
}

func TestWalk(t *testing.T) {
	t.Run("returns ErrNoFiles when directory contains no matching files", func(t *testing.T) {
		dir := t.TempDir()
		writeFile(t, filepath.Join(dir, "README.md"), "not esdm")

		paths, err := loader.Walk(dir)
		assert.ErrorIs(t, err, loader.ErrNoFiles)
		assert.Nil(t, paths)
	})

	t.Run("returns ErrNoFiles for an empty directory", func(t *testing.T) {
		dir := t.TempDir()

		paths, err := loader.Walk(dir)
		assert.ErrorIs(t, err, loader.ErrNoFiles)
		assert.Nil(t, paths)
	})

	t.Run("finds files at the top level", func(t *testing.T) {
		dir := t.TempDir()
		writeFile(t, filepath.Join(dir, "a.esdm.yaml"), "---")
		writeFile(t, filepath.Join(dir, "b.esdm.yaml"), "---")

		paths, err := loader.Walk(dir)
		require.NoError(t, err)

		assert.Equal(t, []string{
			filepath.Join(dir, "a.esdm.yaml"),
			filepath.Join(dir, "b.esdm.yaml"),
		}, paths)
	})

	t.Run("finds files recursively in subdirectories", func(t *testing.T) {
		dir := t.TempDir()
		writeFile(t, filepath.Join(dir, "top.esdm.yaml"), "---")
		writeFile(t, filepath.Join(dir, "sub", "nested.esdm.yaml"), "---")
		writeFile(t, filepath.Join(dir, "sub", "deep", "deepest.esdm.yaml"), "---")

		paths, err := loader.Walk(dir)
		require.NoError(t, err)

		assert.Equal(t, []string{
			filepath.Join(dir, "sub", "deep", "deepest.esdm.yaml"),
			filepath.Join(dir, "sub", "nested.esdm.yaml"),
			filepath.Join(dir, "top.esdm.yaml"),
		}, paths)
	})

	t.Run("ignores files that do not end in .esdm.yaml", func(t *testing.T) {
		dir := t.TempDir()
		writeFile(t, filepath.Join(dir, "a.esdm.yaml"), "---")
		writeFile(t, filepath.Join(dir, "a.yaml"), "---")
		writeFile(t, filepath.Join(dir, "a.esdm"), "---")
		writeFile(t, filepath.Join(dir, "a.esdm.yml"), "---")

		paths, err := loader.Walk(dir)
		require.NoError(t, err)

		assert.Equal(t, []string{filepath.Join(dir, "a.esdm.yaml")}, paths)
	})

	t.Run("returns sorted results", func(t *testing.T) {
		dir := t.TempDir()
		writeFile(t, filepath.Join(dir, "zebra.esdm.yaml"), "---")
		writeFile(t, filepath.Join(dir, "alpha.esdm.yaml"), "---")
		writeFile(t, filepath.Join(dir, "mike.esdm.yaml"), "---")

		paths, err := loader.Walk(dir)
		require.NoError(t, err)

		assert.Equal(t, []string{
			filepath.Join(dir, "alpha.esdm.yaml"),
			filepath.Join(dir, "mike.esdm.yaml"),
			filepath.Join(dir, "zebra.esdm.yaml"),
		}, paths)
	})

	t.Run("fails when directory does not exist", func(t *testing.T) {
		paths, err := loader.Walk(filepath.Join(t.TempDir(), "does-not-exist"))
		assert.Error(t, err)
		assert.NotErrorIs(t, err, loader.ErrNoFiles)
		assert.Nil(t, paths)
	})

	t.Run("fails when given a file instead of a directory", func(t *testing.T) {
		dir := t.TempDir()
		file := filepath.Join(dir, "some.esdm.yaml")
		writeFile(t, file, "---")

		paths, err := loader.Walk(file)
		assert.Error(t, err)
		assert.Nil(t, paths)
	})
}
