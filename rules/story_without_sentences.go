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
		Extension:   "domain-storytelling",
		Severity:    diag.SeverityWarning,
		Description: "Every Domain Story must have at least one Sentence; a story without Sentences tells nothing. The schema already requires this; the rule keeps the requirement in place independently of the schema.",
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
