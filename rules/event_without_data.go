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
		Description: "Every Event must declare a `data` schema. An empty schema is fine and says explicitly that the Event carries no payload beyond the fact that it happened. The schema already requires the field; the rule keeps the requirement in place independently of the schema.",
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
