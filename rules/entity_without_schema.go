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
		Description: "Every entity must declare a schema describing the shape of one instance; mirrors the JSON Schema's required: [schema] constraint as defense in depth.",
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
