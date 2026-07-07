package update_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thenativeweb/esdm/update"
)

func TestIsNewer(t *testing.T) {
	t.Run("returns true when latest is greater than current", func(t *testing.T) {
		isNewer, err := update.IsNewer("v0.9.0", "v0.10.0")
		require.NoError(t, err)
		assert.True(t, isNewer)
	})

	t.Run("returns false when latest equals current", func(t *testing.T) {
		isNewer, err := update.IsNewer("v0.9.0", "v0.9.0")
		require.NoError(t, err)
		assert.False(t, isNewer)
	})

	t.Run("returns false when latest is older than current", func(t *testing.T) {
		isNewer, err := update.IsNewer("v0.10.0", "v0.9.0")
		require.NoError(t, err)
		assert.False(t, isNewer)
	})

	t.Run("accepts versions without a v-prefix", func(t *testing.T) {
		isNewer, err := update.IsNewer("0.9.0", "0.10.0")
		require.NoError(t, err)
		assert.True(t, isNewer)
	})

	t.Run("accepts mixed prefixed and bare inputs", func(t *testing.T) {
		isNewer, err := update.IsNewer("v0.9.0", "0.10.0")
		require.NoError(t, err)
		assert.True(t, isNewer)
	})

	t.Run("returns an error for an invalid current version", func(t *testing.T) {
		_, err := update.IsNewer("(version unavailable)", "0.10.0")
		require.Error(t, err)
	})

	t.Run("returns an error for an invalid latest version", func(t *testing.T) {
		_, err := update.IsNewer("0.9.0", "garbage")
		require.Error(t, err)
	})
}

func TestStripVPrefix(t *testing.T) {
	t.Run("removes a leading v", func(t *testing.T) {
		assert.Equal(t, "0.9.0", update.StripVPrefix("v0.9.0"))
	})

	t.Run("leaves a bare version unchanged", func(t *testing.T) {
		assert.Equal(t, "0.9.0", update.StripVPrefix("0.9.0"))
	})

	t.Run("does not remove other leading characters", func(t *testing.T) {
		assert.Equal(t, "x0.9.0", update.StripVPrefix("x0.9.0"))
	})
}
