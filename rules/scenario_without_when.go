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
		Extension:   "given-when-then",
		Severity:    diag.SeverityWarning,
		Description: "Every Scenario must declare a `when` trigger; a Scenario without one exercises nothing. The schema already requires the field; the rule keeps the requirement in place independently of the schema.",
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
