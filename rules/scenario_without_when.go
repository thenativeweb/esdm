package rules

import (
	"context"
	"fmt"

	"github.com/thenativeweb/esdm/diag"
	"github.com/thenativeweb/esdm/model"
)

const ruleIDScenarioWithoutWhen = "esdm/gwt/scenario-without-when"

type scenarioWithoutWhenRule struct{}

func newScenarioWithoutWhenRule() *scenarioWithoutWhenRule {
	return &scenarioWithoutWhenRule{}
}

func (*scenarioWithoutWhenRule) Meta() Meta {
	return Meta{
		ID:          ruleIDScenarioWithoutWhen,
		Severity:    diag.SeverityWarning,
		Description: "Every scenario must declare a `when` trigger; mirrors the given-when-then schema's required-fields constraint as defense in depth so an accidental schema relaxation still surfaces the modeling issue.",
	}
}

func (*scenarioWithoutWhenRule) Check(ctx context.Context, m *model.Model, report diag.Reporter) {
	for _, feature := range sortedByName(m.Extensions.GivenWhenThen.Features) {
		featureName, _ := feature.Name().Text()
		for _, scenario := range scenariosOf(feature) {
			if scenario.When().Exists() {
				continue
			}
			scenarioName, _ := scenario.Name().Text()
			report.Report(diag.Diagnostic{
				Message:  fmt.Sprintf("feature %q scenario %q has no when", featureName, scenarioName),
				Location: scenario.Name().Location(),
			})
		}
	}
}
