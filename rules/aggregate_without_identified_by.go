package rules

import (
	"context"
	"fmt"

	"github.com/thenativeweb/esdm/diag"
	"github.com/thenativeweb/esdm/model"
)

const ruleIDAggregateWithoutIdentifiedBy = "esdm/modeling/aggregate-without-identified-by"

type aggregateWithoutIdentifiedByRule struct{}

func newAggregateWithoutIdentifiedByRule() *aggregateWithoutIdentifiedByRule {
	return &aggregateWithoutIdentifiedByRule{}
}

func (*aggregateWithoutIdentifiedByRule) Meta() Meta {
	return Meta{
		ID:          ruleIDAggregateWithoutIdentifiedBy,
		Severity:    diag.SeverityWarning,
		Description: "Every aggregate must declare an identifier strategy via identifiedBy; mirrors the JSON Schema's required: [identifiedBy] constraint as defense in depth.",
	}
}

func (*aggregateWithoutIdentifiedByRule) Check(ctx context.Context, m *model.Model, report diag.Reporter) {
	for _, agg := range sortedByName(m.Aggregates) {
		if agg.IdentifiedBy().Exists() {
			continue
		}
		name, _ := agg.Name().Text()
		report.Report(diag.Diagnostic{
			Message:  fmt.Sprintf("aggregate %q has no identifiedBy", name),
			Location: agg.Name().Location(),
		})
	}
}
