package rules

import (
	"context"
	"fmt"

	"github.com/thenativeweb/esdm/diag"
	"github.com/thenativeweb/esdm/model"
)

const ruleIDCommandWithoutData = "esdm/modeling/command-without-data"

type commandWithoutDataRule struct{}

func newCommandWithoutDataRule() *commandWithoutDataRule {
	return &commandWithoutDataRule{}
}

func (*commandWithoutDataRule) Meta() Meta {
	return Meta{
		ID:          ruleIDCommandWithoutData,
		Severity:    diag.SeverityWarning,
		Description: "Every command must declare a data schema (an empty schema is allowed); mirrors the JSON Schema's required: [data] constraint as defense in depth so 'no payload' is expressed deliberately.",
	}
}

func (*commandWithoutDataRule) Check(ctx context.Context, m *model.Model, report diag.Reporter) {
	for _, cmd := range sortedByName(m.Commands) {
		if cmd.Data().Exists() {
			continue
		}
		name, _ := cmd.Name().Text()
		report.Report(diag.Diagnostic{
			Message:  fmt.Sprintf("command %q has no data schema", name),
			Location: cmd.Name().Location(),
		})
	}
}
