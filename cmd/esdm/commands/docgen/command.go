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
	Command.Flags().StringVarP(&output, "output", "o", "", "directory to write the Markdown tree to (required)")
	Command.Flags().BoolVar(&force, "force", false, "clear a non-empty output directory before writing")
	_ = Command.MarkFlagRequired("output")
}

// Command is the cobra command instance registered by the
// root command. It runs the resolver pipeline, renders the
// model's tree as Markdown pages, and writes them to the
// output directory. Linter findings do not block the output;
// only an unresolvable model, an invalid path argument, or a
// non-empty output directory without --force is an error.
var Command = &cobra.Command{
	Use:           "documentation [path]",
	Short:         "Writes an ESDM model as a tree of Markdown pages",
	Long:          "Writes an ESDM model in --directory as a tree of Markdown pages to --output, one page per element.",
	Example:       "  esdm documentation --output docs\n  esdm documentation --output docs <domain>/<bounded-context>\n  esdm documentation --output docs --force",
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
