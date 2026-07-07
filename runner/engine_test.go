package runner_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thenativeweb/esdm/diag"
	"github.com/thenativeweb/esdm/model"
	"github.com/thenativeweb/esdm/reporter"
	"github.com/thenativeweb/esdm/rules"
	"github.com/thenativeweb/esdm/runner"
)

// stubRule is a test-only Rule implementation that lets
// us drive the runner's engine behavior precisely.
type stubRule struct {
	id       string
	severity diag.Severity
	check    func(ctx context.Context, m *model.Model, r diag.Reporter)
}

func (s *stubRule) Meta() rules.Meta {
	return rules.Meta{ID: s.id, Severity: s.severity}
}

func (s *stubRule) Check(ctx context.Context, m *model.Model, r diag.Reporter) {
	s.check(ctx, m, r)
}

func TestRunRules(t *testing.T) {
	t.Run("stamps RuleID and Severity onto every reported diagnostic", func(t *testing.T) {
		rule := &stubRule{
			id:       "esdm/x/y",
			severity: diag.SeverityWarning,
			check: func(ctx context.Context, m *model.Model, r diag.Reporter) {
				r.Report(diag.Diagnostic{
					Message:  "something",
					Location: diag.Location{File: "a.esdm.yaml", Line: 1, Column: 1},
				})
			},
		}
		collector := reporter.NewCollector()
		runner.RunRules(context.Background(), []rules.Rule{rule}, model.NewModel(), collector)

		all := collector.All()
		require.Len(t, all, 1)
		assert.Equal(t, "esdm/x/y", all[0].RuleID)
		assert.Equal(t, diag.SeverityWarning, all[0].Severity)
		assert.Equal(t, "something", all[0].Message)
	})

	t.Run("isolates a panicking rule with a system diagnostic", func(t *testing.T) {
		panicky := &stubRule{
			id: "esdm/x/panicky",
			check: func(ctx context.Context, m *model.Model, r diag.Reporter) {
				panic("boom")
			},
		}
		collector := reporter.NewCollector()
		runner.RunRules(context.Background(), []rules.Rule{panicky}, model.NewModel(), collector)

		all := collector.All()
		require.Len(t, all, 1)
		assert.Equal(t, "esdm/system/rule-panic", all[0].RuleID)
		assert.Equal(t, diag.SeverityError, all[0].Severity)
		assert.Contains(t, all[0].Message, "esdm/x/panicky")
		assert.Contains(t, all[0].Message, "boom")
	})

	t.Run("keeps partial findings from a panicking rule", func(t *testing.T) {
		rule := &stubRule{
			id:       "esdm/x/partial",
			severity: diag.SeverityWarning,
			check: func(ctx context.Context, m *model.Model, r diag.Reporter) {
				r.Report(diag.Diagnostic{Message: "first", Location: diag.Location{File: "f", Line: 1, Column: 1}})
				panic("boom")
			},
		}
		collector := reporter.NewCollector()
		runner.RunRules(context.Background(), []rules.Rule{rule}, model.NewModel(), collector)

		all := collector.All()
		require.Len(t, all, 2)

		kinds := map[string]bool{}
		for _, d := range all {
			kinds[d.RuleID] = true
		}

		assert.True(t, kinds["esdm/x/partial"])
		assert.True(t, kinds["esdm/system/rule-panic"])
	})

	t.Run("runs rules concurrently without racing (-race detects)", func(t *testing.T) {
		const ruleCount = 50
		ruleList := make([]rules.Rule, 0, ruleCount)
		for i := 0; i < ruleCount; i++ {
			ruleList = append(ruleList, &stubRule{
				id:       "esdm/x/concurrent",
				severity: diag.SeverityInfo,
				check: func(ctx context.Context, m *model.Model, r diag.Reporter) {
					r.Report(diag.Diagnostic{
						Message:  "concurrent",
						Location: diag.Location{File: "f", Line: i + 1, Column: 1},
					})
				},
			})
		}

		collector := reporter.NewCollector()
		runner.RunRules(context.Background(), ruleList, model.NewModel(), collector)

		assert.Len(t, collector.All(), ruleCount)
	})

	t.Run("one rule's panic does not prevent other rules from running", func(t *testing.T) {
		panicky := &stubRule{
			id: "esdm/x/panicky",
			check: func(ctx context.Context, m *model.Model, r diag.Reporter) {
				panic("boom")
			},
		}
		healthy := &stubRule{
			id:       "esdm/x/healthy",
			severity: diag.SeverityWarning,
			check: func(ctx context.Context, m *model.Model, r diag.Reporter) {
				r.Report(diag.Diagnostic{Message: "ok", Location: diag.Location{File: "f", Line: 1, Column: 1}})
			},
		}

		collector := reporter.NewCollector()
		runner.RunRules(context.Background(), []rules.Rule{panicky, healthy}, model.NewModel(), collector)

		all := collector.All()
		var hasSeenHealthy, hasSeenPanic bool
		for _, d := range all {
			if d.RuleID == "esdm/x/healthy" {
				hasSeenHealthy = true
			}
			if d.RuleID == "esdm/system/rule-panic" {
				hasSeenPanic = true
			}
		}
		assert.True(t, hasSeenHealthy)
		assert.True(t, hasSeenPanic)
	})
}
