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
		Description: "Every Policy must declare at least one Event it handles. The schema already requires this; the rule keeps the requirement in place independently of the schema.",
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
