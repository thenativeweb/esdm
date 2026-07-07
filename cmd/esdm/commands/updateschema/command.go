package updateschema

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/thenativeweb/esdm/schema"
)

// schemasSubdir mirrors the constant in addschema; the
// two commands share the directory layout but not their
// preconditions, so each carries its own copy rather than
// importing across packages.
const schemasSubdir = "schemas"

var Command = &cobra.Command{
	Use:   "update-schema",
	Short: "Refreshes the local ESDM schemas to match the embedded revision",
	Long:  "Refreshes the local schemas/ directory to match the embedded ESDM schemas.",
	Args:  cobra.NoArgs,
	RunE: func(command *cobra.Command, args []string) error {
		return updateSchema(command.OutOrStdout())
	},
}

func updateSchema(stdout io.Writer) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	schemasRoot := filepath.Join(cwd, schemasSubdir)

	info, err := os.Stat(schemasRoot)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("%q does not exist; run add-schema to create it", schemasRoot)
		}
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("%q exists but is not a directory", schemasRoot)
	}

	needsRewrite, err := planUpdate(schemasRoot)
	if err != nil {
		return err
	}

	if !needsRewrite {
		fmt.Fprintln(stdout, "schemas/ is up to date - nothing to update.")
		return nil
	}

	err = os.RemoveAll(schemasRoot)
	if err != nil {
		return err
	}
	return schema.Write(schemasRoot)
}

// planUpdate inspects the local schemas/ tree and decides
// whether a wipe-and-rewrite is warranted. It returns an
// error when a local revision is strictly higher than the
// embedded one (unsupported downgrade).
func planUpdate(schemasRoot string) (bool, error) {
	embedded := schema.Files()

	embeddedPaths := make(map[string]schema.File, len(embedded))
	for _, file := range embedded {
		embeddedPaths[file.Path] = file
	}

	hasByteDifferences := false
	for _, file := range embedded {
		localPath := filepath.Join(schemasRoot, filepath.FromSlash(file.Path))
		localBytes, err := os.ReadFile(localPath)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				hasByteDifferences = true
				continue
			}
			return false, err
		}

		if !bytes.Equal(localBytes, file.Bytes) {
			hasByteDifferences = true
		}

		localRevision, err := schema.Revision(localBytes)
		if err != nil {
			return false, fmt.Errorf("local %q: %w", file.Path, err)
		}
		embeddedRevision, err := schema.Revision(file.Bytes)
		if err != nil {
			return false, err
		}

		ordering, err := schema.CompareRevisions(localRevision, embeddedRevision)
		if err != nil {
			return false, err
		}
		if ordering > 0 {
			return false, fmt.Errorf(
				"local schema %q has revision %s, which is newer than this binary's %s; install a newer esdm to update",
				file.Path, localRevision, embeddedRevision,
			)
		}
	}

	hasStrayFiles, err := hasUnknownPaths(schemasRoot, embeddedPaths)
	if err != nil {
		return false, err
	}

	return hasByteDifferences || hasStrayFiles, nil
}

func hasUnknownPaths(schemasRoot string, known map[string]schema.File) (bool, error) {
	hasStray := false
	err := filepath.WalkDir(schemasRoot, func(path string, d fs.DirEntry, err error) error {
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
		if _, ok := known[filepath.ToSlash(rel)]; !ok {
			hasStray = true
		}
		return nil
	})
	if err != nil {
		return false, err
	}
	return hasStray, nil
}
