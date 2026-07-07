package rules

import (
	"context"
	"fmt"

	"github.com/thenativeweb/esdm/diag"
	"github.com/thenativeweb/esdm/model"
)

const ruleIDDomainWithoutBoundedContext = "esdm/modeling/domain-without-bounded-context"

type domainWithoutBoundedContextRule struct{}

func newDomainWithoutBoundedContextRule() *domainWithoutBoundedContextRule {
	return &domainWithoutBoundedContextRule{}
}

func (*domainWithoutBoundedContextRule) Meta() Meta {
	return Meta{
		ID:          ruleIDDomainWithoutBoundedContext,
		Severity:    diag.SeverityWarning,
		Description: "A domain should own at least one bounded context; otherwise its presence in the model is decorative.",
	}
}

func (*domainWithoutBoundedContextRule) Check(ctx context.Context, m *model.Model, report diag.Reporter) {
	hasBC := make(map[string]bool)
	for _, bc := range m.BoundedContexts {
		if d := scopeField(bc.Scope(), "domain"); d != "" {
			hasBC[d] = true
		}
	}

	for _, d := range sortedByName(m.Domains) {
		name, _ := d.Name().Text()
		if hasBC[name] {
			continue
		}
		report.Report(diag.Diagnostic{
			Message:  fmt.Sprintf("domain %q has no bounded contexts", name),
			Location: d.Name().Location(),
		})
	}
}
