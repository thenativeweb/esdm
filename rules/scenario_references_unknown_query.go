package rules

import (
	"context"
	"fmt"

	"github.com/thenativeweb/esdm/diag"
	"github.com/thenativeweb/esdm/model"
)

const ruleIDScenarioReferencesUnknownQuery = "esdm/gwt/scenario-references-unknown-query"

type scenarioReferencesUnknownQueryRule struct{}

func newScenarioReferencesUnknownQueryRule() *scenarioReferencesUnknownQueryRule {
	return &scenarioReferencesUnknownQueryRule{}
}

func (*scenarioReferencesUnknownQueryRule) Meta() Meta {
	return Meta{
		ID:          ruleIDScenarioReferencesUnknownQuery,
		Extension:   "given-when-then",
		Severity:    diag.SeverityError,
		Description: "Every Query a Read Model Scenario names in its `when` must be declared in the Bounded Context of the Feature.",
	}
}

func (*scenarioReferencesUnknownQueryRule) Check(ctx context.Context, m *model.Model, report diag.Reporter) {
	for _, feature := range sortedByName(m.Extensions.GivenWhenThen.Features) {
		if featureVariant(feature) != featureVariantReadModel {
			continue
		}
		featureName, _ := feature.Name().Text()
		domain := scopeField(feature.Scope(), "domain")
		boundedContext := scopeField(feature.Scope(), "boundedContext")

		for _, scenario := range scenariosOf(feature) {
			scenarioName, _ := scenario.Name().Text()
			query, ok := scenario.When().Field("query").Text()
			if !ok {
				continue
			}
			if _, found := m.LookupQuery(domain, boundedContext, query); found {
				continue
			}
			report.Report(diag.Diagnostic{
				Message:  fmt.Sprintf("feature %q scenario %q when references query %q which is not declared in the model", featureName, scenarioName, query),
				Location: scenario.When().Field("query").Location(),
			})
		}
	}
}
