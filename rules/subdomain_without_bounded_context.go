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
		Description: "A subdomain should list at least one Bounded Context; a subdomain without one names a part of the Domain that nothing in the model fills. The schema already requires this; the rule keeps the requirement in place independently of the schema.",
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
