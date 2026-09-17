package rules

import (
	"context"
	"fmt"

	"github.com/thenativeweb/esdm/diag"
	"github.com/thenativeweb/esdm/model"
)

const ruleIDScenarioReferencesUnknownActor = "esdm/gwt/scenario-references-unknown-actor"

type scenarioReferencesUnknownActorRule struct{}

func newScenarioReferencesUnknownActorRule() *scenarioReferencesUnknownActorRule {
	return &scenarioReferencesUnknownActorRule{}
}

func (*scenarioReferencesUnknownActorRule) Meta() Meta {
	return Meta{
		ID:          ruleIDScenarioReferencesUnknownActor,
		Extension:   "given-when-then",
		Severity:    diag.SeverityError,
		Description: "An Actor named in the `when` of a Scenario must be declared in the Bounded Context of the Feature. A Scenario cannot exercise an Actor the model does not know.",
	}
}

func (*scenarioReferencesUnknownActorRule) Check(ctx context.Context, m *model.Model, report diag.Reporter) {
	for _, feature := range sortedByName(m.Extensions.GivenWhenThen.Features) {
		variant := featureVariant(feature)
		if variant != featureVariantAggregate && variant != featureVariantDynamicConsistencyBoundary {
			continue
		}
		featureName, _ := feature.Name().Text()
		domain := scopeField(feature.Scope(), "domain")
		boundedContext := scopeField(feature.Scope(), "boundedContext")

		for _, scenario := range scenariosOf(feature) {
			scenarioName, _ := scenario.Name().Text()
			actor, ok := scenario.When().Field("actor").Text()
			if !ok {
				continue
			}
			if _, found := m.LookupActor(domain, boundedContext, actor); found {
				continue
			}
			report.Report(diag.Diagnostic{
				Message:  fmt.Sprintf("feature %q scenario %q when references actor %q which is not declared in bounded-context %q", featureName, scenarioName, actor, boundedContext),
				Location: scenario.When().Field("actor").Location(),
			})
		}
	}
}
