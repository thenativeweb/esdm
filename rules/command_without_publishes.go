package rules

import (
	"context"
	"fmt"

	"github.com/thenativeweb/esdm/diag"
	"github.com/thenativeweb/esdm/model"
)

const ruleIDCommandWithoutPublishes = "esdm/modeling/command-without-publishes"

type commandWithoutPublishesRule struct{}

func newCommandWithoutPublishesRule() *commandWithoutPublishesRule {
	return &commandWithoutPublishesRule{}
}

func (*commandWithoutPublishesRule) Meta() Meta {
	return Meta{
		ID:          ruleIDCommandWithoutPublishes,
		Severity:    diag.SeverityWarning,
		Description: "Every command must publish at least one event; mirrors the JSON Schema's required + minItems: 1 constraint as defense in depth. A command with no publishes is a no-op intent.",
	}
}

func (*commandWithoutPublishesRule) Check(ctx context.Context, m *model.Model, report diag.Reporter) {
	for _, cmd := range sortedByName(m.Commands) {
		if len(cmd.Publishes().Seq()) > 0 {
			continue
		}
		name, _ := cmd.Name().Text()
		report.Report(diag.Diagnostic{
			Message:  fmt.Sprintf("command %q publishes no events", name),
			Location: cmd.Name().Location(),
		})
	}
}
