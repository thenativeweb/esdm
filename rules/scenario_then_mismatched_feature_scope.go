package rules

import (
	"context"
	"fmt"

	"github.com/thenativeweb/esdm/ast"
	"github.com/thenativeweb/esdm/diag"
	"github.com/thenativeweb/esdm/model"
)

const ruleIDScenarioThenMismatchedFeatureScope = "esdm/gwt/scenario-then-mismatched-feature-scope"

type scenarioThenMismatchedFeatureScopeRule struct{}

func newScenarioThenMismatchedFeatureScopeRule() *scenarioThenMismatchedFeatureScopeRule {
	return &scenarioThenMismatchedFeatureScopeRule{}
}

func (*scenarioThenMismatchedFeatureScopeRule) Meta() Meta {
	return Meta{
		ID:          ruleIDScenarioThenMismatchedFeatureScope,
		Severity:    diag.SeverityWarning,
		Description: "A scenario's then shape must match its feature's scope variant: aggregate/DCB features expect events or rejection, process-manager features expect emits/setTimers/cancelTimers/state/ended, read-model features expect result or readModel. Mirrors the schema's if-then bindings as defense in depth so an accidental schema relaxation still surfaces the mismatch.",
	}
}

func (*scenarioThenMismatchedFeatureScopeRule) Check(ctx context.Context, m *model.Model, report diag.Reporter) {
	for _, feature := range sortedByName(m.Extensions.GivenWhenThen.Features) {
		variant := featureVariant(feature)
		if variant == "" {
			continue
		}
		featureName, _ := feature.Name().Text()
		expected := expectedThenFields(variant)
		if len(expected) == 0 {
			continue
		}

		for _, scenario := range scenariosOf(feature) {
			scenarioName, _ := scenario.Name().Text()
			then := scenario.Then()
			extras := unexpectedFields(then, expected)
			if len(extras) == 0 {
				continue
			}
			report.Report(diag.Diagnostic{
				Message:  fmt.Sprintf("feature %q scenario %q then carries %v but a %s feature expects only %v", featureName, scenarioName, extras, variant, expected),
				Location: then.Location(),
			})
		}
	}
}

func expectedThenFields(variant string) []string {
	switch variant {
	case featureVariantAggregate, featureVariantDynamicConsistencyBoundary:
		return []string{"events", "rejection"}
	case featureVariantProcessManager:
		return []string{"emits", "setTimers", "cancelTimers", "state", "ended"}
	case featureVariantReadModel:
		return []string{"result", "readModel"}
	}
	return nil
}

func unexpectedFields(then ast.Node, expected []string) []string {
	allowed := make(map[string]bool, len(expected))
	for _, e := range expected {
		allowed[e] = true
	}
	var extras []string
	for _, entry := range then.Entries() {
		key, ok := entry.Key.Text()
		if !ok {
			continue
		}
		if !allowed[key] {
			extras = append(extras, key)
		}
	}
	return extras
}
