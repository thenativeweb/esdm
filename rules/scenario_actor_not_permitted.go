package rules

import (
	"context"
	"fmt"

	"github.com/thenativeweb/esdm/diag"
	"github.com/thenativeweb/esdm/model"
)

const ruleIDScenarioActorNotPermitted = "esdm/gwt/scenario-actor-not-permitted"

type scenarioActorNotPermittedRule struct{}

func newScenarioActorNotPermittedRule() *scenarioActorNotPermittedRule {
	return &scenarioActorNotPermittedRule{}
}

func (*scenarioActorNotPermittedRule) Meta() Meta {
	return Meta{
		ID:          ruleIDScenarioActorNotPermitted,
		Severity:    diag.SeverityError,
		Description: "When a scenario names an actor on its when (aggregate or DCB feature), that actor must appear in the targeted command's actors list. A test exercising an unauthorized actor specifies a path the model does not permit.",
	}
}

func (*scenarioActorNotPermittedRule) Check(ctx context.Context, m *model.Model, report diag.Reporter) {
	for _, feature := range sortedByName(m.Extensions.GivenWhenThen.Features) {
		variant := featureVariant(feature)
		if variant != featureVariantAggregate && variant != featureVariantDynamicConsistencyBoundary {
			continue
		}
		featureName, _ := feature.Name().Text()
		scope := feature.Scope()
		domain := scopeField(scope, "domain")
		boundedContext := scopeField(scope, "boundedContext")
		var parent string
		switch variant {
		case featureVariantAggregate:
			parent = scopeField(scope, "aggregate")
		case featureVariantDynamicConsistencyBoundary:
			parent = scopeField(scope, "dynamicConsistencyBoundary")
		}

		for _, scenario := range scenariosOf(feature) {
			scenarioName, _ := scenario.Name().Text()
			actor, ok := scenario.When().Field("actor").Text()
			if !ok {
				continue
			}
			command, ok := scenario.When().Field("command").Text()
			if !ok {
				continue
			}
			cmd, ok := m.LookupCommand(domain, boundedContext, parent, command)
			if !ok {
				continue
			}
			if commandPermitsActor(cmd, actor) {
				continue
			}
			report.Report(diag.Diagnostic{
				Message:  fmt.Sprintf("feature %q scenario %q when names actor %q which command %q does not permit", featureName, scenarioName, actor, command),
				Location: scenario.When().Field("actor").Location(),
			})
		}
	}
}

func commandPermitsActor(command model.CommandView, actor string) bool {
	for _, a := range command.Field("actors").Seq() {
		if name, ok := a.Text(); ok && name == actor {
			return true
		}
	}
	return false
}
