package rules

import (
	"context"
	"fmt"
	"strings"

	"github.com/thenativeweb/esdm/diag"
	"github.com/thenativeweb/esdm/model"
)

const ruleIDEventNameWithAggregatePrefix = "esdm/modeling/event-name-with-aggregate-prefix"

type eventNameWithAggregatePrefixRule struct{}

func newEventNameWithAggregatePrefixRule() *eventNameWithAggregatePrefixRule {
	return &eventNameWithAggregatePrefixRule{}
}

func (*eventNameWithAggregatePrefixRule) Meta() Meta {
	return Meta{
		ID:          ruleIDEventNameWithAggregatePrefix,
		Severity:    diag.SeverityWarning,
		Description: "An aggregate-bound event's scope already conveys the aggregate; repeating the aggregate's name at the start of the event name is redundant. BC-scoped events (DCB-emitted) are exempt, because no enclosing aggregate provides the context.",
	}
}

// Check walks the event index and flags aggregate-bound
// events whose name begins with their own aggregate's
// name - a redundancy ESDM modeling prefers to avoid.
// Events without an aggregate in their scope (free-
// standing, BC-scoped events emitted by DCB-bound
// commands) are deliberately skipped.
func (*eventNameWithAggregatePrefixRule) Check(ctx context.Context, m *model.Model, report diag.Reporter) {
	for _, ev := range sortedByName(m.Events) {
		aggName := scopeField(ev.Scope(), "aggregate")
		if aggName == "" {
			continue
		}

		name, _ := ev.Name().Text()
		if name != aggName && !strings.HasPrefix(name, aggName+"-") {
			continue
		}

		report.Report(diag.Diagnostic{
			Message:  fmt.Sprintf("event name %q redundantly starts with its aggregate's name %q; the aggregate scope already conveys that context", name, aggName),
			Location: ev.Name().Location(),
		})
	}
}
