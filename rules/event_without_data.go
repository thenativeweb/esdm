package rules

import (
	"context"
	"fmt"

	"github.com/thenativeweb/esdm/diag"
	"github.com/thenativeweb/esdm/model"
)

const ruleIDEventWithoutData = "esdm/modeling/event-without-data"

type eventWithoutDataRule struct{}

func newEventWithoutDataRule() *eventWithoutDataRule {
	return &eventWithoutDataRule{}
}

func (*eventWithoutDataRule) Meta() Meta {
	return Meta{
		ID:          ruleIDEventWithoutData,
		Severity:    diag.SeverityWarning,
		Description: "Every event must declare a data schema (an empty schema is allowed for tag-only events); mirrors the JSON Schema's required: [data] constraint as defense in depth so 'no payload' is expressed deliberately.",
	}
}

func (*eventWithoutDataRule) Check(ctx context.Context, m *model.Model, report diag.Reporter) {
	for _, ev := range sortedByName(m.Events) {
		if ev.Data().Exists() {
			continue
		}
		name, _ := ev.Name().Text()
		report.Report(diag.Diagnostic{
			Message:  fmt.Sprintf("event %q has no data schema", name),
			Location: ev.Name().Location(),
		})
	}
}
