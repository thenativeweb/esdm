package rules

import (
	"context"
	"fmt"

	"github.com/thenativeweb/esdm/diag"
	"github.com/thenativeweb/esdm/model"
)

const ruleIDDynamicConsistencyBoundaryWithoutConsults = "esdm/modeling/dynamic-consistency-boundary-without-consults"

type dynamicConsistencyBoundaryWithoutConsultsRule struct{}

func newDynamicConsistencyBoundaryWithoutConsultsRule() *dynamicConsistencyBoundaryWithoutConsultsRule {
	return &dynamicConsistencyBoundaryWithoutConsultsRule{}
}

func (*dynamicConsistencyBoundaryWithoutConsultsRule) Meta() Meta {
	return Meta{
		ID:          ruleIDDynamicConsistencyBoundaryWithoutConsults,
		Severity:    diag.SeverityWarning,
		Description: "Every Dynamic Consistency Boundary must consult at least one Event; the consulted Events are what its decisions are based on. The schema already requires this; the rule keeps the requirement in place independently of the schema.",
	}
}

func (*dynamicConsistencyBoundaryWithoutConsultsRule) Check(ctx context.Context, m *model.Model, report diag.Reporter) {
	for _, dcb := range sortedByName(m.DynamicConsistencyBoundaries) {
		if len(dcb.Consults().Seq()) > 0 {
			continue
		}
		name, _ := dcb.Name().Text()
		report.Report(diag.Diagnostic{
			Message:  fmt.Sprintf("DCB %q has no consults entries", name),
			Location: dcb.Name().Location(),
		})
	}
}
