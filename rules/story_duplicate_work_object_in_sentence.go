package rules

import (
	"context"
	"fmt"

	"github.com/thenativeweb/esdm/diag"
	"github.com/thenativeweb/esdm/model"
)

const ruleIDStoryDuplicateWorkObjectInSentence = "esdm/structure/story-duplicate-work-object-in-sentence"

type storyDuplicateWorkObjectInSentenceRule struct{}

func newStoryDuplicateWorkObjectInSentenceRule() *storyDuplicateWorkObjectInSentenceRule {
	return &storyDuplicateWorkObjectInSentenceRule{}
}

func (*storyDuplicateWorkObjectInSentenceRule) Meta() Meta {
	return Meta{
		ID:          ruleIDStoryDuplicateWorkObjectInSentence,
		Severity:    diag.SeverityError,
		Description: "Within a single sentence's workObjects[] list every name appears at most once. Domain Storytelling redraws work objects per sentence, so the same name in a *different* sentence is a fresh instance and not flagged; a redeclaration *inside the same sentence* is a modeling mistake.",
	}
}

func (*storyDuplicateWorkObjectInSentenceRule) Check(ctx context.Context, m *model.Model, report diag.Reporter) {
	for _, story := range sortedByName(m.Extensions.DomainStorytelling.Stories) {
		storyName, _ := story.Name().Text()
		for _, sentence := range story.Sentences().Seq() {
			seq, _ := sentence.Field("sequenceNumber").Int()
			seen := make(map[string]bool)
			for _, wo := range sentence.Field("workObjects").Seq() {
				name, ok := wo.Field("name").Text()
				if !ok {
					continue
				}
				if !seen[name] {
					seen[name] = true
					continue
				}
				report.Report(diag.Diagnostic{
					Message:  fmt.Sprintf("domain-story %q sentence %d declares work-object %q more than once in workObjects[]", storyName, seq, name),
					Location: wo.Field("name").Location(),
				})
			}
		}
	}
}
