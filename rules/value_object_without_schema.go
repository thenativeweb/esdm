package rules

import (
	"context"
	"fmt"

	"github.com/thenativeweb/esdm/diag"
	"github.com/thenativeweb/esdm/model"
)

const ruleIDValueObjectWithoutSchema = "esdm/modeling/value-object-without-schema"

type valueObjectWithoutSchemaRule struct{}

func newValueObjectWithoutSchemaRule() *valueObjectWithoutSchemaRule {
	return &valueObjectWithoutSchemaRule{}
}

func (*valueObjectWithoutSchemaRule) Meta() Meta {
	return Meta{
		ID:          ruleIDValueObjectWithoutSchema,
		Severity:    diag.SeverityWarning,
		Description: "Every Value Object must declare a `schema` field describing the shape of one instance. The schema already requires the field; the rule keeps the requirement in place independently of the schema.",
	}
}

func (*valueObjectWithoutSchemaRule) Check(ctx context.Context, m *model.Model, report diag.Reporter) {
	for _, vo := range sortedByName(m.ValueObjects) {
		if vo.Schema().Exists() {
			continue
		}
		name, _ := vo.Name().Text()
		report.Report(diag.Diagnostic{
			Message:  fmt.Sprintf("value-object %q has no schema", name),
			Location: vo.Name().Location(),
		})
	}
}
