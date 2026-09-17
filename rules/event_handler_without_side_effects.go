package rules

import (
	"context"
	"fmt"

	"github.com/thenativeweb/esdm/diag"
	"github.com/thenativeweb/esdm/model"
)

const ruleIDEventHandlerWithoutSideEffects = "esdm/modeling/event-handler-without-side-effects"

type eventHandlerWithoutSideEffectsRule struct{}

func newEventHandlerWithoutSideEffectsRule() *eventHandlerWithoutSideEffectsRule {
	return &eventHandlerWithoutSideEffectsRule{}
}

func (*eventHandlerWithoutSideEffectsRule) Meta() Meta {
	return Meta{
		ID:          ruleIDEventHandlerWithoutSideEffects,
		Severity:    diag.SeverityWarning,
		Description: "Every Event Handler must declare at least one side effect; reacting to Events with further state changes instead is the job of a Process Manager. The schema already requires this; the rule keeps the requirement in place independently of the schema.",
	}
}

func (*eventHandlerWithoutSideEffectsRule) Check(ctx context.Context, m *model.Model, report diag.Reporter) {
	for _, eh := range sortedByName(m.EventHandlers) {
		if len(eh.SideEffects().Seq()) > 0 {
			continue
		}
		name, _ := eh.Name().Text()
		report.Report(diag.Diagnostic{
			Message:  fmt.Sprintf("event-handler %q has no side effects", name),
			Location: eh.Name().Location(),
		})
	}
}
