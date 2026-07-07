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
		Description: "Aggregates should have at least one associated event; an event-less aggregate is either incomplete or miscategorised.",
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
