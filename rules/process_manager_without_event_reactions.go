package rules

import (
	"context"
	"fmt"

	"github.com/thenativeweb/esdm/diag"
	"github.com/thenativeweb/esdm/model"
)

const ruleIDProcessManagerWithoutEventReactions = "esdm/modeling/process-manager-without-event-reactions"

type processManagerWithoutEventReactionsRule struct{}

func newProcessManagerWithoutEventReactionsRule() *processManagerWithoutEventReactionsRule {
	return &processManagerWithoutEventReactionsRule{}
}

func (*processManagerWithoutEventReactionsRule) Meta() Meta {
	return Meta{
		ID:          ruleIDProcessManagerWithoutEventReactions,
		Severity:    diag.SeverityWarning,
		Description: "A Process Manager whose reactions are all timer-based has no event-driven behavior beyond its starting Event. That is usually a gap: a process that only waits is not coordinating anything.",
	}
}

func (*processManagerWithoutEventReactionsRule) Check(ctx context.Context, m *model.Model, report diag.Reporter) {
	for _, pm := range sortedByName(m.ProcessManagers) {
		reactions := pm.Reactions().Seq()
		if len(reactions) == 0 {
			continue
		}

		hasEventReaction := false
		for _, reaction := range reactions {
			if reaction.Field("when").HasField("event") {
				hasEventReaction = true
				break
			}
		}

		if hasEventReaction {
			continue
		}

		name, _ := pm.Name().Text()
		report.Report(diag.Diagnostic{
			Message:  fmt.Sprintf("process-manager %q has only timer-based reactions", name),
			Location: pm.Name().Location(),
		})
	}
}
