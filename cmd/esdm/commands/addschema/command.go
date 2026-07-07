package addschema

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/thenativeweb/esdm/schema"
)

// schemasSubdir is the single top-level directory the
// add-schema command writes into. Keeping schemas under
// one container leaves the project root clean for the
// user's own .esdm.yaml model files.
const schemasSubdir = "schemas"

var Command = &cobra.Command{
	Use:   "add-schema",
	Short: "Writes the embedded ESDM schemas into the current directory",
	Long:  "Writes the embedded ESDM schemas into a schemas/ directory for editor support.",
	Args:  cobra.NoArgs,
	RunE: func(command *cobra.Command, args []string) error {
		return addSchema()
	},
}

func addSchema() error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	schemasRoot := filepath.Join(cwd, schemasSubdir)

	info, err := os.Stat(schemasRoot)
	if err == nil {
		if info.IsDir() {
			return fmt.Errorf("%q already exists; run update-schema to refresh it", schemasRoot)
		}
		return fmt.Errorf("%q exists but is not a directory", schemasRoot)
	}
	if !errors.Is(err, os.ErrNotExist) {
		return err
	}

	return schema.Write(schemasRoot)
}
