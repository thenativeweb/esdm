package reporter_test

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thenativeweb/esdm/diag"
	"github.com/thenativeweb/esdm/reporter"
)

func TestCollector(t *testing.T) {
	t.Run("returns an empty slice when no diagnostics were reported", func(t *testing.T) {
		c := reporter.NewCollector()
		assert.Empty(t, c.All())
	})

	t.Run("returns reported diagnostics sorted by file, line, column, and rule ID", func(t *testing.T) {
		c := reporter.NewCollector()

		c.Report(diag.Diagnostic{
			RuleID:   "esdm/naming/event-past-tense",
			Location: diag.Location{File: "b.esdm.yaml", Line: 1, Column: 1},
		})
		c.Report(diag.Diagnostic{
			RuleID:   "esdm/structure/duplicate-name",
			Location: diag.Location{File: "a.esdm.yaml", Line: 10, Column: 5},
		})
		c.Report(diag.Diagnostic{
			RuleID:   "esdm/naming/event-past-tense",
			Location: diag.Location{File: "a.esdm.yaml", Line: 5, Column: 1},
		})
		c.Report(diag.Diagnostic{
			RuleID:   "esdm/structure/duplicate-name",
			Location: diag.Location{File: "a.esdm.yaml", Line: 5, Column: 1},
		})

		all := c.All()
		require.Len(t, all, 4)

		assert.Equal(t, "a.esdm.yaml", all[0].Location.File)
		assert.Equal(t, 5, all[0].Location.Line)
		assert.Equal(t, "esdm/naming/event-past-tense", all[0].RuleID)

		assert.Equal(t, "a.esdm.yaml", all[1].Location.File)
		assert.Equal(t, 5, all[1].Location.Line)
		assert.Equal(t, "esdm/structure/duplicate-name", all[1].RuleID)

		assert.Equal(t, "a.esdm.yaml", all[2].Location.File)
		assert.Equal(t, 10, all[2].Location.Line)

		assert.Equal(t, "b.esdm.yaml", all[3].Location.File)
	})

	t.Run("returns a defensive copy so callers cannot mutate internal state", func(t *testing.T) {
		c := reporter.NewCollector()
		c.Report(diag.Diagnostic{RuleID: "esdm/x/y"})

		first := c.All()
		first[0].RuleID = "mutated"

		second := c.All()
		assert.Equal(t, "esdm/x/y", second[0].RuleID)
	})

	t.Run("reports concurrently without racing", func(t *testing.T) {
		c := reporter.NewCollector()

		var wg sync.WaitGroup
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				c.Report(diag.Diagnostic{
					RuleID:   "esdm/x/y",
					Location: diag.Location{File: "f.esdm.yaml", Line: i + 1, Column: 1},
				})
			}(i)
		}
		wg.Wait()

		assert.Len(t, c.All(), 100)
	})

	t.Run("HasErrors returns true when at least one diagnostic has SeverityError", func(t *testing.T) {
		c := reporter.NewCollector()
		c.Report(diag.Diagnostic{Severity: diag.SeverityWarning})
		c.Report(diag.Diagnostic{Severity: diag.SeverityError})

		assert.True(t, c.HasErrors())
	})

	t.Run("HasErrors returns false when no diagnostic has SeverityError", func(t *testing.T) {
		c := reporter.NewCollector()
		c.Report(diag.Diagnostic{Severity: diag.SeverityWarning})
		c.Report(diag.Diagnostic{Severity: diag.SeverityInfo})

		assert.False(t, c.HasErrors())
	})
}
