package rules

import (
	"context"
	"fmt"

	"github.com/thenativeweb/esdm/diag"
	"github.com/thenativeweb/esdm/model"
)

const ruleIDFeatureWithoutScenarios = "esdm/gwt/feature-without-scenarios"

type featureWithoutScenariosRule struct{}

func newFeatureWithoutScenariosRule() *featureWithoutScenariosRule {
	return &featureWithoutScenariosRule{}
}

func (*featureWithoutScenariosRule) Meta() Meta {
	return Meta{
		ID:          ruleIDFeatureWithoutScenarios,
		Severity:    diag.SeverityWarning,
		Description: "Every feature must declare at least one scenario; mirrors the given-when-then schema's required + minItems: 1 constraint as defense in depth so an accidental schema relaxation still surfaces the modeling issue.",
	}
}

func (*featureWithoutScenariosRule) Check(ctx context.Context, m *model.Model, report diag.Reporter) {
	for _, feature := range sortedByName(m.Extensions.GivenWhenThen.Features) {
		if len(feature.Scenarios().Seq()) > 0 {
			continue
		}
		name, _ := feature.Name().Text()
		report.Report(diag.Diagnostic{
			Message:  fmt.Sprintf("feature %q has no scenarios", name),
			Location: feature.Name().Location(),
		})
	}
}
