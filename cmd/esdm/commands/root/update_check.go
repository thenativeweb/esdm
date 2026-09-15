package root

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/thenativeweb/esdm/cmd/cmdutils"
	"github.com/thenativeweb/esdm/update"
)

const (
	devVersionMarker = "(version unavailable)"
	disabledValue    = "true"
)

func runUpdateCheck(cmd *cobra.Command) {
	stderr := cmd.ErrOrStderr()
	isStdoutTerminal := cmdutils.WriterIsTerminal(cmd.OutOrStdout())
	isStderrTerminal := cmdutils.WriterIsTerminal(stderr)
	if shouldSkipUpdateCheck(cmdutils.Version, isStdoutTerminal, isStderrTerminal) {
		return
	}

	cacheDir, err := userCacheDir()
	if err != nil {
		return
	}

	update.Run(cmd.Context(), update.RunOptions{
		CurrentVersion: cmdutils.Version,
		Endpoint:       update.DefaultEndpoint,
		CacheDir:       cacheDir,
		Stderr:         stderr,
		Color:          isStderrTerminal,
	})
}

// shouldSkipUpdateCheck decides whether the version check
// stays quiet. The hint is written to stderr, but stderr
// being a terminal is not enough to conclude that a person
// is watching: a pipe, a redirect, or a shell that sources
// `esdm completion` at startup captures stdout while stderr
// still points at the screen. The hint is only useful when
// someone is actually looking at the command's output, so
// both streams have to be terminals.
func shouldSkipUpdateCheck(version string, isStdoutTerminal, isStderrTerminal bool) bool {
	if version == devVersionMarker {
		return true
	}
	if os.Getenv(update.DisableEnvVar) == disabledValue {
		return true
	}
	if os.Getenv("CI") != "" {
		return true
	}
	if !isStdoutTerminal || !isStderrTerminal {
		return true
	}
	return false
}

func userCacheDir() (string, error) {
	base, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(base, "esdm"), nil
}
