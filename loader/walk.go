package loader

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// FileSuffix is the filename suffix that marks a file as
// part of an ESDM model.
const FileSuffix = ".esdm.yaml"

// ErrNoFiles is returned by Walk when the given directory
// exists but contains no files matching the ESDM file
// suffix.
var ErrNoFiles = errors.New("no esdm files found")

// Walk recursively collects every file under dir whose
// name ends in ".esdm.yaml". The returned slice is sorted
// lexicographically so that subsequent processing (parsing,
// diagnostics) is deterministic across runs.
//
// Walk returns an error if dir does not exist, cannot be
// read, or contains no matching files (ErrNoFiles).
func Walk(dir string) ([]string, error) {
	info, err := os.Stat(dir)
	if err != nil {
		return nil, fmt.Errorf("cannot access %q: %w", dir, err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("%q is not a directory", dir)
	}

	var paths []string
	walkErr := filepath.WalkDir(dir, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		if !strings.HasSuffix(entry.Name(), FileSuffix) {
			return nil
		}

		paths = append(paths, path)
		return nil
	})
	if walkErr != nil {
		return nil, walkErr
	}

	if len(paths) == 0 {
		return nil, ErrNoFiles
	}

	sort.Strings(paths)
	return paths, nil
}
