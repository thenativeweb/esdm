package reporter_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thenativeweb/esdm/diag"
	"github.com/thenativeweb/esdm/reporter"
)

func TestHumanFormatter(t *testing.T) {
	t.Run("renders a diagnostic as a multi-line block", func(t *testing.T) {
		f := reporter.NewHumanFormatter()

		var buf bytes.Buffer
		err := f.Format(&buf, []diag.Diagnostic{
			{
				RuleID:   "esdm/naming/event-past-tense",
				Severity: diag.SeverityWarning,
				Message:  `event name "CreateOrder" is not in past tense`,
				Location: diag.Location{File: "a.esdm.yaml", Line: 3, Column: 7},
			},
		})
		require.NoError(t, err)

		expected := "warning: esdm/naming/event-past-tense\n" +
			"  at a.esdm.yaml:3:7\n" +
			"  event name \"CreateOrder\" is not in past tense\n"
		assert.Equal(t, expected, buf.String())
	})

	t.Run("renders a diagnostic without location with the <internal> placeholder", func(t *testing.T) {
		f := reporter.NewHumanFormatter()

		var buf bytes.Buffer
		err := f.Format(&buf, []diag.Diagnostic{
			{
				RuleID:   "esdm/system/rule-panic",
				Severity: diag.SeverityError,
				Message:  `rule "esdm/naming/event-past-tense" panicked: boom`,
			},
		})
		require.NoError(t, err)

		assert.Contains(t, buf.String(), "error: esdm/system/rule-panic")
		assert.Contains(t, buf.String(), "at <internal>")
		assert.Contains(t, buf.String(), "panicked: boom")
	})

	t.Run("separates multiple diagnostics with a blank line", func(t *testing.T) {
		f := reporter.NewHumanFormatter()

		var buf bytes.Buffer
		err := f.Format(&buf, []diag.Diagnostic{
			{RuleID: "r1", Severity: diag.SeverityError, Message: "m1", Location: diag.Location{File: "a", Line: 1, Column: 1}},
			{RuleID: "r2", Severity: diag.SeverityError, Message: "m2", Location: diag.Location{File: "b", Line: 2, Column: 2}},
		})
		require.NoError(t, err)

		lines := strings.Split(strings.TrimRight(buf.String(), "\n"), "\n")
		var blanks int
		for _, line := range lines {
			if line == "" {
				blanks++
			}
		}
		assert.Equal(t, 1, blanks, "expected exactly one blank line between two diagnostics, got %q", buf.String())
	})

	t.Run("writes nothing for an empty slice", func(t *testing.T) {
		f := reporter.NewHumanFormatter()

		var buf bytes.Buffer
		err := f.Format(&buf, nil)
		require.NoError(t, err)

		assert.Empty(t, buf.String())
	})

	t.Run("renders related entries as indented note lines", func(t *testing.T) {
		f := reporter.NewHumanFormatter()

		var buf bytes.Buffer
		err := f.Format(&buf, []diag.Diagnostic{
			{
				RuleID:   "esdm/structure/unresolved-reference",
				Severity: diag.SeverityError,
				Message:  `unresolved aggregate "ordr"`,
				Location: diag.Location{File: "event.esdm.yaml", Line: 7, Column: 14},
				Related: []diag.Related{
					{
						Message:  `did you mean "order"?`,
						Location: diag.Location{File: "aggregate.esdm.yaml", Line: 3, Column: 7},
					},
				},
			},
		})
		require.NoError(t, err)

		expected := "error: esdm/structure/unresolved-reference\n" +
			"  at event.esdm.yaml:7:14\n" +
			"  unresolved aggregate \"ordr\"\n" +
			"  note: did you mean \"order\"? (aggregate.esdm.yaml:3:7)\n"
		assert.Equal(t, expected, buf.String())
	})

	t.Run("wraps the severity label and rule id in ANSI escape codes when Colors is enabled", func(t *testing.T) {
		f := reporter.NewHumanFormatter()
		f.Colors = true

		var buf bytes.Buffer
		err := f.Format(&buf, []diag.Diagnostic{
			{
				RuleID:   "esdm/x/y",
				Severity: diag.SeverityError,
				Message:  "m",
				Location: diag.Location{File: "a", Line: 1, Column: 1},
			},
		})
		require.NoError(t, err)

		assert.Contains(t, buf.String(), "\x1b[")
		assert.Contains(t, buf.String(), "error")
		assert.Contains(t, buf.String(), "esdm/x/y")
	})

	t.Run("emits no ANSI escape codes when Colors is disabled", func(t *testing.T) {
		f := reporter.NewHumanFormatter()

		var buf bytes.Buffer
		err := f.Format(&buf, []diag.Diagnostic{
			{
				RuleID:   "esdm/x/y",
				Severity: diag.SeverityError,
				Message:  "m",
				Location: diag.Location{File: "a", Line: 1, Column: 1},
			},
		})
		require.NoError(t, err)

		assert.NotContains(t, buf.String(), "\x1b[")
	})
}
