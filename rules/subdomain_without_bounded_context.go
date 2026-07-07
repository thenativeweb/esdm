package rules

import (
	"context"
	"fmt"

	"github.com/thenativeweb/esdm/diag"
	"github.com/thenativeweb/esdm/model"
)

const ruleIDSubdomainWithoutBoundedContext = "esdm/modeling/subdomain-without-bounded-context"

type subdomainWithoutBoundedContextRule struct{}

func newSubdomainWithoutBoundedContextRule() *subdomainWithoutBoundedContextRule {
	return &subdomainWithoutBoundedContextRule{}
}

func (*subdomainWithoutBoundedContextRule) Meta() Meta {
	return Meta{
		ID:          ruleIDSubdomainWithoutBoundedContext,
		Severity:    diag.SeverityWarning,
		Description: "A subdomain should list at least one bounded context; mirrors the JSON Schema's minItems: 1 constraint as defense in depth so an accidental schema relaxation still surfaces the modeling issue.",
	}
}

func (*subdomainWithoutBoundedContextRule) Check(ctx context.Context, m *model.Model, report diag.Reporter) {
	for _, sub := range sortedByName(m.Subdomains) {
		bcs := sub.BoundedContexts().Seq()
		if len(bcs) > 0 {
			continue
		}
		name, _ := sub.Name().Text()
		report.Report(diag.Diagnostic{
			Message:  fmt.Sprintf("subdomain %q lists no bounded contexts", name),
			Location: sub.Name().Location(),
		})
	}
}
