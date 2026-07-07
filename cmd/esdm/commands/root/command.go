package root

import (
	"github.com/spf13/cobra"
	"github.com/thenativeweb/esdm/cmd/esdm/commands/addschema"
	"github.com/thenativeweb/esdm/cmd/esdm/commands/glossary"
	"github.com/thenativeweb/esdm/cmd/esdm/commands/lint"
	"github.com/thenativeweb/esdm/cmd/esdm/commands/updateschema"
	"github.com/thenativeweb/esdm/cmd/esdm/commands/version"
	"github.com/thenativeweb/esdm/cmd/esdm/commands/view"
)

func init() {
	Command.AddCommand(addschema.Command)
	Command.AddCommand(glossary.Command)
	Command.AddCommand(lint.Command)
	Command.AddCommand(updateschema.Command)
	Command.AddCommand(version.Command)
	Command.AddCommand(view.Command)
}

var Command = &cobra.Command{
	Use:   "esdm",
	Short: "esdm - tools for event-sourced domain models",
	Long:  "esdm - tools for event-sourced domain models.",
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		runUpdateCheck(cmd)
	},
}
