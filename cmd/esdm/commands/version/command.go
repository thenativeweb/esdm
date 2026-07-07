package version

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/thenativeweb/esdm/cmd/cmdutils"
)

var Command = &cobra.Command{
	Use:   "version",
	Short: "Prints the esdm version",
	Long:  "Prints the esdm version.",
	RunE: func(command *cobra.Command, args []string) error {
		fmt.Fprintln(command.OutOrStdout(), "esdm: "+cmdutils.Version)
		fmt.Fprintln(command.OutOrStdout(), "Revision: "+cmdutils.GitVersion)
		fmt.Fprintln(command.OutOrStdout())

		return nil
	},
}
