package rules

import (
	"context"
	"fmt"

	"github.com/thenativeweb/esdm/diag"
	"github.com/thenativeweb/esdm/model"
)

const ruleIDCrossAggregateCommandEmission = "esdm/modeling/cross-aggregate-command-emission"

type crossAggregateCommandEmissionRule struct{}

func newCrossAggregateCommandEmissionRule() *crossAggregateCommandEmissionRule {
	return &crossAggregateCommandEmissionRule{}
}

func (*crossAggregateCommandEmissionRule) Meta() Meta {
	return Meta{
		ID:          ruleIDCrossAggregateCommandEmission,
		Severity:    diag.SeverityError,
		Description: "A command must only publish events owned by its own aggregate; emitting another aggregate's events crosses an aggregate boundary and breaks ownership.",
	}
}

func (*crossAggregateCommandEmissionRule) Check(ctx context.Context, m *model.Model, report diag.Reporter) {
	for _, cmd := range sortedByName(m.Commands) {
		cmdScope := cmd.Scope()
		cmdDomain := scopeField(cmdScope, "domain")
		cmdBC := scopeField(cmdScope, "boundedContext")
		cmdAggregate := scopeField(cmdScope, "aggregate")
		if cmdAggregate == "" {
			continue
		}
		cmdName, _ := cmd.Name().Text()

		for _, item := range cmd.Publishes().Seq() {
			eventName, ok := item.Text()
			if !ok {
				continue
			}
			// First check whether the published event lives
			// in the command's own aggregate - if so,
			// publishing is in-bounds and there is nothing
			// to flag.
			if _, exists := m.LookupEvent(cmdDomain, cmdBC, cmdAggregate, eventName); exists {
				continue
			}
			// Otherwise look in the same BC for an event of
			// that bare name; if one exists, publishing
			// crosses an aggregate boundary.
			for _, candidate := range m.FindEventsByName(eventName) {
				cScope := candidate.Scope()
				if scopeField(cScope, "domain") != cmdDomain {
					continue
				}
				if scopeField(cScope, "boundedContext") != cmdBC {
					continue
				}
				eventAggregate := scopeField(cScope, "aggregate")
				if eventAggregate == "" || eventAggregate == cmdAggregate {
					continue
				}
				report.Report(diag.Diagnostic{
					Message: fmt.Sprintf(
						"command %q (aggregate %q) publishes event %q owned by aggregate %q",
						cmdName, cmdAggregate, eventName, eventAggregate,
					),
					Location: item.Location(),
				})
				break
			}
		}
	}
}
