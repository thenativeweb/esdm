package update_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thenativeweb/esdm/update"
)

func TestCache(t *testing.T) {
	t.Run("Read returns nil when the cache file does not exist", func(t *testing.T) {
		dir := t.TempDir()
		cache := update.NewCache(dir)

		entry, err := cache.Read()
		require.NoError(t, err)
		assert.Nil(t, entry)
	})

	t.Run("Write followed by Read round-trips the entry", func(t *testing.T) {
		dir := t.TempDir()
		cache := update.NewCache(dir)

		original := &update.CacheEntry{
			NextCheckAt:   time.Date(2026, 5, 10, 12, 0, 0, 0, time.UTC),
			LatestVersion: "0.10.0",
			FirstRunShown: true,
		}
		require.NoError(t, cache.Write(original))

		loaded, err := cache.Read()
		require.NoError(t, err)
		require.NotNil(t, loaded)
		assert.True(t, original.NextCheckAt.Equal(loaded.NextCheckAt))
		assert.Equal(t, original.LatestVersion, loaded.LatestVersion)
		assert.Equal(t, original.FirstRunShown, loaded.FirstRunShown)
	})

	t.Run("Write creates the parent directory if it does not exist", func(t *testing.T) {
		dir := filepath.Join(t.TempDir(), "missing", "subdir")
		cache := update.NewCache(dir)

		require.NoError(t, cache.Write(&update.CacheEntry{
			NextCheckAt: time.Now(),
		}))

		_, err := os.Stat(filepath.Join(dir, "version-check.json"))
		assert.NoError(t, err)
	})

	t.Run("a corrupt cache file is treated as no entry", func(t *testing.T) {
		dir := t.TempDir()
		require.NoError(t, os.WriteFile(filepath.Join(dir, "version-check.json"), []byte("not json"), 0o600))

		cache := update.NewCache(dir)
		entry, err := cache.Read()
		require.NoError(t, err)
		assert.Nil(t, entry)
	})

	t.Run("Write does not leave temp files behind on success", func(t *testing.T) {
		dir := t.TempDir()
		cache := update.NewCache(dir)

		require.NoError(t, cache.Write(&update.CacheEntry{
			NextCheckAt: time.Now(),
		}))

		entries, err := os.ReadDir(dir)
		require.NoError(t, err)
		require.Len(t, entries, 1)
		assert.Equal(t, "version-check.json", entries[0].Name())
	})
}
