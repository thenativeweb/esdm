package schema_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thenativeweb/esdm/schema"
)

func TestCompareRevisions(t *testing.T) {
	cases := []struct {
		name string
		a, b string
		want int
	}{
		{"equal", "1.0.0", "1.0.0", 0},
		{"older patch", "1.0.0", "1.0.1", -1},
		{"newer patch", "1.0.2", "1.0.1", 1},
		{"older minor", "1.0.9", "1.1.0", -1},
		{"newer minor", "1.2.0", "1.1.9", 1},
		{"older major", "1.9.9", "2.0.0", -1},
		{"newer major", "2.0.0", "1.9.9", 1},
		{"multi-digit components", "1.10.0", "1.9.0", 1},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := schema.CompareRevisions(c.a, c.b)
			require.NoError(t, err)
			assert.Equal(t, c.want, got)
		})
	}

	t.Run("rejects non-semver inputs", func(t *testing.T) {
		_, err := schema.CompareRevisions("1.0", "1.0.0")
		assert.Error(t, err)

		_, err = schema.CompareRevisions("1.0.0", "1.x.0")
		assert.Error(t, err)
	})
}
