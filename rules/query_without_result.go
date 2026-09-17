package rules

import (
	"context"
	"fmt"

	"github.com/thenativeweb/esdm/diag"
	"github.com/thenativeweb/esdm/model"
)

const ruleIDQueryWithoutResult = "esdm/modeling/query-without-result"

type queryWithoutResultRule struct{}

func newQueryWithoutResultRule() *queryWithoutResultRule {
	return &queryWithoutResultRule{}
}

func (*queryWithoutResultRule) Meta() Meta {
	return Meta{
		ID:          ruleIDQueryWithoutResult,
		Severity:    diag.SeverityWarning,
		Description: "Every Query must declare a `result` schema. The schema already requires the field; the rule keeps the requirement in place independently of the schema.",
	}
}

func (*queryWithoutResultRule) Check(ctx context.Context, m *model.Model, report diag.Reporter) {
	for _, q := range sortedByName(m.Queries) {
		if q.Result().Exists() {
			continue
		}
		name, _ := q.Name().Text()
		report.Report(diag.Diagnostic{
			Message:  fmt.Sprintf("query %q has no result schema", name),
			Location: q.Name().Location(),
		})
	}
}
