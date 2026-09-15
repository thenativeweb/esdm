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
	for _, ev := range sortedByName(m.Events) {
		if len(m.PublishersOf(ev)) > 0 {
			continue
		}
		name, _ := ev.Name().Text()
		report.Report(diag.Diagnostic{
			Message:  fmt.Sprintf("event %q is not published by any command", name),
			Location: ev.Name().Location(),
		})
	}
}
