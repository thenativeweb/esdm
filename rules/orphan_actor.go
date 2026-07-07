package rules

import (
	"context"
	"fmt"

	"github.com/thenativeweb/esdm/diag"
	"github.com/thenativeweb/esdm/model"
)

const ruleIDOrphanActor = "esdm/modeling/orphan-actor"

type orphanActorRule struct{}

func newOrphanActorRule() *orphanActorRule {
	return &orphanActorRule{}
}

func (*orphanActorRule) Meta() Meta {
	return Meta{
		ID:          ruleIDOrphanActor,
		Severity:    diag.SeverityWarning,
		Description: "Every actor should be named by at least one command.actors or query.actors entry; an actor that nothing references is decorative.",
	}
}

func (*orphanActorRule) Check(ctx context.Context, m *model.Model, report diag.Reporter) {
	// Actors are bounded-context-scoped: a `command.actors`
	// or `query.actors` entry resolves within the command's
	// or query's own BC. Two actors of the same bare name
	// in different BCs are distinct entities, so the
	// referenced set keys on (domain, BC, bare name).
	referenced := make(map[string]bool)
	addRef := func(domain, bc, name string) {
		referenced[domain+"/"+bc+"/"+name] = true
	}

	for _, cmd := range m.Commands {
		domain := scopeField(cmd.Scope(), "domain")
		bc := scopeField(cmd.Scope(), "boundedContext")
		for _, item := range cmd.Actors().Seq() {
			if name, ok := item.Text(); ok {
				addRef(domain, bc, name)
			}
		}
	}
	for _, q := range m.Queries {
		domain := scopeField(q.Scope(), "domain")
		bc := scopeField(q.Scope(), "boundedContext")
		for _, item := range q.Actors().Seq() {
			if name, ok := item.Text(); ok {
				addRef(domain, bc, name)
			}
		}
	}

	for _, a := range sortedByName(m.Actors) {
		domain := scopeField(a.Scope(), "domain")
		bc := scopeField(a.Scope(), "boundedContext")
		name, _ := a.Name().Text()
		if referenced[domain+"/"+bc+"/"+name] {
			continue
		}
		report.Report(diag.Diagnostic{
			Message:  fmt.Sprintf("actor %q is not referenced by any command or query", name),
			Location: a.Name().Location(),
		})
	}
}
