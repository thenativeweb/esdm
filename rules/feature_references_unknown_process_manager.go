package rules

import (
	"context"
	"fmt"

	"github.com/thenativeweb/esdm/diag"
	"github.com/thenativeweb/esdm/model"
)

const ruleIDFeatureReferencesUnknownProcessManager = "esdm/gwt/feature-references-unknown-process-manager"

type featureReferencesUnknownProcessManagerRule struct{}

func newFeatureReferencesUnknownProcessManagerRule() *featureReferencesUnknownProcessManagerRule {
	return &featureReferencesUnknownProcessManagerRule{}
}

func (*featureReferencesUnknownProcessManagerRule) Meta() Meta {
	return Meta{
		ID:          ruleIDFeatureReferencesUnknownProcessManager,
		Severity:    diag.SeverityError,
		Description: "A process-manager-scoped feature must point at a declared process-manager; an unresolved scope means every scenario underneath is testing nothing.",
	}
}

func (*featureReferencesUnknownProcessManagerRule) Check(ctx context.Context, m *model.Model, report diag.Reporter) {
	for _, feature := range sortedByName(m.Extensions.GivenWhenThen.Features) {
		if featureVariant(feature) != featureVariantProcessManager {
			continue
		}
		scope := feature.Scope()
		domain := scopeField(scope, "domain")
		processManager := scopeField(scope, "processManager")
		if _, ok := m.LookupProcessManager(domain, processManager); ok {
			continue
		}
		featureName, _ := feature.Name().Text()
		report.Report(diag.Diagnostic{
			Message:  fmt.Sprintf("feature %q scope references process-manager %q in domain %q which is not declared in the model", featureName, processManager, domain),
			Location: scope.Field("processManager").Location(),
		})
	}
}
