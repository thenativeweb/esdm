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
		Extension:   "given-when-then",
		Severity:    diag.SeverityWarning,
		Description: "Every Scenario must declare a `then` outcome; a Scenario without one asserts nothing. The schema already requires the field; the rule keeps the requirement in place independently of the schema.",
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
