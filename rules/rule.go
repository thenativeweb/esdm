package rules

import (
	"context"

	"github.com/thenativeweb/esdm/diag"
	"github.com/thenativeweb/esdm/model"
)

// Meta carries the stable metadata of a rule: its ID,
// default severity, and a short human-readable
// description. Severity is per-rule, not per-finding, so
// every Diagnostic produced by the rule shares it.
type Meta struct {
	ID          string
	Severity    diag.Severity
	Description string
}

// Rule is the contract every check implements. Check runs
// synchronously against a resolved model; the runner
// wraps it in a goroutine for parallelism and in a
// recover block for panic isolation.
//
// Rules call report.Report with Diagnostics that carry
// Message, Location and (optionally) Related; the runner
// fills in RuleID and Severity from Meta(), so rule code
// does not have to.
type Rule interface {
	Meta() Meta
	Check(ctx context.Context, m *model.Model, report diag.Reporter)
}
