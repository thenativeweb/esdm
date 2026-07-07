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
		Description: "Every event-handler must declare at least one side effect; mirrors the JSON Schema's required + minItems: 1 constraint as defense in depth. An event-handler whose only purpose was further state changes should be a process manager.",
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
