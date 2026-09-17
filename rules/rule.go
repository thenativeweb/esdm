package rules

import (
	"context"
	"strings"

	"github.com/thenativeweb/esdm/diag"
	"github.com/thenativeweb/esdm/model"
)

// Meta carries the stable metadata of a rule: its ID,
// default severity, the extension it belongs to, and a
// human-readable description. Severity is per-rule, not
// per-finding, so every Diagnostic produced by the rule
// shares it.
//
// Description is documentation prose: the Linter Rules
// pages of the documentation are generated from it, so it
// addresses the reader of a finding, not the maintainer of
// the rule. It says what the rule checks and why the model
// is better off when the rule does not throw.
//
// Extension names the schema extension a rule belongs to,
// using the extension's directory name under schema/ (for
// example "given-when-then"). It is empty for rules over
// the core schema. The documentation uses it to place a
// rule on the Linter Rules page of its extension; the ID
// alone cannot carry that information, because the
// domain-storytelling rules share their categories with
// the core rules.
type Meta struct {
	ID          string
	Severity    diag.Severity
	Extension   string
	Description string
}

// Category returns the middle segment of the ID, for
// example "structure" for "esdm/structure/ambiguous-name".
func (m Meta) Category() string {
	category, _ := m.splitID()
	return category
}

// Name returns the last segment of the ID, for example
// "ambiguous-name" for "esdm/structure/ambiguous-name".
func (m Meta) Name() string {
	_, name := m.splitID()
	return name
}

// Anchor returns the fragment identifier of the rule's
// entry on its Linter Rules page: the category and the
// name joined by a hyphen. The category is part of the
// anchor so that two rules of different categories can
// never claim the same anchor, however their names evolve.
func (m Meta) Anchor() string {
	return m.Category() + "-" + m.Name()
}

// splitID separates the ID into its category and name.
// The leading "esdm" segment carries no information beyond
// the namespace and is dropped.
func (m Meta) splitID() (category, name string) {
	segments := strings.Split(m.ID, "/")
	if len(segments) != 3 {
		return "", ""
	}
	return segments[1], segments[2]
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
