package reporter_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thenativeweb/esdm/diag"
	"github.com/thenativeweb/esdm/reporter"
)

func TestJSONFormatter(t *testing.T) {
	t.Run("renders a diagnostic to a stable JSON shape", func(t *testing.T) {
		f := reporter.NewJSONFormatter()

		var buf bytes.Buffer
		err := f.Format(&buf, []diag.Diagnostic{
			{
				RuleID:   "esdm/naming/event-past-tense",
				Severity: diag.SeverityWarning,
				Message:  "m",
				Location: diag.Location{File: "a.esdm.yaml", Line: 3, Column: 7},
			},
		})
		require.NoError(t, err)

		var got []map[string]any
		require.NoError(t, json.Unmarshal(buf.Bytes(), &got))
		require.Len(t, got, 1)

		assert.Equal(t, "esdm/naming/event-past-tense", got[0]["ruleId"])
		assert.Equal(t, "warning", got[0]["severity"])
		assert.Equal(t, "m", got[0]["message"])

		loc, ok := got[0]["location"].(map[string]any)
		require.True(t, ok)
		assert.Equal(t, "a.esdm.yaml", loc["file"])
		assert.Equal(t, float64(3), loc["line"])
		assert.Equal(t, float64(7), loc["column"])
	})

	t.Run("renders an empty slice as an empty JSON array", func(t *testing.T) {
		f := reporter.NewJSONFormatter()

		var buf bytes.Buffer
		err := f.Format(&buf, nil)
		require.NoError(t, err)

		assert.Equal(t, "[]\n", buf.String())
	})

	t.Run("includes related entries when present", func(t *testing.T) {
		f := reporter.NewJSONFormatter()

		var buf bytes.Buffer
		err := f.Format(&buf, []diag.Diagnostic{
			{
				RuleID:   "esdm/structure/duplicate-name",
				Severity: diag.SeverityError,
				Message:  "duplicate name",
				Location: diag.Location{File: "a.esdm.yaml", Line: 1, Column: 1},
				Related: []diag.Related{
					{Message: "first defined here", Location: diag.Location{File: "b.esdm.yaml", Line: 5, Column: 3}},
				},
			},
		})
		require.NoError(t, err)

		var got []map[string]any
		require.NoError(t, json.Unmarshal(buf.Bytes(), &got))
		require.Len(t, got, 1)

		related, ok := got[0]["related"].([]any)
		require.True(t, ok)
		require.Len(t, related, 1)

		relatedEntry, ok := related[0].(map[string]any)
		require.True(t, ok)
		assert.Equal(t, "first defined here", relatedEntry["message"])
	})

	t.Run("omits related field when empty", func(t *testing.T) {
		f := reporter.NewJSONFormatter()

		var buf bytes.Buffer
		err := f.Format(&buf, []diag.Diagnostic{
			{RuleID: "r", Severity: diag.SeverityError, Message: "m"},
		})
		require.NoError(t, err)

		var got []map[string]any
		require.NoError(t, json.Unmarshal(buf.Bytes(), &got))
		require.Len(t, got, 1)

		_, hasRelated := got[0]["related"]
		assert.False(t, hasRelated)
	})
}
