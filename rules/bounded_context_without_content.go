package rules

import (
	"context"
	"fmt"

	"github.com/thenativeweb/esdm/ast"
	"github.com/thenativeweb/esdm/diag"
	"github.com/thenativeweb/esdm/model"
)

const ruleIDBoundedContextWithoutContent = "esdm/modeling/bounded-context-without-content"

type boundedContextWithoutContentRule struct{}

func newBoundedContextWithoutContentRule() *boundedContextWithoutContentRule {
	return &boundedContextWithoutContentRule{}
}

func (*boundedContextWithoutContentRule) Meta() Meta {
	return Meta{
		ID:          ruleIDBoundedContextWithoutContent,
		Severity:    diag.SeverityWarning,
		Description: "A Bounded Context should host at least one Aggregate, Dynamic Consistency Boundary, or Read Model. These are the kinds that carry behavior or derived state; a Bounded Context holding only supporting elements such as Entities, Value Objects, or Actors is a placeholder.",
	}
}

// markBoundedContexts records the bounded context each entity of the
// given map lives in. Aggregates, dynamic-consistency-boundaries and
// read-models are the three kinds that count as content, and they
// share no common view type beyond their scope, so the marking is
// generic over it.
func markBoundedContexts[V interface{ Scope() ast.Node }](entities map[string]V, hasContent map[string]bool) {
	for _, entity := range entities {
		key := scopeField(entity.Scope(), "domain") + "/" + scopeField(entity.Scope(), "boundedContext")
		hasContent[key] = true
	}
}

func (*boundedContextWithoutContentRule) Check(ctx context.Context, m *model.Model, report diag.Reporter) {
	// Read models count alongside the two consistency units because a
	// projection-only context is a deliberate modeling choice: it derives
	// every view from events owned by other bounded contexts and needs no
	// write side of its own.
	hasContent := make(map[string]bool)
	markBoundedContexts(m.Aggregates, hasContent)
	markBoundedContexts(m.DynamicConsistencyBoundaries, hasContent)
	markBoundedContexts(m.ReadModels, hasContent)

	for _, bc := range sortedByName(m.BoundedContexts) {
		name, _ := bc.Name().Text()
		domain := scopeField(bc.Scope(), "domain")
		key := domain + "/" + name
		if hasContent[key] {
			continue
		}
		report.Report(diag.Diagnostic{
			Message:  fmt.Sprintf("bounded-context %q has no aggregates, dynamic-consistency-boundaries or read-models", name),
			Location: bc.Name().Location(),
		})
	}
}
