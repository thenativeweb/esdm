package rules

import (
	"context"
	"fmt"

	"github.com/thenativeweb/esdm/diag"
	"github.com/thenativeweb/esdm/model"
)

const ruleIDProcessManagerTimerAtField = "esdm/structure/process-manager-timer-at-field"

type processManagerTimerAtFieldRule struct{}

func newProcessManagerTimerAtFieldRule() *processManagerTimerAtFieldRule {
	return &processManagerTimerAtFieldRule{}
}

func (*processManagerTimerAtFieldRule) Meta() Meta {
	return Meta{
		ID:          ruleIDProcessManagerTimerAtField,
		Severity:    diag.SeverityError,
		Description: "When a process-manager timer uses the absolute `at` shape, the named field must be declared in the process-manager's state.properties.",
	}
}

func (*processManagerTimerAtFieldRule) Check(ctx context.Context, m *model.Model, report diag.Reporter) {
	for _, pm := range sortedByName(m.ProcessManagers) {
		pmName, _ := pm.Name().Text()

		state := pm.State()
		for _, timer := range pm.Timers().Seq() {
			atNode := timer.Field("at")
			if !atNode.Exists() {
				continue
			}
			fieldName, ok := atNode.Text()
			if !ok {
				continue
			}

			timerName, _ := timer.Field("name").Text()
			if SchemaHasProperty(state, fieldName) {
				continue
			}
			report.Report(diag.Diagnostic{
				Message: fmt.Sprintf(
					"process-manager %q timer %q references state field %q, which is not declared in state.properties",
					pmName, timerName, fieldName,
				),
				Location: atNode.Location(),
			})
		}
	}
}
