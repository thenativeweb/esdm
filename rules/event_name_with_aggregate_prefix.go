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
		Description: "An Event bound to an Aggregate should not repeat the name of the Aggregate at the start of its own name. The scope already conveys the Aggregate, and ESDM composes the full name from Aggregate and Event wherever it is needed, so `book-registered` on the Aggregate `book` would read as `BookBookRegistered`. Events scoped to a Bounded Context are exempt, because no enclosing Aggregate provides the context.",
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
