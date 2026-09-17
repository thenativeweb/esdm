package rules

import (
	"context"
	"fmt"

	"github.com/thenativeweb/esdm/diag"
	"github.com/thenativeweb/esdm/model"
)

const ruleIDTranslationIntoOwnLanguage = "esdm/structure/translation-into-own-language"

type translationIntoOwnLanguageRule struct{}

func newTranslationIntoOwnLanguageRule() *translationIntoOwnLanguageRule {
	return &translationIntoOwnLanguageRule{}
}

func (*translationIntoOwnLanguageRule) Meta() Meta {
	return Meta{
		ID:          ruleIDTranslationIntoOwnLanguage,
		Severity:    diag.SeverityError,
		Description: "A translation into the language the Bounded Context itself is written in is a second term in that language, which the ubiquitous language does not allow.",
	}
}

// Check compares every translation's language with the
// language of the bounded context that owns the term. The
// main entry already is the term in that language, so a
// translation into it would be an alias, not a translation.
func (*translationIntoOwnLanguageRule) Check(ctx context.Context, m *model.Model, report diag.Reporter) {
	for _, bc := range sortedByName(m.BoundedContexts) {
		bcName, _ := bc.Name().Text()
		bcLanguage, ok := bc.Language().Text()
		if !ok {
			continue
		}
		for _, entry := range bc.UbiquitousLanguage().Seq() {
			term, _ := entry.Field("term").Text()
			for _, translation := range entry.Field("translations").Seq() {
				languageNode := translation.Field("language")
				language, ok := languageNode.Text()
				if !ok || language != bcLanguage {
					continue
				}
				report.Report(diag.Diagnostic{
					Message:  fmt.Sprintf("term %q is translated into %q, which is the language of bounded context %q", term, language, bcName),
					Location: languageNode.Location(),
				})
			}
		}
	}
}
