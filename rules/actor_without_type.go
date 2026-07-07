package rules

import (
	"context"
	"fmt"

	"github.com/thenativeweb/esdm/diag"
	"github.com/thenativeweb/esdm/model"
)

const ruleIDActorWithoutType = "esdm/modeling/actor-without-type"

type actorWithoutTypeRule struct{}

func newActorWithoutTypeRule() *actorWithoutTypeRule {
	return &actorWithoutTypeRule{}
}

func (*actorWithoutTypeRule) Meta() Meta {
	return Meta{
		ID:          ruleIDActorWithoutType,
		Severity:    diag.SeverityWarning,
		Description: "Every actor must declare a type (human or system); mirrors the JSON Schema's required: [type] constraint as defense in depth.",
	}
}

func (*actorWithoutTypeRule) Check(ctx context.Context, m *model.Model, report diag.Reporter) {
	for _, a := range sortedByName(m.Actors) {
		if a.Type().Exists() {
			continue
		}
		name, _ := a.Name().Text()
		report.Report(diag.Diagnostic{
			Message:  fmt.Sprintf("actor %q has no type", name),
			Location: a.Name().Location(),
		})
	}
}
