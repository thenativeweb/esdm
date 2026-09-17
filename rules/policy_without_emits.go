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
		Description: "Every Policy must emit at least one Command; a Policy exists to turn Events into Commands, and a Policy that only observes is an Event Handler. The schema already requires this; the rule keeps the requirement in place independently of the schema.",
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
