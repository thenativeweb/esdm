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
		Description: "Every DCB must declare at least one consulted event; mirrors the JSON Schema's required + minItems: 1 constraint as defense in depth.",
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
