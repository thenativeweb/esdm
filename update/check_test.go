package update_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thenativeweb/esdm/update"
)

func newVersionServer(t *testing.T, version string) *httptest.Server {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"version": version})
	}))
	t.Cleanup(server.Close)
	return server
}

func newFailingServer(t *testing.T, status int) *httptest.Server {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
	}))
	t.Cleanup(server.Close)
	return server
}

func fixedNow(at time.Time) func() time.Time {
	return func() time.Time {
		return at
	}
}

func TestRun(t *testing.T) {
	t.Run("renders the first-run notification when no cache exists yet", func(t *testing.T) {
		server := newVersionServer(t, "0.10.0")
		var stderr bytes.Buffer

		update.Run(context.Background(), update.RunOptions{
			CurrentVersion: "v0.9.0",
			Endpoint:       server.URL,
			CacheDir:       t.TempDir(),
			Stderr:         &stderr,
			Now:            fixedNow(time.Date(2026, 5, 9, 12, 0, 0, 0, time.UTC)),
		})

		out := stderr.String()
		assert.Contains(t, out, "0.9.0 → 0.10.0")
		assert.Contains(t, out, update.UpgradeURL)
		assert.Contains(t, out, update.DisableEnvVar+"=true")
	})

	t.Run("marks first-run shown so the next call omits the disable hint", func(t *testing.T) {
		server := newVersionServer(t, "0.10.0")
		cacheDir := t.TempDir()

		var firstStderr bytes.Buffer
		now := time.Date(2026, 5, 9, 12, 0, 0, 0, time.UTC)
		update.Run(context.Background(), update.RunOptions{
			CurrentVersion: "v0.9.0",
			Endpoint:       server.URL,
			CacheDir:       cacheDir,
			Stderr:         &firstStderr,
			Now:            fixedNow(now),
		})
		require.Contains(t, firstStderr.String(), update.DisableEnvVar)

		var secondStderr bytes.Buffer
		update.Run(context.Background(), update.RunOptions{
			CurrentVersion: "v0.9.0",
			Endpoint:       server.URL,
			CacheDir:       cacheDir,
			Stderr:         &secondStderr,
			Now:            fixedNow(now.Add(time.Minute)),
		})

		out := secondStderr.String()
		assert.Contains(t, out, "0.9.0 → 0.10.0")
		assert.NotContains(t, out, update.DisableEnvVar)
	})

	t.Run("emits no notification when the latest version equals the current one", func(t *testing.T) {
		server := newVersionServer(t, "0.9.0")
		var stderr bytes.Buffer

		update.Run(context.Background(), update.RunOptions{
			CurrentVersion: "v0.9.0",
			Endpoint:       server.URL,
			CacheDir:       t.TempDir(),
			Stderr:         &stderr,
		})

		assert.Empty(t, stderr.String())
	})

	t.Run("emits no notification when the latest version is older than the current one", func(t *testing.T) {
		server := newVersionServer(t, "0.8.0")
		var stderr bytes.Buffer

		update.Run(context.Background(), update.RunOptions{
			CurrentVersion: "v0.9.0",
			Endpoint:       server.URL,
			CacheDir:       t.TempDir(),
			Stderr:         &stderr,
		})

		assert.Empty(t, stderr.String())
	})

	t.Run("does not call the endpoint while the cache is fresh", func(t *testing.T) {
		var calls atomic.Int32
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			calls.Add(1)
			_ = json.NewEncoder(w).Encode(map[string]string{"version": "0.10.0"})
		}))
		t.Cleanup(server.Close)

		cacheDir := t.TempDir()
		now := time.Date(2026, 5, 9, 12, 0, 0, 0, time.UTC)

		update.Run(context.Background(), update.RunOptions{
			CurrentVersion: "v0.9.0",
			Endpoint:       server.URL,
			CacheDir:       cacheDir,
			Stderr:         &bytes.Buffer{},
			Now:            fixedNow(now),
		})
		update.Run(context.Background(), update.RunOptions{
			CurrentVersion: "v0.9.0",
			Endpoint:       server.URL,
			CacheDir:       cacheDir,
			Stderr:         &bytes.Buffer{},
			Now:            fixedNow(now.Add(1 * time.Hour)),
		})

		assert.Equal(t, int32(1), calls.Load())
	})

	t.Run("refetches once the cache passes its scheduled refresh time", func(t *testing.T) {
		var calls atomic.Int32
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			calls.Add(1)
			_ = json.NewEncoder(w).Encode(map[string]string{"version": "0.10.0"})
		}))
		t.Cleanup(server.Close)

		cacheDir := t.TempDir()
		now := time.Date(2026, 5, 9, 12, 0, 0, 0, time.UTC)

		update.Run(context.Background(), update.RunOptions{
			CurrentVersion: "v0.9.0",
			Endpoint:       server.URL,
			CacheDir:       cacheDir,
			Stderr:         &bytes.Buffer{},
			Now:            fixedNow(now),
		})
		update.Run(context.Background(), update.RunOptions{
			CurrentVersion: "v0.9.0",
			Endpoint:       server.URL,
			CacheDir:       cacheDir,
			Stderr:         &bytes.Buffer{},
			Now:            fixedNow(now.Add(25 * time.Hour)),
		})

		assert.Equal(t, int32(2), calls.Load())
	})

	t.Run("emits no notification when the endpoint fails and there is no previous cache", func(t *testing.T) {
		server := newFailingServer(t, http.StatusInternalServerError)
		var stderr bytes.Buffer

		update.Run(context.Background(), update.RunOptions{
			CurrentVersion: "v0.9.0",
			Endpoint:       server.URL,
			CacheDir:       t.TempDir(),
			Stderr:         &stderr,
		})

		assert.Empty(t, stderr.String())
	})

	t.Run("falls back to the previous cache value when the endpoint fails", func(t *testing.T) {
		cacheDir := t.TempDir()
		require.NoError(t, update.NewCache(cacheDir).Write(&update.CacheEntry{
			NextCheckAt:   time.Date(2026, 5, 9, 0, 0, 0, 0, time.UTC),
			LatestVersion: "0.10.0",
			FirstRunShown: true,
		}))

		server := newFailingServer(t, http.StatusInternalServerError)
		var stderr bytes.Buffer

		update.Run(context.Background(), update.RunOptions{
			CurrentVersion: "v0.9.0",
			Endpoint:       server.URL,
			CacheDir:       cacheDir,
			Stderr:         &stderr,
			Now:            fixedNow(time.Date(2026, 5, 10, 0, 0, 0, 0, time.UTC)),
		})

		assert.Contains(t, stderr.String(), "0.9.0 → 0.10.0")
	})

	t.Run("schedules a short retry after a failed fetch", func(t *testing.T) {
		cacheDir := t.TempDir()
		server := newFailingServer(t, http.StatusInternalServerError)
		now := time.Date(2026, 5, 9, 12, 0, 0, 0, time.UTC)

		update.Run(context.Background(), update.RunOptions{
			CurrentVersion: "v0.9.0",
			Endpoint:       server.URL,
			CacheDir:       cacheDir,
			Stderr:         &bytes.Buffer{},
			Now:            fixedNow(now),
		})

		entry, err := update.NewCache(cacheDir).Read()
		require.NoError(t, err)
		require.NotNil(t, entry)
		assert.True(t, entry.NextCheckAt.Sub(now) <= time.Hour, "next check should be scheduled within an hour after failure")
		assert.True(t, entry.NextCheckAt.After(now), "next check should be in the future")
	})

	t.Run("emits no notification when the current version cannot be parsed", func(t *testing.T) {
		server := newVersionServer(t, "0.10.0")
		var stderr bytes.Buffer

		update.Run(context.Background(), update.RunOptions{
			CurrentVersion: "(version unavailable)",
			Endpoint:       server.URL,
			CacheDir:       t.TempDir(),
			Stderr:         &stderr,
		})

		assert.Empty(t, stderr.String())
	})

	t.Run("ignores a malformed JSON payload from the endpoint", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("not json"))
		}))
		t.Cleanup(server.Close)
		var stderr bytes.Buffer

		update.Run(context.Background(), update.RunOptions{
			CurrentVersion: "v0.9.0",
			Endpoint:       server.URL,
			CacheDir:       t.TempDir(),
			Stderr:         &stderr,
		})

		assert.Empty(t, stderr.String())
	})

	t.Run("writes the cache file under the configured directory", func(t *testing.T) {
		server := newVersionServer(t, "0.10.0")
		cacheDir := t.TempDir()

		update.Run(context.Background(), update.RunOptions{
			CurrentVersion: "v0.9.0",
			Endpoint:       server.URL,
			CacheDir:       cacheDir,
			Stderr:         &bytes.Buffer{},
		})

		entry, err := update.NewCache(cacheDir).Read()
		require.NoError(t, err)
		require.NotNil(t, entry)
		assert.Equal(t, "0.10.0", entry.LatestVersion)
		assert.True(t, entry.FirstRunShown)
		assert.FileExists(t, filepath.Join(cacheDir, "version-check.json"))
	})
}
