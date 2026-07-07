package rules

import (
	"context"
	"fmt"

	"github.com/thenativeweb/esdm/diag"
	"github.com/thenativeweb/esdm/model"
)

const ruleIDScenarioWithoutThen = "esdm/gwt/scenario-without-then"

type scenarioWithoutThenRule struct{}

func newScenarioWithoutThenRule() *scenarioWithoutThenRule {
	return &scenarioWithoutThenRule{}
}

func (*scenarioWithoutThenRule) Meta() Meta {
	return Meta{
		ID:          ruleIDScenarioWithoutThen,
		Severity:    diag.SeverityWarning,
		Description: "Every scenario must declare a `then` outcome; mirrors the given-when-then schema's required-fields constraint as defense in depth so an accidental schema relaxation still surfaces the modeling issue.",
	}
}

func (*scenarioWithoutThenRule) Check(ctx context.Context, m *model.Model, report diag.Reporter) {
	for _, feature := range sortedByName(m.Extensions.GivenWhenThen.Features) {
		featureName, _ := feature.Name().Text()
		for _, scenario := range scenariosOf(feature) {
			if scenario.Then().Exists() {
				continue
			}
			scenarioName, _ := scenario.Name().Text()
			report.Report(diag.Diagnostic{
				Message:  fmt.Sprintf("feature %q scenario %q has no then", featureName, scenarioName),
				Location: scenario.Name().Location(),
			})
		}
	}
}
