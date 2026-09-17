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
		Description: "Every Aggregate must declare how its instances are identified, via `identifiedBy`. The schema already requires the field; the rule keeps the requirement in place independently of the schema.",
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
