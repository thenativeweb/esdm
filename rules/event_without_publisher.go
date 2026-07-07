package rules

import (
	"context"
	"fmt"

	"github.com/thenativeweb/esdm/diag"
	"github.com/thenativeweb/esdm/model"
)

const ruleIDEventWithoutPublisher = "esdm/modeling/event-without-publisher"

type eventWithoutPublisherRule struct{}

func newEventWithoutPublisherRule() *eventWithoutPublisherRule {
	return &eventWithoutPublisherRule{}
}

func (*eventWithoutPublisherRule) Meta() Meta {
	return Meta{
		ID:          ruleIDEventWithoutPublisher,
		Severity:    diag.SeverityWarning,
		Description: "Every event should be published by at least one command; a publisher-less event has no path to come into existence in a running system.",
	}
}

func (*eventWithoutPublisherRule) Check(ctx context.Context, m *model.Model, report diag.Reporter) {
	// An event's identity is (domain, BC, aggregate, name)
	// for aggregate-bound events, and (domain, BC, "",
	// name) for free-standing BC-scoped events emitted by
	// DCB-bound commands. The published set keys on the
	// same composite identity so equally-named events in
	// different scopes are tracked independently.
	published := make(map[string]bool)
	for _, cmd := range m.Commands {
		cmdScope := cmd.Scope()
		domain := scopeField(cmdScope, "domain")
		bc := scopeField(cmdScope, "boundedContext")
		// Aggregate-bound commands publish into their
		// aggregate; DCB-bound commands publish
		// free-standing BC-scoped events (empty aggregate).
		eventAggregate := scopeField(cmdScope, "aggregate")
		for _, item := range cmd.Publishes().Seq() {
			if name, ok := item.Text(); ok {
				published[model.EventKey(domain, bc, eventAggregate, name)] = true
			}
		}
	}

	for _, ev := range sortedByName(m.Events) {
		evScope := ev.Scope()
		domain := scopeField(evScope, "domain")
		bc := scopeField(evScope, "boundedContext")
		aggregate := scopeField(evScope, "aggregate")
		name, _ := ev.Name().Text()
		if published[model.EventKey(domain, bc, aggregate, name)] {
			continue
		}
		report.Report(diag.Diagnostic{
			Message:  fmt.Sprintf("event %q is not published by any command", name),
			Location: ev.Name().Location(),
		})
	}
}
