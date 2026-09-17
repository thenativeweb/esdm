package rules

import (
	"context"
	"fmt"

	"github.com/thenativeweb/esdm/diag"
	"github.com/thenativeweb/esdm/model"
)

const ruleIDProcessManagerWithoutStartsWhen = "esdm/modeling/process-manager-without-starts-when"

type processManagerWithoutStartsWhenRule struct{}

func newProcessManagerWithoutStartsWhenRule() *processManagerWithoutStartsWhenRule {
	return &processManagerWithoutStartsWhenRule{}
}

func (*processManagerWithoutStartsWhenRule) Meta() Meta {
	return Meta{
		ID:          ruleIDProcessManagerWithoutStartsWhen,
		Severity:    diag.SeverityWarning,
		Description: "Every Process Manager must declare at least one starting Event in `startsWhen`. The schema already requires this; the rule keeps the requirement in place independently of the schema.",
	}
}

func (*processManagerWithoutStartsWhenRule) Check(ctx context.Context, m *model.Model, report diag.Reporter) {
	for _, pm := range sortedByName(m.ProcessManagers) {
		if len(pm.StartsWhen().Seq()) > 0 {
			continue
		}
		name, _ := pm.Name().Text()
		report.Report(diag.Diagnostic{
			Message:  fmt.Sprintf("process-manager %q has no startsWhen entries", name),
			Location: pm.Name().Location(),
		})
	}
}
