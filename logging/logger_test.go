package logging_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thenativeweb/esdm/logging"
)

// These are smoke tests that verify the logging functions can be called
// without panicking. They do not verify the actual log output, as the
// logger uses a package-level variable that writes to os.Stdout, making
// it difficult to capture and verify output in tests.

func TestDebug(t *testing.T) {
	t.Run("does not panic when called", func(t *testing.T) {
		assert.NotPanics(t, func() {
			logging.Debug("test debug message")
		})
	})

	t.Run("does not panic when called with arguments", func(t *testing.T) {
		assert.NotPanics(t, func() {
			logging.Debug("test debug message", "key", "value", "count", 42)
		})
	})
}

func TestInfo(t *testing.T) {
	t.Run("does not panic when called", func(t *testing.T) {
		assert.NotPanics(t, func() {
			logging.Info("test info message")
		})
	})

	t.Run("does not panic when called with arguments", func(t *testing.T) {
		assert.NotPanics(t, func() {
			logging.Info("test info message", "key", "value", "count", 42)
		})
	})
}

func TestWarn(t *testing.T) {
	t.Run("does not panic when called", func(t *testing.T) {
		assert.NotPanics(t, func() {
			logging.Warn("test warning message")
		})
	})

	t.Run("does not panic when called with arguments", func(t *testing.T) {
		assert.NotPanics(t, func() {
			logging.Warn("test warning message", "key", "value", "count", 42)
		})
	})
}

func TestError(t *testing.T) {
	t.Run("does not panic when called", func(t *testing.T) {
		assert.NotPanics(t, func() {
			logging.Error("test error message")
		})
	})

	t.Run("does not panic when called with arguments", func(t *testing.T) {
		assert.NotPanics(t, func() {
			logging.Error("test error message", "key", "value", "count", 42)
		})
	})
}
