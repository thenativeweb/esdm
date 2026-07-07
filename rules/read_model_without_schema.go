package rules

import (
	"context"
	"fmt"

	"github.com/thenativeweb/esdm/diag"
	"github.com/thenativeweb/esdm/model"
)

const ruleIDReadModelWithoutSchema = "esdm/modeling/read-model-without-schema"

type readModelWithoutSchemaRule struct{}

func newReadModelWithoutSchemaRule() *readModelWithoutSchemaRule {
	return &readModelWithoutSchemaRule{}
}

func (*readModelWithoutSchemaRule) Meta() Meta {
	return Meta{
		ID:          ruleIDReadModelWithoutSchema,
		Severity:    diag.SeverityWarning,
		Description: "Every read-model must declare a schema describing what is materialized; mirrors the JSON Schema's required: [schema] constraint as defense in depth.",
	}
}

func (*readModelWithoutSchemaRule) Check(ctx context.Context, m *model.Model, report diag.Reporter) {
	for _, rm := range sortedByName(m.ReadModels) {
		if rm.Schema().Exists() {
			continue
		}
		name, _ := rm.Name().Text()
		report.Report(diag.Diagnostic{
			Message:  fmt.Sprintf("read-model %q has no schema", name),
			Location: rm.Name().Location(),
		})
	}
}
