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
		Description: "Every event-handler must declare at least one event it handles; mirrors the JSON Schema's required + minItems: 1 constraint as defense in depth.",
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
