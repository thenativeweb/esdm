package rules

import (
	"context"
	"fmt"

	"github.com/thenativeweb/esdm/diag"
	"github.com/thenativeweb/esdm/model"
)

const ruleIDPolicyWithoutHandles = "esdm/modeling/policy-without-handles"

type policyWithoutHandlesRule struct{}

func newPolicyWithoutHandlesRule() *policyWithoutHandlesRule {
	return &policyWithoutHandlesRule{}
}

func (*policyWithoutHandlesRule) Meta() Meta {
	return Meta{
		ID:          ruleIDPolicyWithoutHandles,
		Severity:    diag.SeverityWarning,
		Description: "Every policy must declare at least one event it handles; mirrors the JSON Schema's required + minItems: 1 constraint as defense in depth.",
	}
}

func (*policyWithoutHandlesRule) Check(ctx context.Context, m *model.Model, report diag.Reporter) {
	for _, p := range sortedByName(m.Policies) {
		if len(p.Handles().Seq()) > 0 {
			continue
		}
		name, _ := p.Name().Text()
		report.Report(diag.Diagnostic{
			Message:  fmt.Sprintf("policy %q handles no events", name),
			Location: p.Name().Location(),
		})
	}
}
