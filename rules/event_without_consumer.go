package rules

import (
	"context"
	"fmt"

	"github.com/thenativeweb/esdm/ast"
	"github.com/thenativeweb/esdm/diag"
	"github.com/thenativeweb/esdm/model"
)

const ruleIDEventWithoutConsumer = "esdm/modeling/event-without-consumer"

type eventWithoutConsumerRule struct{}

func newEventWithoutConsumerRule() *eventWithoutConsumerRule {
	return &eventWithoutConsumerRule{}
}

func (*eventWithoutConsumerRule) Meta() Meta {
	return Meta{
		ID:          ruleIDEventWithoutConsumer,
		Severity:    diag.SeverityWarning,
		Description: "Events should be consumed somewhere (event-handler, policy, process-manager, read-model, or DCB); a never-consumed event is usually a modeling leftover.",
	}
}

func (*eventWithoutConsumerRule) Check(ctx context.Context, m *model.Model, report diag.Reporter) {
	// Consumer references carry a {boundedContext,
	// [aggregate], event} triple; the consumer's own
	// scope provides the domain. A consumed event's
	// identity is therefore (domain, BC, aggregate, name)
	// for aggregate-bound events, or (domain, BC, "",
	// name) for free-standing BC-scoped events. The
	// consumed set keys on that same composite identity
	// so equally-named events in different scopes track
	// independently.
	consumed := make(map[string]bool)
	mark := func(domain string, ref ast.Node) {
		bc, _ := ref.Field("boundedContext").Text()
		name, _ := ref.Field("event").Text()
		if bc == "" || name == "" {
			return
		}
		var aggregate string
		if aggNode := ref.Field("aggregate"); aggNode.Exists() {
			aggregate, _ = aggNode.Text()
		}
		consumed[model.EventKey(domain, bc, aggregate, name)] = true
	}

	for _, eh := range m.EventHandlers {
		domain := scopeField(eh.Scope(), "domain")
		for _, ref := range eh.Handles().Seq() {
			mark(domain, ref)
		}
	}
	for _, p := range m.Policies {
		domain := scopeField(p.Scope(), "domain")
		for _, ref := range p.Handles().Seq() {
			mark(domain, ref)
		}
	}
	for _, pm := range m.ProcessManagers {
		domain := scopeField(pm.Scope(), "domain")
		for _, ref := range pm.StartsWhen().Seq() {
			mark(domain, ref)
		}
		for _, reaction := range pm.Reactions().Seq() {
			when := reaction.Field("when")
			if when.HasField("event") {
				mark(domain, when)
			}
		}
	}
	for _, rm := range m.ReadModels {
		domain := scopeField(rm.Scope(), "domain")
		for _, proj := range rm.Projections().Seq() {
			mark(domain, proj)
		}
	}
	for _, dcb := range m.DynamicConsistencyBoundaries {
		domain := scopeField(dcb.Scope(), "domain")
		for _, c := range dcb.Consults().Seq() {
			mark(domain, c)
		}
	}

	for _, ev := range sortedByName(m.Events) {
		evScope := ev.Scope()
		domain := scopeField(evScope, "domain")
		bc := scopeField(evScope, "boundedContext")
		aggregate := scopeField(evScope, "aggregate")
		name, _ := ev.Name().Text()
		if consumed[model.EventKey(domain, bc, aggregate, name)] {
			continue
		}
		report.Report(diag.Diagnostic{
			Message:  fmt.Sprintf("event %q has no consumer", name),
			Location: ev.Name().Location(),
		})
	}
}
