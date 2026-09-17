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
		Description: "Every Process Manager must declare at least one termination condition in `endsWhen`; a process that never ends is a design smell the model should not hide. The schema already requires this; the rule keeps the requirement in place independently of the schema.",
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
