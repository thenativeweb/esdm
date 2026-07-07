package rules

import (
	"context"
	"fmt"

	"github.com/thenativeweb/esdm/diag"
	"github.com/thenativeweb/esdm/model"
)

const ruleIDExternalSystemWithoutDirection = "esdm/modeling/external-system-without-direction"

type externalSystemWithoutDirectionRule struct{}

func newExternalSystemWithoutDirectionRule() *externalSystemWithoutDirectionRule {
	return &externalSystemWithoutDirectionRule{}
}

func (*externalSystemWithoutDirectionRule) Meta() Meta {
	return Meta{
		ID:          ruleIDExternalSystemWithoutDirection,
		Severity:    diag.SeverityWarning,
		Description: "Every external-system must declare a direction (inbound, outbound, or bidirectional); mirrors the JSON Schema's required: [direction] constraint as defense in depth.",
	}
}

func (*externalSystemWithoutDirectionRule) Check(ctx context.Context, m *model.Model, report diag.Reporter) {
	for _, es := range sortedByName(m.ExternalSystems) {
		if es.Direction().Exists() {
			continue
		}
		name, _ := es.Name().Text()
		report.Report(diag.Diagnostic{
			Message:  fmt.Sprintf("external-system %q has no direction", name),
			Location: es.Name().Location(),
		})
	}
}
