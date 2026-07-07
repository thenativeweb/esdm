package lint

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"

	"github.com/thenativeweb/esdm/cmd/cmdutils"
	"github.com/thenativeweb/esdm/diag"
	"github.com/thenativeweb/esdm/reporter"
	"github.com/thenativeweb/esdm/runner"
)

var (
	format           string
	colorMode        string
	directory        string
	warningsAsErrors bool
)

const (
	formatHuman = "human"
	formatJSON  = "json"
)

// ExitCodeHasErrors is the exit code used when lint
// completes successfully but produced one or more
// error-severity diagnostics.
const ExitCodeHasErrors = 1

// exitCode captures the intended process exit code so
// main can read it via ExitCode() after Execute returns
// and propagate it via os.Exit. Reset to 0 at the start
// of every Run so a fresh invocation never inherits a
// previous run's outcome.
var exitCode = 0

// ExitCode returns the exit code the lint command wants
// the process to exit with: 0 on success, ExitCodeHasErrors
// when any error-severity diagnostic was produced. main is
// the only consumer in production; tests use it to assert
// on the propagation.
func ExitCode() int {
	return exitCode
}

func init() {
	Command.Flags().StringVar(&format, "format", formatHuman, "output format: human or json")
	Command.Flags().StringVar(&colorMode, "color", cmdutils.ColorAuto, "colorize human output: auto, always, or never")
	Command.Flags().StringVarP(&directory, "directory", "d", ".", "directory containing the model to lint")
	Command.Flags().BoolVar(&warningsAsErrors, "warnings-as-errors", false, "treat warning-severity findings as errors for exit-code purposes")
}

var Command = &cobra.Command{
	Use:           "lint",
	Short:         "Lints an ESDM model",
	Long:          "Lints all .esdm.yaml files in --directory (default: the current directory).",
	Args:          cobra.NoArgs,
	SilenceUsage:  true,
	SilenceErrors: false,
	RunE: func(command *cobra.Command, args []string) error {
		exitCode = 0

		diagnostics, err := runner.Run(command.Context(), directory)
		if err != nil {
			return err
		}

		out := command.OutOrStdout()
		formatter, err := selectFormatter(format, colorMode, out)
		if err != nil {
			return err
		}

		err = formatter.Format(out, diagnostics)
		if err != nil {
			return err
		}

		for _, d := range diagnostics {
			isError := d.Severity == diag.SeverityError
			isWarningEscalated := warningsAsErrors && d.Severity == diag.SeverityWarning
			if isError || isWarningEscalated {
				exitCode = ExitCodeHasErrors
				return nil
			}
		}

		return nil
	},
}

func selectFormatter(name, color string, out io.Writer) (reporter.Formatter, error) {
	switch name {
	case formatHuman:
		humanFormatter := reporter.NewHumanFormatter()
		shouldUseColor, err := cmdutils.ResolveColor(color, out)
		if err != nil {
			return nil, err
		}
		humanFormatter.Colors = shouldUseColor
		return humanFormatter, nil
	case formatJSON:
		return reporter.NewJSONFormatter(), nil
	default:
		return nil, fmt.Errorf("unknown --format value %q (expected %q or %q)", name, formatHuman, formatJSON)
	}
}
