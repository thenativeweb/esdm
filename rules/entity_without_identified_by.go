package rules

import (
	"context"
	"fmt"

	"github.com/thenativeweb/esdm/diag"
	"github.com/thenativeweb/esdm/model"
)

const ruleIDEntityWithoutIdentifiedBy = "esdm/modeling/entity-without-identified-by"

type entityWithoutIdentifiedByRule struct{}

func newEntityWithoutIdentifiedByRule() *entityWithoutIdentifiedByRule {
	return &entityWithoutIdentifiedByRule{}
}

func (*entityWithoutIdentifiedByRule) Meta() Meta {
	return Meta{
		ID:          ruleIDEntityWithoutIdentifiedBy,
		Severity:    diag.SeverityWarning,
		Description: "Every Entity must declare how its instances are identified, via `identifiedBy`; identity is what distinguishes an Entity from a Value Object. The schema already requires the field; the rule keeps the requirement in place independently of the schema.",
	}
}

func (*entityWithoutIdentifiedByRule) Check(ctx context.Context, m *model.Model, report diag.Reporter) {
	for _, e := range sortedByName(m.Entities) {
		if e.IdentifiedBy().Exists() {
			continue
		}
		name, _ := e.Name().Text()
		report.Report(diag.Diagnostic{
			Message:  fmt.Sprintf("entity %q has no identifiedBy", name),
			Location: e.Name().Location(),
		})
	}
}
