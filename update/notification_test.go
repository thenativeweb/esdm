package update_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thenativeweb/esdm/update"
)

func TestRenderNotification(t *testing.T) {
	t.Run("includes both versions and the upgrade URL", func(t *testing.T) {
		out := update.RenderNotification("v0.9.0", "v0.10.0", false, false)

		assert.Contains(t, out, "0.9.0")
		assert.Contains(t, out, "0.10.0")
		assert.Contains(t, out, "0.9.0 → 0.10.0")
		assert.Contains(t, out, update.UpgradeURL)
	})

	t.Run("strips a leading v from both versions", func(t *testing.T) {
		out := update.RenderNotification("v0.9.0", "v0.10.0", false, false)

		assert.NotContains(t, out, "v0.9.0")
		assert.NotContains(t, out, "v0.10.0")
	})

	t.Run("renders a leading lightning marker", func(t *testing.T) {
		out := update.RenderNotification("0.9.0", "0.10.0", false, false)

		assert.True(t, strings.HasPrefix(out, "⚡"))
	})

	t.Run("includes the disable hint when requested", func(t *testing.T) {
		out := update.RenderNotification("0.9.0", "0.10.0", true, false)

		assert.Contains(t, out, update.DisableEnvVar+"=true")
		assert.Contains(t, out, "once a day")
	})

	t.Run("omits the disable hint when not requested", func(t *testing.T) {
		out := update.RenderNotification("0.9.0", "0.10.0", false, false)

		assert.NotContains(t, out, update.DisableEnvVar)
		assert.NotContains(t, out, "disable")
	})

	t.Run("contains ANSI escape sequences when color is enabled", func(t *testing.T) {
		out := update.RenderNotification("0.9.0", "0.10.0", true, true)

		assert.Contains(t, out, "\x1b[")
	})

	t.Run("contains no ANSI escape sequences when color is disabled", func(t *testing.T) {
		out := update.RenderNotification("0.9.0", "0.10.0", true, false)

		assert.NotContains(t, out, "\x1b[")
	})
}
