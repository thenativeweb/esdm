package rules

import (
	"context"
	"fmt"

	"github.com/thenativeweb/esdm/diag"
	"github.com/thenativeweb/esdm/model"
)

const ruleIDEventHandlerWithoutHandles = "esdm/modeling/event-handler-without-handles"

type eventHandlerWithoutHandlesRule struct{}

func newEventHandlerWithoutHandlesRule() *eventHandlerWithoutHandlesRule {
	return &eventHandlerWithoutHandlesRule{}
}

func (*eventHandlerWithoutHandlesRule) Meta() Meta {
	return Meta{
		ID:          ruleIDEventHandlerWithoutHandles,
		Severity:    diag.SeverityWarning,
		Description: "Every Event Handler must declare at least one Event it handles. The schema already requires this; the rule keeps the requirement in place independently of the schema.",
	}
}

func (*eventHandlerWithoutHandlesRule) Check(ctx context.Context, m *model.Model, report diag.Reporter) {
	for _, eh := range sortedByName(m.EventHandlers) {
		if len(eh.Handles().Seq()) > 0 {
			continue
		}
		name, _ := eh.Name().Text()
		report.Report(diag.Diagnostic{
			Message:  fmt.Sprintf("event-handler %q handles no events", name),
			Location: eh.Name().Location(),
		})
	}
}
