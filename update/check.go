package update

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	successTTL  = 24 * time.Hour
	failureTTL  = 1 * time.Hour
	httpTimeout = 2 * time.Second
	maxBodySize = 4 * 1024
)

// RunOptions bundles the dependencies of the version
// check. Stderr is where the notification is written;
// CacheDir is the directory holding the cache file. The
// HTTPClient and Now fields are optional and exist purely
// for tests.
type RunOptions struct {
	CurrentVersion string
	Endpoint       string
	CacheDir       string
	Stderr         io.Writer
	Color          bool
	HTTPClient     *http.Client
	Now            func() time.Time
}

// Run performs the version check end-to-end. It reads
// the cache, refreshes it via HTTP when the schedule
// allows, renders the notification when a newer version
// is available, and updates the first-run state. The
// function never returns an error: every failure path
// degrades silently. It blocks at most for the HTTP
// timeout when a refresh is due.
func Run(ctx context.Context, options RunOptions) {
	now := nowOrDefault(options.Now)
	cache := NewCache(options.CacheDir)

	entry, _ := cache.Read()
	entry = refreshIfStale(ctx, entry, options, cache, now)

	if entry == nil || entry.LatestVersion == "" {
		return
	}

	isNewer, err := IsNewer(options.CurrentVersion, entry.LatestVersion)
	if err != nil || !isNewer {
		return
	}

	notification := RenderNotification(options.CurrentVersion, entry.LatestVersion, !entry.FirstRunShown, options.Color)
	fmt.Fprintln(options.Stderr, notification)

	if !entry.FirstRunShown {
		entry.FirstRunShown = true
		_ = cache.Write(entry)
	}
}

func refreshIfStale(ctx context.Context, current *CacheEntry, options RunOptions, cache *Cache, now time.Time) *CacheEntry {
	if current != nil && now.Before(current.NextCheckAt) {
		return current
	}

	fetched, err := fetchLatestVersion(ctx, options.HTTPClient, options.Endpoint)

	next := &CacheEntry{}
	if current != nil {
		next.LatestVersion = current.LatestVersion
		next.FirstRunShown = current.FirstRunShown
	}
	if err != nil {
		next.NextCheckAt = now.Add(failureTTL)
	} else {
		next.NextCheckAt = now.Add(successTTL)
		next.LatestVersion = fetched
	}

	_ = cache.Write(next)
	return next
}

func fetchLatestVersion(ctx context.Context, client *http.Client, endpoint string) (string, error) {
	if client == nil {
		client = &http.Client{Timeout: httpTimeout}
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return "", err
	}

	response, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status %d", response.StatusCode)
	}

	body, err := io.ReadAll(io.LimitReader(response.Body, maxBodySize))
	if err != nil {
		return "", err
	}

	var payload struct {
		Version string `json:"version"`
	}
	err = json.Unmarshal(body, &payload)
	if err != nil {
		return "", err
	}
	if payload.Version == "" {
		return "", fmt.Errorf("missing version field in payload")
	}
	return payload.Version, nil
}

func nowOrDefault(now func() time.Time) time.Time {
	if now == nil {
		return time.Now()
	}
	return now()
}
