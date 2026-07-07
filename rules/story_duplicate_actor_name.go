package rules

import (
	"context"
	"fmt"

	"github.com/thenativeweb/esdm/diag"
	"github.com/thenativeweb/esdm/model"
)

const ruleIDStoryDuplicateActorName = "esdm/structure/story-duplicate-actor-name"

type storyDuplicateActorNameRule struct{}

func newStoryDuplicateActorNameRule() *storyDuplicateActorNameRule {
	return &storyDuplicateActorNameRule{}
}

func (*storyDuplicateActorNameRule) Meta() Meta {
	return Meta{
		ID:          ruleIDStoryDuplicateActorName,
		Severity:    diag.SeverityError,
		Description: "Within a domain-story's actors[] list every name appears at most once. Domain Storytelling is built around the rule \"one icon per actor in a story\"; a redeclaration with the same name confuses the diagram and any downstream reference resolution.",
	}
}

func (*storyDuplicateActorNameRule) Check(ctx context.Context, m *model.Model, report diag.Reporter) {
	for _, story := range sortedByName(m.Extensions.DomainStorytelling.Stories) {
		storyName, _ := story.Name().Text()
		seen := make(map[string]bool)
		for _, actor := range story.Actors().Seq() {
			name, ok := actor.Field("name").Text()
			if !ok {
				continue
			}
			if !seen[name] {
				seen[name] = true
				continue
			}
			report.Report(diag.Diagnostic{
				Message:  fmt.Sprintf("domain-story %q declares actor %q more than once in actors[]", storyName, name),
				Location: actor.Field("name").Location(),
			})
		}
	}
}
