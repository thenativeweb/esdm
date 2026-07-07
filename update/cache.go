package update

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

const cacheFileName = "version-check.json"

// CacheEntry is the on-disk record of the most recent
// version check. NextCheckAt schedules when the next HTTP
// request may run; LatestVersion holds the most recent
// successfully fetched version (empty if no successful
// fetch has happened yet); FirstRunShown becomes true once
// the first-time variant of the notification has been
// rendered.
type CacheEntry struct {
	NextCheckAt   time.Time `json:"next_check_at"`
	LatestVersion string    `json:"latest_version,omitempty"`
	FirstRunShown bool      `json:"first_run_shown,omitempty"`
}

// Cache reads and writes the version-check cache file.
// Concurrent calls are safe at the OS level because writes
// go through an atomic rename.
type Cache struct {
	dir string
}

// NewCache returns a Cache rooted at dir. The directory
// is created on first write.
func NewCache(dir string) *Cache {
	return &Cache{dir: dir}
}

// Path returns the absolute path of the cache file.
func (c *Cache) Path() string {
	return filepath.Join(c.dir, cacheFileName)
}

// Read returns the stored entry, or nil if no entry is
// available. Missing or corrupt files are treated as no
// entry; callers do not need to distinguish.
func (c *Cache) Read() (*CacheEntry, error) {
	data, err := os.ReadFile(c.Path())
	if errors.Is(err, fs.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var entry CacheEntry
	err = json.Unmarshal(data, &entry)
	if err != nil {
		return nil, nil
	}
	return &entry, nil
}

// Write persists entry to the cache file via a temp file
// and atomic rename.
func (c *Cache) Write(entry *CacheEntry) error {
	err := os.MkdirAll(c.dir, 0o755)
	if err != nil {
		return err
	}

	data, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	temporaryFile, err := os.CreateTemp(c.dir, cacheFileName+".tmp.*")
	if err != nil {
		return err
	}
	temporaryFilePath := temporaryFile.Name()

	_, err = temporaryFile.Write(data)
	if err != nil {
		temporaryFile.Close()
		os.Remove(temporaryFilePath)
		return err
	}
	err = temporaryFile.Close()
	if err != nil {
		os.Remove(temporaryFilePath)
		return err
	}

	return os.Rename(temporaryFilePath, c.Path())
}
