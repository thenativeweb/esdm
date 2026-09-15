package docgen

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thenativeweb/esdm/modelpath"
	"github.com/thenativeweb/esdm/runner"
)

var (
	directory string
	output    string
	force     bool
)

func init() {
	Command.Flags().StringVarP(&directory, "directory", "d", ".", "directory containing the model to document")
	Command.Flags().StringVarP(&output, "output", "o", "", "directory to write the Markdown tree to")
	Command.Flags().BoolVar(&force, "force", false, "clear and rewrite the output directory when it is not empty")
	_ = Command.MarkFlagRequired("output")
}

// Command is the cobra command instance registered by the root
// command. It runs the resolver pipeline and writes the model as a
// Markdown directory tree to the output directory. Linter findings do
// not block the output; only an unresolvable model or invalid input is
// treated as an error.
var Command = &cobra.Command{
	Use:           "documentation [path]",
	Short:         "Renders an ESDM model as a Markdown directory tree",
	Long:          "Renders the ESDM model in --directory as a Markdown directory tree written to --output. Each element becomes a page at its containment path, so GitHub and MkDocs can both read the result.",
	Example:       "  esdm documentation --output ./docs\n  esdm documentation --output ./docs <domain>/<bounded-context>\n  esdm documentation --directory ./model --output ./docs --force",
	Args:          cobra.MaximumNArgs(1),
	SilenceUsage:  true,
	SilenceErrors: false,
	RunE: func(command *cobra.Command, args []string) error {
		var rawPath string
		if len(args) > 0 {
			rawPath = args[0]
		}
		path, err := modelpath.ParsePath(rawPath)
		if err != nil {
			return err
		}

		_, m, err := runner.RunWithModel(command.Context(), directory)
		if err != nil {
			return err
		}
		if m == nil {
			return fmt.Errorf("the resolver could not produce a model for %q; run `esdm lint` for diagnostics", directory)
		}

		pages, err := Build(m, path)
		if err != nil {
			return err
		}

		return Write(pages, output, force)
	},
}
