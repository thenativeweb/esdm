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
		assert.True(t, shouldSkipUpdateCheck("(version unavailable)", true))
	})

	t.Run("skips when the disable env var is set to true", func(t *testing.T) {
		t.Setenv(update.DisableEnvVar, "true")
		t.Setenv("CI", "")
		assert.True(t, shouldSkipUpdateCheck("v0.9.0", true))
	})

	t.Run("does not skip when the disable env var is set to anything other than true", func(t *testing.T) {
		t.Setenv(update.DisableEnvVar, "yes")
		t.Setenv("CI", "")
		assert.False(t, shouldSkipUpdateCheck("v0.9.0", true))
	})

	t.Run("skips when the CI env var is set to any non-empty value", func(t *testing.T) {
		t.Setenv(update.DisableEnvVar, "")
		t.Setenv("CI", "true")
		assert.True(t, shouldSkipUpdateCheck("v0.9.0", true))
	})

	t.Run("skips when stderr is not a terminal", func(t *testing.T) {
		t.Setenv(update.DisableEnvVar, "")
		t.Setenv("CI", "")
		assert.True(t, shouldSkipUpdateCheck("v0.9.0", false))
	})

	t.Run("does not skip when no skip condition is met", func(t *testing.T) {
		t.Setenv(update.DisableEnvVar, "")
		t.Setenv("CI", "")
		assert.False(t, shouldSkipUpdateCheck("v0.9.0", true))
	})
}
