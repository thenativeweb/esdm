package hint_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thenativeweb/esdm/hint"
)

func TestLevenshtein(t *testing.T) {
	t.Run("returns 0 for identical strings", func(t *testing.T) {
		assert.Equal(t, 0, hint.Levenshtein("order", "order"))
	})

	t.Run("returns the length when one side is empty", func(t *testing.T) {
		assert.Equal(t, 5, hint.Levenshtein("", "order"))
		assert.Equal(t, 5, hint.Levenshtein("order", ""))
	})

	t.Run("counts single-character substitutions", func(t *testing.T) {
		assert.Equal(t, 1, hint.Levenshtein("order", "ordor"))
	})

	t.Run("counts insertions and deletions", func(t *testing.T) {
		assert.Equal(t, 1, hint.Levenshtein("order", "orders"))
		assert.Equal(t, 1, hint.Levenshtein("orders", "order"))
	})

	t.Run("operates on rune boundaries, not bytes", func(t *testing.T) {
		assert.Equal(t, 1, hint.Levenshtein("café", "cafe"))
	})
}

func TestClosestMatch(t *testing.T) {
	t.Run("returns ok=false for an empty candidate list", func(t *testing.T) {
		_, _, ok := hint.ClosestMatch("anything", nil)
		assert.False(t, ok)
	})

	t.Run("picks the candidate with the smallest distance", func(t *testing.T) {
		best, distance, ok := hint.ClosestMatch("ordr", []string{"order", "shipment", "invoice"})
		require.True(t, ok)
		assert.Equal(t, "order", best)
		assert.Equal(t, 1, distance)
	})
}

func TestSuggestionThreshold(t *testing.T) {
	t.Run("short names tolerate at least 2 edits", func(t *testing.T) {
		assert.Equal(t, 2, hint.SuggestionThreshold("ab"))
		assert.Equal(t, 2, hint.SuggestionThreshold("order"))
	})

	t.Run("long names tolerate up to a third of their length", func(t *testing.T) {
		assert.Equal(t, 3, hint.SuggestionThreshold("nine-runes"))
		assert.Equal(t, 4, hint.SuggestionThreshold("twelve-runes"))
	})
}

func TestBest(t *testing.T) {
	t.Run("returns the closest candidate when within threshold", func(t *testing.T) {
		best, ok := hint.Best("ordr", []string{"order", "shipment"})
		require.True(t, ok)
		assert.Equal(t, "order", best)
	})

	t.Run("returns false when no candidate is close enough", func(t *testing.T) {
		_, ok := hint.Best("wildly-different", []string{"short", "also-short"})
		assert.False(t, ok)
	})

	t.Run("returns false when the input matches a candidate exactly", func(t *testing.T) {
		_, ok := hint.Best("order", []string{"order", "shipment"})
		assert.False(t, ok)
	})

	t.Run("returns false for an empty candidate list", func(t *testing.T) {
		_, ok := hint.Best("anything", nil)
		assert.False(t, ok)
	})
}
