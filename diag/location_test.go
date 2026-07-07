package diag_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thenativeweb/esdm/diag"
)

func TestLocationIsZero(t *testing.T) {
	t.Run("returns true for a zero-valued Location", func(t *testing.T) {
		assert.True(t, diag.Location{}.IsZero())
	})

	t.Run("returns false when any field is set", func(t *testing.T) {
		assert.False(t, diag.Location{File: "x.yaml"}.IsZero())
		assert.False(t, diag.Location{Line: 1}.IsZero())
		assert.False(t, diag.Location{Column: 1}.IsZero())
	})
}
