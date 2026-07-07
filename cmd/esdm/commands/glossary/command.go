package glossary

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"

	"github.com/thenativeweb/esdm/modelpath"
	"github.com/thenativeweb/esdm/runner"
)

var directory string

func init() {
	Command.Flags().StringVarP(&directory, "directory", "d", ".", "directory containing the model to read the glossary from")
}

// Command is the cobra command instance registered by the
// root command. It runs the resolver pipeline, extracts the
// ubiquitous language from the selected bounded contexts,
// and writes the glossary to stdout as Markdown. Linter
// findings do not block the output; only an unresolvable
// model or an invalid path argument is treated as an error.
var Command = &cobra.Command{
	Use:           "glossary [path]",
	Short:         "Extracts the ubiquitous language of an ESDM model as Markdown",
	Long:          "Extracts the ubiquitous language of an ESDM model in --directory and writes it to stdout as Markdown.",
	Example:       "  esdm glossary\n  esdm glossary <domain>\n  esdm glossary <domain>/<bounded-context>",
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

		g, err := Build(m, path)
		if err != nil {
			return err
		}

		out := command.OutOrStdout()
		_, err = io.WriteString(out, Render(g))
		if err != nil {
			return err
		}
		return nil
	},
}
