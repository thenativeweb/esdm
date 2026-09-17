package rules

import (
	"context"
	"fmt"

	"github.com/thenativeweb/esdm/diag"
	"github.com/thenativeweb/esdm/model"
)

const ruleIDAggregateWithoutCommands = "esdm/modeling/aggregate-without-commands"

type aggregateWithoutCommandsRule struct{}

func newAggregateWithoutCommandsRule() *aggregateWithoutCommandsRule {
	return &aggregateWithoutCommandsRule{}
}

func (*aggregateWithoutCommandsRule) Meta() Meta {
	return Meta{
		ID:          ruleIDAggregateWithoutCommands,
		Severity:    diag.SeverityWarning,
		Description: "An Aggregate should expose at least one Command. Without one, nothing can drive it from the outside, so it never changes state and never publishes an Event.",
	}
}

func (*aggregateWithoutCommandsRule) Check(ctx context.Context, m *model.Model, report diag.Reporter) {
	for _, agg := range sortedByName(m.Aggregates) {
		name, _ := agg.Name().Text()
		aggDomain := scopeField(agg.Scope(), "domain")
		aggBC := scopeField(agg.Scope(), "boundedContext")

		hasCommand := false
		for _, cmd := range m.Commands {
			cmdScope := cmd.Scope()
			if scopeField(cmdScope, "aggregate") == name &&
				scopeField(cmdScope, "boundedContext") == aggBC &&
				scopeField(cmdScope, "domain") == aggDomain {
				hasCommand = true
				break
			}
		}

		if hasCommand {
			continue
		}

		report.Report(diag.Diagnostic{
			Message:  fmt.Sprintf("aggregate %q has no commands", name),
			Location: agg.Name().Location(),
		})
	}
}
