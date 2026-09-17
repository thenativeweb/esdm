package rules

import (
	"context"
	"fmt"

	"github.com/thenativeweb/esdm/diag"
	"github.com/thenativeweb/esdm/model"
)

const ruleIDEntityWithoutSchema = "esdm/modeling/entity-without-schema"

type entityWithoutSchemaRule struct{}

func newEntityWithoutSchemaRule() *entityWithoutSchemaRule {
	return &entityWithoutSchemaRule{}
}

func (*entityWithoutSchemaRule) Meta() Meta {
	return Meta{
		ID:          ruleIDEntityWithoutSchema,
		Severity:    diag.SeverityWarning,
		Description: "Every Entity must declare a `schema` field describing the shape of one instance. The schema already requires the field; the rule keeps the requirement in place independently of the schema.",
	}
}

func (*entityWithoutSchemaRule) Check(ctx context.Context, m *model.Model, report diag.Reporter) {
	for _, e := range sortedByName(m.Entities) {
		if e.Schema().Exists() {
			continue
		}
		name, _ := e.Name().Text()
		report.Report(diag.Diagnostic{
			Message:  fmt.Sprintf("entity %q has no schema", name),
			Location: e.Name().Location(),
		})
	}
}
