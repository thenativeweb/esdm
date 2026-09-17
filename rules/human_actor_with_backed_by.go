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
		Description: "A `human` Actor must not declare `backedBy`. The field names the External Systems that implement the channel of a `system` Actor and has no meaning for a person. The schema already forbids the combination; the rule keeps the restriction in place independently of the schema.",
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
