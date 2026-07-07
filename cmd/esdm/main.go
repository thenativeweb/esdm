package main

import (
	"os"

	"github.com/thenativeweb/esdm/cmd/cmdutils"
	"github.com/thenativeweb/esdm/cmd/esdm/commands/lint"
	"github.com/thenativeweb/esdm/cmd/esdm/commands/root"
	"github.com/thenativeweb/esdm/logging"
)

func main() {
	cmdutils.Ensure64BitArchitecture()

	err := root.Command.Execute()
	if err != nil {
		logging.Fatal("failed to execute command", "error", err)
	}

	if code := lint.ExitCode(); code != 0 {
		os.Exit(code)
	}
}
