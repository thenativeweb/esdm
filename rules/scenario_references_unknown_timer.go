package rules

import (
	"context"
	"fmt"

	"github.com/thenativeweb/esdm/ast"
	"github.com/thenativeweb/esdm/diag"
	"github.com/thenativeweb/esdm/model"
)

const ruleIDScenarioReferencesUnknownTimer = "esdm/gwt/scenario-references-unknown-timer"

type scenarioReferencesUnknownTimerRule struct{}

func newScenarioReferencesUnknownTimerRule() *scenarioReferencesUnknownTimerRule {
	return &scenarioReferencesUnknownTimerRule{}
}

func (*scenarioReferencesUnknownTimerRule) Meta() Meta {
	return Meta{
		ID:          ruleIDScenarioReferencesUnknownTimer,
		Severity:    diag.SeverityError,
		Description: "Every timer referenced from a process-manager scenario's when, then.setTimers, or then.cancelTimers must be declared in the targeted process-manager's timers list.",
	}
}

func (r *scenarioReferencesUnknownTimerRule) Check(ctx context.Context, m *model.Model, report diag.Reporter) {
	for _, feature := range sortedByName(m.Extensions.GivenWhenThen.Features) {
		if featureVariant(feature) != featureVariantProcessManager {
			continue
		}
		featureName, _ := feature.Name().Text()
		domain := scopeField(feature.Scope(), "domain")
		processManager := scopeField(feature.Scope(), "processManager")
		pm, ok := m.LookupProcessManager(domain, processManager)
		if !ok {
			continue
		}
		known := timerNamesOf(pm.Field("timers"))

		for _, scenario := range scenariosOf(feature) {
			scenarioName, _ := scenario.Name().Text()
			r.checkTimer(scenario.When().Field("timer"), known, featureName, scenarioName, "when", report)
			for _, t := range scenario.Then().Field("setTimers").Seq() {
				r.checkTimer(t, known, featureName, scenarioName, "then.setTimers", report)
			}
			for _, t := range scenario.Then().Field("cancelTimers").Seq() {
				r.checkTimer(t, known, featureName, scenarioName, "then.cancelTimers", report)
			}
		}
	}
}

func (*scenarioReferencesUnknownTimerRule) checkTimer(node ast.Node, known map[string]bool, featureName, scenarioName, slot string, report diag.Reporter) {
	name, ok := node.Text()
	if !ok {
		return
	}
	if known[name] {
		return
	}
	report.Report(diag.Diagnostic{
		Message:  fmt.Sprintf("feature %q scenario %q %s references timer %q which is not declared on the process-manager", featureName, scenarioName, slot, name),
		Location: node.Location(),
	})
}

func timerNamesOf(timers ast.Node) map[string]bool {
	out := make(map[string]bool)
	for _, t := range timers.Seq() {
		if name, ok := t.Field("name").Text(); ok {
			out[name] = true
		}
	}
	return out
}
