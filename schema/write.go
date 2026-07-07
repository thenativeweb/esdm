package schema

import (
	"os"
	"path/filepath"
)

// Write materializes every embedded schema file into
// schemasRoot, creating the directory and any required
// subdirectories. It overwrites existing files at the
// same paths but does not remove unrelated files; the
// `add-schema` and `update-schema` commands decide on
// those preconditions before calling Write.
func Write(schemasRoot string) error {
	err := os.MkdirAll(schemasRoot, 0o755)
	if err != nil {
		return err
	}

	for _, file := range Files() {
		destination := filepath.Join(schemasRoot, filepath.FromSlash(file.Path))
		err := os.MkdirAll(filepath.Dir(destination), 0o755)
		if err != nil {
			return err
		}
		err = os.WriteFile(destination, file.Bytes, 0o644)
		if err != nil {
			return err
		}
	}

	return nil
}
