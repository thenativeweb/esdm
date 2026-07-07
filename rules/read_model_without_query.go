package rules

import (
	"context"
	"fmt"

	"github.com/thenativeweb/esdm/diag"
	"github.com/thenativeweb/esdm/model"
)

const ruleIDReadModelWithoutQuery = "esdm/modeling/read-model-without-query"

type readModelWithoutQueryRule struct{}

func newReadModelWithoutQueryRule() *readModelWithoutQueryRule {
	return &readModelWithoutQueryRule{}
}

func (*readModelWithoutQueryRule) Meta() Meta {
	return Meta{
		ID:          ruleIDReadModelWithoutQuery,
		Severity:    diag.SeverityWarning,
		Description: "A read-model should have at least one query reading from it; otherwise its purpose is unclear.",
	}
}

func (*readModelWithoutQueryRule) Check(ctx context.Context, m *model.Model, report diag.Reporter) {
	// Queries and the read-models they read from are
	// co-located in the same bounded context. The used
	// set keys on (domain, BC, bare name) so equally-named
	// read-models in different BCs are checked against
	// their own queries only.
	used := make(map[string]bool)
	for _, q := range m.Queries {
		domain := scopeField(q.Scope(), "domain")
		bc := scopeField(q.Scope(), "boundedContext")
		if name, ok := q.ReadModel().Text(); ok {
			used[domain+"/"+bc+"/"+name] = true
		}
	}

	for _, rm := range sortedByName(m.ReadModels) {
		domain := scopeField(rm.Scope(), "domain")
		bc := scopeField(rm.Scope(), "boundedContext")
		name, _ := rm.Name().Text()
		if used[domain+"/"+bc+"/"+name] {
			continue
		}
		report.Report(diag.Diagnostic{
			Message:  fmt.Sprintf("read-model %q is not accessed by any query", name),
			Location: rm.Name().Location(),
		})
	}
}
