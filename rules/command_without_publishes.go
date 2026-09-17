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
		Description: "Every Command must publish at least one Event. A Command that publishes nothing expresses an intent without any consequence the model can see. The schema already requires this; the rule keeps the requirement in place independently of the schema.",
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
