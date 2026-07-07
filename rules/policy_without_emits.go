package rules

import (
	"context"
	"fmt"

	"github.com/thenativeweb/esdm/diag"
	"github.com/thenativeweb/esdm/model"
)

const ruleIDPolicyWithoutEmits = "esdm/modeling/policy-without-emits"

type policyWithoutEmitsRule struct{}

func newPolicyWithoutEmitsRule() *policyWithoutEmitsRule {
	return &policyWithoutEmitsRule{}
}

func (*policyWithoutEmitsRule) Meta() Meta {
	return Meta{
		ID:          ruleIDPolicyWithoutEmits,
		Severity:    diag.SeverityWarning,
		Description: "Every policy must emit at least one command; mirrors the JSON Schema's required + minItems: 1 constraint as defense in depth. A policy that handles events without emitting anything is observation, which is event-handler territory.",
	}
}

func (*policyWithoutEmitsRule) Check(ctx context.Context, m *model.Model, report diag.Reporter) {
	for _, p := range sortedByName(m.Policies) {
		if len(p.Emits().Seq()) > 0 {
			continue
		}
		name, _ := p.Name().Text()
		report.Report(diag.Diagnostic{
			Message:  fmt.Sprintf("policy %q emits no commands", name),
			Location: p.Name().Location(),
		})
	}
}
