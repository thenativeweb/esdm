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
		Description: "Every Actor must declare whether it is `human` or `system`; the type decides which other fields make sense on it. The schema already requires the field; the rule keeps the requirement in place independently of the schema.",
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
