package rules

import (
	"context"
	"fmt"

	"github.com/thenativeweb/esdm/diag"
	"github.com/thenativeweb/esdm/model"
)

const ruleIDBoundedContextWithoutConsistencyUnit = "esdm/modeling/bounded-context-without-consistency-unit"

type boundedContextWithoutConsistencyUnitRule struct{}

func newBoundedContextWithoutConsistencyUnitRule() *boundedContextWithoutConsistencyUnitRule {
	return &boundedContextWithoutConsistencyUnitRule{}
}

func (*boundedContextWithoutConsistencyUnitRule) Meta() Meta {
	return Meta{
		ID:          ruleIDBoundedContextWithoutConsistencyUnit,
		Severity:    diag.SeverityWarning,
		Description: "A bounded context should host at least one consistency unit (aggregate or dynamic-consistency-boundary); otherwise it is a modeling placeholder.",
	}
}

func (*boundedContextWithoutConsistencyUnitRule) Check(ctx context.Context, m *model.Model, report diag.Reporter) {
	hasUnit := make(map[string]bool)
	for _, agg := range m.Aggregates {
		key := scopeField(agg.Scope(), "domain") + "/" + scopeField(agg.Scope(), "boundedContext")
		hasUnit[key] = true
	}
	for _, dcb := range m.DynamicConsistencyBoundaries {
		key := scopeField(dcb.Scope(), "domain") + "/" + scopeField(dcb.Scope(), "boundedContext")
		hasUnit[key] = true
	}

	for _, bc := range sortedByName(m.BoundedContexts) {
		name, _ := bc.Name().Text()
		domain := scopeField(bc.Scope(), "domain")
		key := domain + "/" + name
		if hasUnit[key] {
			continue
		}
		report.Report(diag.Diagnostic{
			Message:  fmt.Sprintf("bounded-context %q has no aggregates or dynamic-consistency-boundaries", name),
			Location: bc.Name().Location(),
		})
	}
}
