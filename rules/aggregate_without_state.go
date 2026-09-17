package rules

import (
	"context"
	"fmt"

	"github.com/thenativeweb/esdm/diag"
	"github.com/thenativeweb/esdm/model"
)

const ruleIDAggregateWithoutState = "esdm/modeling/aggregate-without-state"

type aggregateWithoutStateRule struct{}

func newAggregateWithoutStateRule() *aggregateWithoutStateRule {
	return &aggregateWithoutStateRule{}
}

func (*aggregateWithoutStateRule) Meta() Meta {
	return Meta{
		ID:          ruleIDAggregateWithoutState,
		Severity:    diag.SeverityWarning,
		Description: "Every Aggregate must declare a `state` schema. An empty schema (`type: object`) is fine and says explicitly that the Aggregate has no observable state. The schema already requires the field; the rule keeps the requirement in place independently of the schema.",
	}
}

func (*aggregateWithoutStateRule) Check(ctx context.Context, m *model.Model, report diag.Reporter) {
	for _, agg := range sortedByName(m.Aggregates) {
		if agg.State().Exists() {
			continue
		}
		name, _ := agg.Name().Text()
		report.Report(diag.Diagnostic{
			Message:  fmt.Sprintf("aggregate %q has no state schema", name),
			Location: agg.Name().Location(),
		})
	}
}
