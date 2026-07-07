package rules

import (
	"context"
	"fmt"

	"github.com/thenativeweb/esdm/diag"
	"github.com/thenativeweb/esdm/model"
)

const ruleIDReadModelWithoutProjections = "esdm/modeling/read-model-without-projections"

type readModelWithoutProjectionsRule struct{}

func newReadModelWithoutProjectionsRule() *readModelWithoutProjectionsRule {
	return &readModelWithoutProjectionsRule{}
}

func (*readModelWithoutProjectionsRule) Meta() Meta {
	return Meta{
		ID:          ruleIDReadModelWithoutProjections,
		Severity:    diag.SeverityWarning,
		Description: "A read-model should have at least one projection; mirrors the JSON Schema's minItems: 1 constraint as defense in depth so an accidental schema relaxation still surfaces the modeling issue.",
	}
}

func (*readModelWithoutProjectionsRule) Check(ctx context.Context, m *model.Model, report diag.Reporter) {
	for _, rm := range sortedByName(m.ReadModels) {
		if len(rm.Projections().Seq()) > 0 {
			continue
		}
		name, _ := rm.Name().Text()
		report.Report(diag.Diagnostic{
			Message:  fmt.Sprintf("read-model %q has no projections", name),
			Location: rm.Name().Location(),
		})
	}
}
