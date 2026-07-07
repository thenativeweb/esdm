package rules

import (
	"context"
	"fmt"

	"github.com/thenativeweb/esdm/diag"
	"github.com/thenativeweb/esdm/model"
)

const ruleIDHumanActorWithBackedBy = "esdm/modeling/human-actor-with-backed-by"

type humanActorWithBackedByRule struct{}

func newHumanActorWithBackedByRule() *humanActorWithBackedByRule {
	return &humanActorWithBackedByRule{}
}

func (*humanActorWithBackedByRule) Meta() Meta {
	return Meta{
		ID:          ruleIDHumanActorWithBackedBy,
		Severity:    diag.SeverityWarning,
		Description: "A human actor must not declare backedBy; the field documents which external systems implement a system actor's channel and is meaningful only for type: system. Mirrors the JSON Schema's allOf if-then-not constraint as defense in depth.",
	}
}

func (*humanActorWithBackedByRule) Check(ctx context.Context, m *model.Model, report diag.Reporter) {
	for _, a := range sortedByName(m.Actors) {
		actorType, _ := a.Type().Text()
		if actorType != "human" {
			continue
		}
		if !a.BackedBy().Exists() {
			continue
		}
		name, _ := a.Name().Text()
		report.Report(diag.Diagnostic{
			Message:  fmt.Sprintf("human actor %q has a backedBy list, which is only meaningful for type: system", name),
			Location: a.BackedBy().Location(),
		})
	}
}
