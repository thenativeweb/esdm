package rules

import (
	"context"
	"fmt"

	"github.com/thenativeweb/esdm/ast"
	"github.com/thenativeweb/esdm/diag"
	"github.com/thenativeweb/esdm/model"
)

const ruleIDScenarioRejectionReferencesUnknownInvariant = "esdm/gwt/scenario-rejection-references-unknown-invariant"

type scenarioRejectionReferencesUnknownInvariantRule struct{}

func newScenarioRejectionReferencesUnknownInvariantRule() *scenarioRejectionReferencesUnknownInvariantRule {
	return &scenarioRejectionReferencesUnknownInvariantRule{}
}

func (*scenarioRejectionReferencesUnknownInvariantRule) Meta() Meta {
	return Meta{
		ID:          ruleIDScenarioRejectionReferencesUnknownInvariant,
		Severity:    diag.SeverityError,
		Description: "When a scenario's then.rejection points at a named invariant, that invariant must appear in the targeted unit's invariants list.",
	}
}

func (*scenarioRejectionReferencesUnknownInvariantRule) Check(ctx context.Context, m *model.Model, report diag.Reporter) {
	for _, feature := range sortedByName(m.Extensions.GivenWhenThen.Features) {
		variant := featureVariant(feature)
		if variant != featureVariantAggregate && variant != featureVariantDynamicConsistencyBoundary {
			continue
		}
		featureName, _ := feature.Name().Text()
		scope := feature.Scope()
		domain := scopeField(scope, "domain")
		boundedContext := scopeField(scope, "boundedContext")

		var unitInvariants ast.Node
		switch variant {
		case featureVariantAggregate:
			aggregate := scopeField(scope, "aggregate")
			unit, ok := m.LookupAggregate(domain, boundedContext, aggregate)
			if !ok {
				continue
			}
			unitInvariants = unit.Field("invariants")
		case featureVariantDynamicConsistencyBoundary:
			dcb := scopeField(scope, "dynamicConsistencyBoundary")
			unit, ok := m.LookupDynamicConsistencyBoundary(domain, boundedContext, dcb)
			if !ok {
				continue
			}
			unitInvariants = unit.Field("invariants")
		}
		known := invariantNamesOf(unitInvariants)

		for _, scenario := range scenariosOf(feature) {
			scenarioName, _ := scenario.Name().Text()
			invariant, ok := scenario.Then().Field("rejection").Field("invariant").Text()
			if !ok {
				continue
			}
			if known[invariant] {
				continue
			}
			report.Report(diag.Diagnostic{
				Message:  fmt.Sprintf("feature %q scenario %q rejection references invariant %q which the targeted unit does not declare", featureName, scenarioName, invariant),
				Location: scenario.Then().Field("rejection").Field("invariant").Location(),
			})
		}
	}
}

func invariantNamesOf(invariants ast.Node) map[string]bool {
	out := make(map[string]bool)
	for _, inv := range invariants.Seq() {
		if name, ok := inv.Field("name").Text(); ok {
			out[name] = true
		}
	}
	return out
}
