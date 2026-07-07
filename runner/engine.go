package runner

import (
	"context"
	"fmt"
	"sync"

	"github.com/thenativeweb/esdm/diag"
	"github.com/thenativeweb/esdm/model"
	"github.com/thenativeweb/esdm/rules"
)

// ruleReporter wraps a diag.Reporter for a specific rule.
// It stamps the configured RuleID and Severity onto every
// Diagnostic the rule emits, so rule implementations can
// focus on Message, Location, and Related.
type ruleReporter struct {
	inner    diag.Reporter
	ruleID   string
	severity diag.Severity
}

func (r *ruleReporter) Report(d diag.Diagnostic) {
	d.RuleID = r.ruleID
	d.Severity = r.severity
	r.inner.Report(d)
}

// runRule executes a single rule with panic isolation.
// A panic inside Check is caught and converted into an
// esdm/system/rule-panic diagnostic routed through the
// shared reporter, so one buggy rule cannot abort the
// entire linter run.
func runRule(ctx context.Context, r rules.Rule, m *model.Model, shared diag.Reporter) {
	meta := r.Meta()

	defer func() {
		recovered := recover()
		if recovered == nil {
			return
		}

		shared.Report(diag.Diagnostic{
			RuleID:   "esdm/system/rule-panic",
			Severity: diag.SeverityError,
			Message:  fmt.Sprintf("rule %q panicked: %v", meta.ID, recovered),
		})
	}()

	ruleReport := &ruleReporter{
		inner:    shared,
		ruleID:   meta.ID,
		severity: meta.Severity,
	}

	r.Check(ctx, m, ruleReport)
}

// RunRules executes the given rules concurrently, each in
// its own goroutine. Diagnostics from all rules land in
// the same reporter, which must be concurrency-safe.
func RunRules(ctx context.Context, ruleSet []rules.Rule, m *model.Model, shared diag.Reporter) {
	var wg sync.WaitGroup
	for _, r := range ruleSet {
		wg.Add(1)
		go func() {
			defer wg.Done()
			runRule(ctx, r, m, shared)
		}()
	}
	wg.Wait()
}
