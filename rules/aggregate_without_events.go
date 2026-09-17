package rules

import (
	"context"
	"fmt"

	"github.com/thenativeweb/esdm/diag"
	"github.com/thenativeweb/esdm/model"
)

const ruleIDAggregateWithoutEvents = "esdm/modeling/aggregate-without-events"

type aggregateWithoutEventsRule struct{}

func newAggregateWithoutEventsRule() *aggregateWithoutEventsRule {
	return &aggregateWithoutEventsRule{}
}

func (*aggregateWithoutEventsRule) Meta() Meta {
	return Meta{
		ID:          ruleIDAggregateWithoutEvents,
		Severity:    diag.SeverityWarning,
		Description: "An Aggregate should have at least one Event. An Aggregate that records nothing is either unfinished or belongs to a different kind, such as a Read Model or a Value Object.",
	}
}

func (*aggregateWithoutEventsRule) Check(ctx context.Context, m *model.Model, report diag.Reporter) {
	for _, agg := range sortedByName(m.Aggregates) {
		name, _ := agg.Name().Text()
		aggDomain := scopeField(agg.Scope(), "domain")
		aggBC := scopeField(agg.Scope(), "boundedContext")

		hasEvent := false
		for _, ev := range m.Events {
			evScope := ev.Scope()
			if scopeField(evScope, "aggregate") == name &&
				scopeField(evScope, "boundedContext") == aggBC &&
				scopeField(evScope, "domain") == aggDomain {
				hasEvent = true
				break
			}
		}

		if hasEvent {
			continue
		}

		report.Report(diag.Diagnostic{
			Message:  fmt.Sprintf("aggregate %q has no events", name),
			Location: agg.Name().Location(),
		})
	}
}
