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
	isTerminal := cmdutils.WriterIsTerminal(stderr)
	if shouldSkipUpdateCheck(cmdutils.Version, isTerminal) {
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
		Color:          isTerminal,
	})
}

func shouldSkipUpdateCheck(version string, isTerminal bool) bool {
	if version == devVersionMarker {
		return true
	}
	if os.Getenv(update.DisableEnvVar) == disabledValue {
		return true
	}
	if os.Getenv("CI") != "" {
		return true
	}
	if !isTerminal {
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
