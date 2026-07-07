package rules

import (
	"context"
	"fmt"

	"github.com/thenativeweb/esdm/diag"
	"github.com/thenativeweb/esdm/model"
)

const ruleIDDynamicConsistencyBoundaryWithoutIdentifiedBy = "esdm/modeling/dynamic-consistency-boundary-without-identified-by"

type dynamicConsistencyBoundaryWithoutIdentifiedByRule struct{}

func newDynamicConsistencyBoundaryWithoutIdentifiedByRule() *dynamicConsistencyBoundaryWithoutIdentifiedByRule {
	return &dynamicConsistencyBoundaryWithoutIdentifiedByRule{}
}

func (*dynamicConsistencyBoundaryWithoutIdentifiedByRule) Meta() Meta {
	return Meta{
		ID:          ruleIDDynamicConsistencyBoundaryWithoutIdentifiedBy,
		Severity:    diag.SeverityWarning,
		Description: "Every DCB must declare at least one identifiedBy entry; mirrors the JSON Schema's required + minItems: 1 constraint as defense in depth.",
	}
}

func (*dynamicConsistencyBoundaryWithoutIdentifiedByRule) Check(ctx context.Context, m *model.Model, report diag.Reporter) {
	for _, dcb := range sortedByName(m.DynamicConsistencyBoundaries) {
		if len(dcb.IdentifiedBy().Seq()) > 0 {
			continue
		}
		name, _ := dcb.Name().Text()
		report.Report(diag.Diagnostic{
			Message:  fmt.Sprintf("DCB %q has no identifiedBy entries", name),
			Location: dcb.Name().Location(),
		})
	}
}
