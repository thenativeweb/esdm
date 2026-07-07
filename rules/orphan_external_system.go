package rules

import (
	"context"
	"fmt"

	"github.com/thenativeweb/esdm/ast"
	"github.com/thenativeweb/esdm/diag"
	"github.com/thenativeweb/esdm/model"
)

const ruleIDOrphanExternalSystem = "esdm/modeling/orphan-external-system"

type orphanExternalSystemRule struct{}

func newOrphanExternalSystemRule() *orphanExternalSystemRule {
	return &orphanExternalSystemRule{}
}

func (*orphanExternalSystemRule) Meta() Meta {
	return Meta{
		ID:          ruleIDOrphanExternalSystem,
		Severity:    diag.SeverityWarning,
		Description: "Every external-system should be referenced from an actor (backedBy), an event-handler side-effect (external-call), or a context-mapping endpoint; otherwise it has no role in the model.",
	}
}

func (*orphanExternalSystemRule) Check(ctx context.Context, m *model.Model, report diag.Reporter) {
	// External systems are domain-scoped: a `backedBy`
	// entry on an actor or an `external-call` side effect
	// on an event-handler resolves within the surrounding
	// entity's domain. Context-mapping endpoints carry
	// their own `domain` field. The used set therefore
	// keys on (domain, bare name).
	used := make(map[string]bool)
	mark := func(domain, name string) {
		used[domain+"/"+name] = true
	}

	for _, actor := range m.Actors {
		domain := scopeField(actor.Scope(), "domain")
		for _, item := range actor.BackedBy().Seq() {
			if name, ok := item.Text(); ok {
				mark(domain, name)
			}
		}
	}
	for _, eh := range m.EventHandlers {
		domain := scopeField(eh.Scope(), "domain")
		for _, se := range eh.SideEffects().Seq() {
			if name, ok := se.Field("externalSystem").Text(); ok {
				mark(domain, name)
			}
		}
	}
	for _, cm := range m.ContextMappings {
		for _, ep := range []ast.Node{
			cm.Customer(), cm.Supplier(), cm.Conformist(),
			cm.Upstream(), cm.Downstream(), cm.Host(),
			cm.Consumer(), cm.Publisher(),
		} {
			if name, ok := ep.Field("externalSystem").Text(); ok {
				domain, _ := ep.Field("domain").Text()
				mark(domain, name)
			}
		}
		for _, p := range cm.Participants().Seq() {
			if name, ok := p.Field("externalSystem").Text(); ok {
				domain, _ := p.Field("domain").Text()
				mark(domain, name)
			}
		}
	}

	for _, es := range sortedByName(m.ExternalSystems) {
		domain := scopeField(es.Scope(), "domain")
		name, _ := es.Name().Text()
		if used[domain+"/"+name] {
			continue
		}
		report.Report(diag.Diagnostic{
			Message:  fmt.Sprintf("external-system %q is not referenced anywhere", name),
			Location: es.Name().Location(),
		})
	}
}
