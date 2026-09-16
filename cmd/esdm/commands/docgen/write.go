package docgen

import (
	"fmt"
	"os"
	"path/filepath"
)

// Write puts the pages into the output directory. An
// existing, non-empty directory is refused unless force is
// set, in which case its contents are removed first, so the
// written tree mirrors the model exactly and no page from an
// earlier run survives.
func Write(pages []Page, output string, force bool) error {
	entries, err := os.ReadDir(output)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if len(entries) > 0 {
		if !force {
			return fmt.Errorf("output directory %q is not empty; use --force to clear it first", output)
		}
		for _, entry := range entries {
			err = os.RemoveAll(filepath.Join(output, entry.Name()))
			if err != nil {
				return err
			}
		}
	}

	for _, page := range pages {
		target := filepath.Join(output, filepath.FromSlash(page.Path))
		err = os.MkdirAll(filepath.Dir(target), 0o755)
		if err != nil {
			return err
		}
		err = os.WriteFile(target, []byte(page.Content), 0o644)
		if err != nil {
			return err
		}
	}
	return nil
}
