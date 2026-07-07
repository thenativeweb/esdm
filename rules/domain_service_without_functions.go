package rules

import (
	"context"
	"fmt"

	"github.com/thenativeweb/esdm/diag"
	"github.com/thenativeweb/esdm/model"
)

const ruleIDDomainServiceWithoutFunctions = "esdm/modeling/domain-service-without-functions"

type domainServiceWithoutFunctionsRule struct{}

func newDomainServiceWithoutFunctionsRule() *domainServiceWithoutFunctionsRule {
	return &domainServiceWithoutFunctionsRule{}
}

func (*domainServiceWithoutFunctionsRule) Meta() Meta {
	return Meta{
		ID:          ruleIDDomainServiceWithoutFunctions,
		Severity:    diag.SeverityWarning,
		Description: "Every domain-service must declare at least one function; mirrors the JSON Schema's required + minItems: 1 constraint as defense in depth.",
	}
}

func (*domainServiceWithoutFunctionsRule) Check(ctx context.Context, m *model.Model, report diag.Reporter) {
	for _, ds := range sortedByName(m.DomainServices) {
		if len(ds.Functions().Seq()) > 0 {
			continue
		}
		name, _ := ds.Name().Text()
		report.Report(diag.Diagnostic{
			Message:  fmt.Sprintf("domain-service %q has no functions", name),
			Location: ds.Name().Location(),
		})
	}
}
