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
		Description: "Every process manager must declare at least one starting event; mirrors the JSON Schema's required + minItems: 1 constraint as defense in depth.",
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
