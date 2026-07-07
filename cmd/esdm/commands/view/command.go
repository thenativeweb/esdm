package view

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"

	"github.com/thenativeweb/esdm/cmd/cmdutils"
	"github.com/thenativeweb/esdm/modelpath"
	"github.com/thenativeweb/esdm/runner"
)

var (
	withDetails bool
	colorMode   string
	directory   string
)

func init() {
	Command.Flags().BoolVar(&withDetails, "with-details", false, "show node-level details (schemas, invariants, rule prose) in addition to the skeleton")
	Command.Flags().StringVar(&colorMode, "color", cmdutils.ColorAuto, "colorize output: auto, always, or never")
	Command.Flags().StringVarP(&directory, "directory", "d", ".", "directory containing the model to summarize")
}

// Command is the cobra command instance registered by
// the root command. The implementation runs the resolver
// and rule pipeline implicitly, builds a render tree
// from the resolved model, annotates each node with
// matching linter diagnostics, and writes the rendered
// text to stdout.
var Command = &cobra.Command{
	Use:           "view [path]",
	Short:         "Renders a hierarchical summary of an ESDM model",
	Long:          "Renders a hierarchical summary of an ESDM model in --directory.",
	Example:       "  esdm view\n  esdm view <domain>/<bounded-context>/<aggregate>\n  esdm view --with-details",
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

		diagnostics, m, err := runner.RunWithModel(command.Context(), directory)
		if err != nil {
			return err
		}
		if m == nil {
			return fmt.Errorf("the resolver could not produce a model for %q; run `esdm lint` for diagnostics", directory)
		}

		root, err := BuildTree(m, path, withDetails)
		if err != nil {
			return err
		}
		Annotate(root, diagnostics)

		out := command.OutOrStdout()
		shouldUseColor, err := cmdutils.ResolveColor(colorMode, out)
		if err != nil {
			return err
		}
		opts := RenderOptions{
			Colors:      shouldUseColor,
			ShowDetails: withDetails,
		}
		_, err = io.WriteString(out, Render(root, opts))
		if err != nil {
			return err
		}
		return nil
	},
}
