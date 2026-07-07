package rules

import (
	"context"
	"fmt"

	"github.com/thenativeweb/esdm/diag"
	"github.com/thenativeweb/esdm/model"
)

const ruleIDProcessManagerWithoutEndsWhen = "esdm/modeling/process-manager-without-ends-when"

type processManagerWithoutEndsWhenRule struct{}

func newProcessManagerWithoutEndsWhenRule() *processManagerWithoutEndsWhenRule {
	return &processManagerWithoutEndsWhenRule{}
}

func (*processManagerWithoutEndsWhenRule) Meta() Meta {
	return Meta{
		ID:          ruleIDProcessManagerWithoutEndsWhen,
		Severity:    diag.SeverityWarning,
		Description: "Every process manager must declare at least one termination condition; mirrors the JSON Schema's required + minItems: 1 constraint on endsWhen as defense in depth. If the schema is later relaxed to allow long-lived process managers, drop this rule alongside.",
	}
}

func (*processManagerWithoutEndsWhenRule) Check(ctx context.Context, m *model.Model, report diag.Reporter) {
	for _, pm := range sortedByName(m.ProcessManagers) {
		if len(pm.EndsWhen().Seq()) > 0 {
			continue
		}
		name, _ := pm.Name().Text()
		report.Report(diag.Diagnostic{
			Message:  fmt.Sprintf("process-manager %q has no endsWhen entries", name),
			Location: pm.Name().Location(),
		})
	}
}
