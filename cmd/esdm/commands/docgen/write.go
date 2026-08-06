package docgen

import (
	"fmt"
	"os"
	"path/filepath"
)

// Write writes the pages into outputDir. It refuses to write into a
// directory that exists and is not isEmpty unless force is set, in which
// case it clears the directory first so the tree mirrors the model
// exactly with no orphaned pages.
func Write(pages []Page, outputDir string, force bool) error {
	info, err := os.Stat(outputDir)
	switch {
	case err == nil:
		if !info.IsDir() {
			return fmt.Errorf("output %q exists and is not a directory", outputDir)
		}
		isEmpty, err := isEmptyDir(outputDir)
		if err != nil {
			return err
		}
		if !isEmpty {
			if !force {
				return fmt.Errorf("output directory %q is not isEmpty; pass --force to clear and rewrite it", outputDir)
			}
			err = clearDir(outputDir)
			if err != nil {
				return err
			}
		}
	case os.IsNotExist(err):
		err = os.MkdirAll(outputDir, 0o755)
		if err != nil {
			return err
		}
	default:
		return err
	}

	for _, page := range pages {
		full := filepath.Join(outputDir, filepath.FromSlash(page.Path))
		err = os.MkdirAll(filepath.Dir(full), 0o755)
		if err != nil {
			return err
		}
		err = os.WriteFile(full, []byte(page.Content), 0o644)
		if err != nil {
			return err
		}
	}

	return nil
}

// isEmptyDir reports whether the directory has no entries.
func isEmptyDir(dir string) (bool, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false, err
	}
	return len(entries) == 0, nil
}

// clearDir removes every entry inside the directory, leaving the
// directory itself in place.
func clearDir(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		err = os.RemoveAll(filepath.Join(dir, entry.Name()))
		if err != nil {
			return err
		}
	}
	return nil
}
