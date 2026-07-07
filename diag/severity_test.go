package diag_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thenativeweb/esdm/diag"
)

func TestSeverityString(t *testing.T) {
	t.Run("renders known severities", func(t *testing.T) {
		assert.Equal(t, "error", diag.SeverityError.String())
		assert.Equal(t, "warning", diag.SeverityWarning.String())
		assert.Equal(t, "info", diag.SeverityInfo.String())
		assert.Equal(t, "hint", diag.SeverityHint.String())
	})

	t.Run("renders unknown severities as unknown", func(t *testing.T) {
		assert.Equal(t, "unknown", diag.Severity(99).String())
	})
}
