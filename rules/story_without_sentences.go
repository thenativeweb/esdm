package rules

import (
	"context"
	"fmt"

	"github.com/thenativeweb/esdm/diag"
	"github.com/thenativeweb/esdm/model"
)

const ruleIDStoryWithoutSentences = "esdm/structure/story-without-sentences"

type storyWithoutSentencesRule struct{}

func newStoryWithoutSentencesRule() *storyWithoutSentencesRule {
	return &storyWithoutSentencesRule{}
}

func (*storyWithoutSentencesRule) Meta() Meta {
	return Meta{
		ID:          ruleIDStoryWithoutSentences,
		Severity:    diag.SeverityWarning,
		Description: "Every domain-story must declare at least one sentence; mirrors the storytelling schema's required + minItems: 1 constraint as defense in depth so an accidental schema relaxation still surfaces the modeling issue.",
	}
}

func (*storyWithoutSentencesRule) Check(ctx context.Context, m *model.Model, report diag.Reporter) {
	for _, story := range sortedByName(m.Extensions.DomainStorytelling.Stories) {
		if len(story.Sentences().Seq()) > 0 {
			continue
		}
		name, _ := story.Name().Text()
		report.Report(diag.Diagnostic{
			Message:  fmt.Sprintf("domain-story %q has no sentences", name),
			Location: story.Name().Location(),
		})
	}
}
