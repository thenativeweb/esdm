package rules

import (
	"context"
	"fmt"

	"github.com/thenativeweb/esdm/diag"
	"github.com/thenativeweb/esdm/model"
)

const ruleIDQueryWithoutReadModel = "esdm/modeling/query-without-read-model"

type queryWithoutReadModelRule struct{}

func newQueryWithoutReadModelRule() *queryWithoutReadModelRule {
	return &queryWithoutReadModelRule{}
}

func (*queryWithoutReadModelRule) Meta() Meta {
	return Meta{
		ID:          ruleIDQueryWithoutReadModel,
		Severity:    diag.SeverityWarning,
		Description: "Every Query must name the Read Model it reads from. The schema already requires the field; the rule keeps the requirement in place independently of the schema.",
	}
}

func (*queryWithoutReadModelRule) Check(ctx context.Context, m *model.Model, report diag.Reporter) {
	for _, q := range sortedByName(m.Queries) {
		if q.ReadModel().Exists() {
			continue
		}
		name, _ := q.Name().Text()
		report.Report(diag.Diagnostic{
			Message:  fmt.Sprintf("query %q has no readModel reference", name),
			Location: q.Name().Location(),
		})
	}
}
