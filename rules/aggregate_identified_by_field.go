package rules

import (
	"context"
	"fmt"

	"github.com/thenativeweb/esdm/diag"
	"github.com/thenativeweb/esdm/model"
)

const ruleIDAggregateIdentifiedByField = "esdm/structure/aggregate-identified-by-field"

type aggregateIdentifiedByFieldRule struct{}

func newAggregateIdentifiedByFieldRule() *aggregateIdentifiedByFieldRule {
	return &aggregateIdentifiedByFieldRule{}
}

func (*aggregateIdentifiedByFieldRule) Meta() Meta {
	return Meta{
		ID:          ruleIDAggregateIdentifiedByField,
		Severity:    diag.SeverityError,
		Description: "When aggregate.identifiedBy uses source: state, the named `field` must be declared in the aggregate's state.properties.",
	}
}

func (*aggregateIdentifiedByFieldRule) Check(ctx context.Context, m *model.Model, report diag.Reporter) {
	for _, agg := range sortedByName(m.Aggregates) {
		name, _ := agg.Name().Text()
		identifiedBy := agg.IdentifiedBy()

		source, _ := identifiedBy.Field("source").Text()
		if source != "state" {
			continue
		}
		fieldNode := identifiedBy.Field("field")
		fieldName, ok := fieldNode.Text()
		if !ok {
			continue
		}

		if SchemaHasProperty(agg.State(), fieldName) {
			continue
		}
		report.Report(diag.Diagnostic{
			Message:  fmt.Sprintf("aggregate %q identifiedBy.field %q is not declared in state.properties", name, fieldName),
			Location: fieldNode.Location(),
		})
	}
}
