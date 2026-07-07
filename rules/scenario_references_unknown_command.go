package rules

import (
	"context"
	"fmt"

	"github.com/thenativeweb/esdm/ast"
	"github.com/thenativeweb/esdm/diag"
	"github.com/thenativeweb/esdm/model"
)

const ruleIDScenarioReferencesUnknownCommand = "esdm/gwt/scenario-references-unknown-command"

type scenarioReferencesUnknownCommandRule struct{}

func newScenarioReferencesUnknownCommandRule() *scenarioReferencesUnknownCommandRule {
	return &scenarioReferencesUnknownCommandRule{}
}

func (*scenarioReferencesUnknownCommandRule) Meta() Meta {
	return Meta{
		ID:          ruleIDScenarioReferencesUnknownCommand,
		Severity:    diag.SeverityError,
		Description: "Every command referenced in a scenario's when (aggregate/DCB features) or in a process-manager scenario's then.emits must be declared in the model. Bare command names resolve through the feature's scope; scoped command references resolve directly through their {boundedContext, aggregate|dynamicConsistencyBoundary, command} triple.",
	}
}

func (r *scenarioReferencesUnknownCommandRule) Check(ctx context.Context, m *model.Model, report diag.Reporter) {
	for _, feature := range sortedByName(m.Extensions.GivenWhenThen.Features) {
		variant := featureVariant(feature)
		if variant == "" {
			continue
		}
		featureName, _ := feature.Name().Text()

		for _, scenario := range scenariosOf(feature) {
			scenarioName, _ := scenario.Name().Text()
			when := scenario.When()

			switch variant {
			case featureVariantAggregate, featureVariantDynamicConsistencyBoundary:
				if name, ok := when.Field("command").Text(); ok {
					r.checkBare(m, feature, variant, name, featureName, scenarioName, "when", when.Field("command").Location(), report)
				}
			case featureVariantProcessManager:
				domain := scopeField(feature.Scope(), "domain")
				for _, emit := range scenario.Then().Field("emits").Seq() {
					r.checkScoped(m, domain, emit, featureName, scenarioName, "then.emits", report)
				}
			}
		}
	}
}

func (*scenarioReferencesUnknownCommandRule) checkBare(m *model.Model, feature model.FeatureView, variant, command, featureName, scenarioName, slot string, loc diag.Location, report diag.Reporter) {
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
	if _, ok := m.LookupCommand(domain, boundedContext, parent, command); ok {
		return
	}
	report.Report(diag.Diagnostic{
		Message:  fmt.Sprintf("feature %q scenario %q %s references command %q which is not declared in the model", featureName, scenarioName, slot, command),
		Location: loc,
	})
}

func (*scenarioReferencesUnknownCommandRule) checkScoped(m *model.Model, domain string, emit ast.Node, featureName, scenarioName, slot string, report diag.Reporter) {
	command, _ := emit.Field("command").Text()
	if command == "" {
		return
	}
	boundedContext, _ := emit.Field("boundedContext").Text()
	parent, _ := emit.Field("aggregate").Text()
	if parent == "" {
		parent, _ = emit.Field("dynamicConsistencyBoundary").Text()
	}
	if _, ok := m.LookupCommand(domain, boundedContext, parent, command); ok {
		return
	}
	report.Report(diag.Diagnostic{
		Message:  fmt.Sprintf("feature %q scenario %q %s references command %q in bounded-context %q which is not declared in the model", featureName, scenarioName, slot, command, boundedContext),
		Location: emit.Field("command").Location(),
	})
}
