package rules

import (
	"context"
	"fmt"

	"github.com/thenativeweb/esdm/diag"
	"github.com/thenativeweb/esdm/model"
)

const ruleIDDuplicateTranslationLanguage = "esdm/structure/duplicate-translation-language"

type duplicateTranslationLanguageRule struct{}

func newDuplicateTranslationLanguageRule() *duplicateTranslationLanguageRule {
	return &duplicateTranslationLanguageRule{}
}

func (*duplicateTranslationLanguageRule) Meta() Meta {
	return Meta{
		ID:          ruleIDDuplicateTranslationLanguage,
		Severity:    diag.SeverityError,
		Description: "A term has exactly one translation per language. Two translations into the same language contradict the one-term principle the ubiquitous language rests on.",
	}
}

// Check walks every term of every bounded context and
// reports the second and any further translation into a
// language that an earlier translation of the same term
// already uses. Like every array in the schema, the
// translations list uses its natural identifier - here the
// language - as an implicit key the schema does not enforce.
func (*duplicateTranslationLanguageRule) Check(ctx context.Context, m *model.Model, report diag.Reporter) {
	for _, bc := range sortedByName(m.BoundedContexts) {
		for _, entry := range bc.UbiquitousLanguage().Seq() {
			term, _ := entry.Field("term").Text()
			seen := make(map[string]bool)
			for _, translation := range entry.Field("translations").Seq() {
				languageNode := translation.Field("language")
				language, ok := languageNode.Text()
				if !ok {
					continue
				}
				if !seen[language] {
					seen[language] = true
					continue
				}
				report.Report(diag.Diagnostic{
					Message:  fmt.Sprintf("term %q has more than one translation into %q", term, language),
					Location: languageNode.Location(),
				})
			}
		}
	}
}
