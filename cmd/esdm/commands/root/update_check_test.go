package root

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thenativeweb/esdm/update"
)

func TestShouldSkipUpdateCheck(t *testing.T) {
	t.Run("skips when the version is the dev marker", func(t *testing.T) {
		t.Setenv(update.DisableEnvVar, "")
		t.Setenv("CI", "")
		assert.True(t, shouldSkipUpdateCheck("(version unavailable)", true, true))
	})

	t.Run("skips when the disable env var is set to true", func(t *testing.T) {
		t.Setenv(update.DisableEnvVar, "true")
		t.Setenv("CI", "")
		assert.True(t, shouldSkipUpdateCheck("v0.9.0", true, true))
	})

	t.Run("does not skip when the disable env var is set to anything other than true", func(t *testing.T) {
		t.Setenv(update.DisableEnvVar, "yes")
		t.Setenv("CI", "")
		assert.False(t, shouldSkipUpdateCheck("v0.9.0", true, true))
	})

	t.Run("skips when the CI env var is set to any non-empty value", func(t *testing.T) {
		t.Setenv(update.DisableEnvVar, "")
		t.Setenv("CI", "true")
		assert.True(t, shouldSkipUpdateCheck("v0.9.0", true, true))
	})

	t.Run("skips when stderr is not a terminal", func(t *testing.T) {
		t.Setenv(update.DisableEnvVar, "")
		t.Setenv("CI", "")
		assert.True(t, shouldSkipUpdateCheck("v0.9.0", true, false))
	})

	// A shell that loads the completion script via `source <(esdm completion zsh)`
	// captures stdout while stderr stays attached to the terminal. The same holds
	// for any pipe or redirect of a command's output. Nobody is reading the
	// screen in these cases, so the hint must stay quiet.
	t.Run("skips when stdout is not a terminal even though stderr is", func(t *testing.T) {
		t.Setenv(update.DisableEnvVar, "")
		t.Setenv("CI", "")
		assert.True(t, shouldSkipUpdateCheck("v0.9.0", false, true))
	})

	t.Run("does not skip when no skip condition is met", func(t *testing.T) {
		t.Setenv(update.DisableEnvVar, "")
		t.Setenv("CI", "")
		assert.False(t, shouldSkipUpdateCheck("v0.9.0", true, true))
	})
}
