package rules

import (
	"context"
	"fmt"

	"github.com/thenativeweb/esdm/ast"
	"github.com/thenativeweb/esdm/diag"
	"github.com/thenativeweb/esdm/model"
)

const ruleIDScenarioReferencesUnknownEvent = "esdm/gwt/scenario-references-unknown-event"

type scenarioReferencesUnknownEventRule struct{}

func newScenarioReferencesUnknownEventRule() *scenarioReferencesUnknownEventRule {
	return &scenarioReferencesUnknownEventRule{}
}

func (*scenarioReferencesUnknownEventRule) Meta() Meta {
	return Meta{
		ID:          ruleIDScenarioReferencesUnknownEvent,
		Severity:    diag.SeverityError,
		Description: "Every event referenced in a scenario's given or then.events must be declared in the model. Bare event names are resolved through the feature's scope; scoped event references resolve directly through their {boundedContext, aggregate?, event} triple.",
	}
}

func (r *scenarioReferencesUnknownEventRule) Check(ctx context.Context, m *model.Model, report diag.Reporter) {
	for _, feature := range sortedByName(m.Extensions.GivenWhenThen.Features) {
		variant := featureVariant(feature)
		if variant == "" {
			continue
		}
		featureName, _ := feature.Name().Text()

		for _, scenario := range scenariosOf(feature) {
			scenarioName, _ := scenario.Name().Text()
			for _, entry := range scenario.Given().Seq() {
				r.checkEntry(m, feature, variant, entry, featureName, scenarioName, "given", report)
			}
			if variant == featureVariantAggregate || variant == featureVariantDynamicConsistencyBoundary {
				for _, entry := range scenario.Then().Field("events").Seq() {
					r.checkEntry(m, feature, variant, entry, featureName, scenarioName, "then.events", report)
				}
			}
		}
	}
}

func (*scenarioReferencesUnknownEventRule) checkEntry(m *model.Model, feature model.FeatureView, variant string, entry ast.Node, featureName, scenarioName, slot string, report diag.Reporter) {
	eventName, _ := entry.Field("event").Text()
	if eventName == "" {
		return
	}
	scope := feature.Scope()
	domain := scopeField(scope, "domain")

	var (
		boundedContext string
		aggregate      string
	)

	if entry.Field("boundedContext").Exists() {
		boundedContext, _ = entry.Field("boundedContext").Text()
		aggregate, _ = entry.Field("aggregate").Text()
	} else {
		switch variant {
		case featureVariantAggregate:
			boundedContext = scopeField(scope, "boundedContext")
			aggregate = scopeField(scope, "aggregate")
		case featureVariantDynamicConsistencyBoundary:
			boundedContext = scopeField(scope, "boundedContext")
			aggregate = ""
		default:
			return
		}
	}

	if _, ok := m.LookupEvent(domain, boundedContext, aggregate, eventName); ok {
		return
	}

	report.Report(diag.Diagnostic{
		Message:  fmt.Sprintf("feature %q scenario %q %s references event %q which is not declared in the model", featureName, scenarioName, slot, eventName),
		Location: entry.Field("event").Location(),
	})
}
