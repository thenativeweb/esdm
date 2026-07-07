package rules

import (
	"context"
	"fmt"
	"strings"

	"github.com/thenativeweb/esdm/ast"
	"github.com/thenativeweb/esdm/diag"
	"github.com/thenativeweb/esdm/model"
)

const ruleIDScenarioWhenMismatchedFeatureScope = "esdm/gwt/scenario-when-mismatched-feature-scope"

type scenarioWhenMismatchedFeatureScopeRule struct{}

func newScenarioWhenMismatchedFeatureScopeRule() *scenarioWhenMismatchedFeatureScopeRule {
	return &scenarioWhenMismatchedFeatureScopeRule{}
}

func (*scenarioWhenMismatchedFeatureScopeRule) Meta() Meta {
	return Meta{
		ID:          ruleIDScenarioWhenMismatchedFeatureScope,
		Severity:    diag.SeverityWarning,
		Description: "A scenario's when shape must match its feature's scope variant: aggregate/DCB features take a command, process-manager features take an event or timer, read-model features take a query. Mirrors the schema's if-then bindings as defense in depth so an accidental schema relaxation still surfaces the mismatch.",
	}
}

func (*scenarioWhenMismatchedFeatureScopeRule) Check(ctx context.Context, m *model.Model, report diag.Reporter) {
	for _, feature := range sortedByName(m.Extensions.GivenWhenThen.Features) {
		variant := featureVariant(feature)
		if variant == "" {
			continue
		}
		featureName, _ := feature.Name().Text()
		expected := expectedWhenSelectors(variant)
		if len(expected) == 0 {
			continue
		}

		for _, scenario := range scenariosOf(feature) {
			scenarioName, _ := scenario.Name().Text()
			actual := actualWhenSelector(scenario.When())
			if actual == "" {
				continue
			}
			if containsString(expected, actual) {
				continue
			}
			report.Report(diag.Diagnostic{
				Message:  fmt.Sprintf("feature %q scenario %q when uses %q but a %s feature expects %s", featureName, scenarioName, actual, variant, joinOr(expected)),
				Location: scenario.When().Location(),
			})
		}
	}
}

// expectedWhenSelectors returns the set of discriminator
// fields a scenario's when may carry for the given feature
// variant.
func expectedWhenSelectors(variant string) []string {
	switch variant {
	case featureVariantAggregate, featureVariantDynamicConsistencyBoundary:
		return []string{"command"}
	case featureVariantProcessManager:
		return []string{"event", "timer"}
	case featureVariantReadModel:
		return []string{"query"}
	}
	return nil
}

// actualWhenSelector returns the first recognized
// discriminator field present on the when node, or "" if
// none is present.
func actualWhenSelector(when ast.Node) string {
	for _, candidate := range []string{"command", "event", "timer", "query"} {
		if when.Field(candidate).Exists() {
			return candidate
		}
	}
	return ""
}

func containsString(haystack []string, needle string) bool {
	for _, h := range haystack {
		if h == needle {
			return true
		}
	}
	return false
}

func joinOr(values []string) string {
	return strings.Join(values, " or ")
}
