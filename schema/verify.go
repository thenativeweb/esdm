package schema

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
)

// VerifyError describes the first deviation Verify
// detected between the on-disk schemas/ directory and the
// embedded schema set. Reason carries a human-readable
// explanation; Path is the offending project-relative
// path inside schemas/, when applicable.
type VerifyError struct {
	Path   string
	Reason string
}

func (e *VerifyError) Error() string {
	if e.Path == "" {
		return e.Reason
	}
	return fmt.Sprintf("%s: %s", e.Path, e.Reason)
}

// Verify checks that schemasRoot is byte-for-byte equal
// to the embedded schema set: every embedded file is
// present with the exact bytes, and no other files exist
// underneath schemasRoot. Any deviation - missing files,
// extra files, content drift - yields a VerifyError. The
// linter calls this at startup so that local copies that
// drift from what the binary expects fail loudly instead
// of silently masking schema inconsistencies.
func Verify(schemasRoot string) error {
	embedded := Files()

	expected := make(map[string][]byte, len(embedded))
	for _, file := range embedded {
		expected[file.Path] = file.Bytes
	}

	seen := make(map[string]bool, len(expected))

	walkErr := filepath.WalkDir(schemasRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		rel, err := filepath.Rel(schemasRoot, path)
		if err != nil {
			return err
		}
		relSlash := filepath.ToSlash(rel)

		want, ok := expected[relSlash]
		if !ok {
			return &VerifyError{Path: relSlash, Reason: "unexpected file (not part of the embedded schema set)"}
		}
		seen[relSlash] = true

		got, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if !bytes.Equal(got, want) {
			return &VerifyError{Path: relSlash, Reason: "content differs from the embedded schema"}
		}
		return nil
	})
	if walkErr != nil {
		var ve *VerifyError
		if errors.As(walkErr, &ve) {
			return ve
		}
		return walkErr
	}

	missing := missingPaths(expected, seen)
	if len(missing) > 0 {
		return &VerifyError{Path: missing[0], Reason: "expected file is missing from the local schemas directory"}
	}

	return nil
}

func missingPaths(expected map[string][]byte, seen map[string]bool) []string {
	var out []string
	for path := range expected {
		if !seen[path] {
			out = append(out, path)
		}
	}
	sort.Strings(out)
	return out
}
