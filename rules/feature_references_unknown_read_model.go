package rules

import (
	"context"
	"fmt"

	"github.com/thenativeweb/esdm/diag"
	"github.com/thenativeweb/esdm/model"
)

const ruleIDFeatureReferencesUnknownReadModel = "esdm/gwt/feature-references-unknown-read-model"

type featureReferencesUnknownReadModelRule struct{}

func newFeatureReferencesUnknownReadModelRule() *featureReferencesUnknownReadModelRule {
	return &featureReferencesUnknownReadModelRule{}
}

func (*featureReferencesUnknownReadModelRule) Meta() Meta {
	return Meta{
		ID:          ruleIDFeatureReferencesUnknownReadModel,
		Severity:    diag.SeverityError,
		Description: "A read-model-scoped feature must point at a declared read-model; an unresolved scope means every scenario underneath is testing nothing.",
	}
}

func (*featureReferencesUnknownReadModelRule) Check(ctx context.Context, m *model.Model, report diag.Reporter) {
	for _, feature := range sortedByName(m.Extensions.GivenWhenThen.Features) {
		if featureVariant(feature) != featureVariantReadModel {
			continue
		}
		scope := feature.Scope()
		domain := scopeField(scope, "domain")
		boundedContext := scopeField(scope, "boundedContext")
		readModel := scopeField(scope, "readModel")
		if _, ok := m.LookupReadModel(domain, boundedContext, readModel); ok {
			continue
		}
		featureName, _ := feature.Name().Text()
		report.Report(diag.Diagnostic{
			Message:  fmt.Sprintf("feature %q scope references read-model %q in bounded-context %q which is not declared in the model", featureName, readModel, boundedContext),
			Location: scope.Field("readModel").Location(),
		})
	}
}
