package rules

import (
	"context"
	"fmt"

	"github.com/thenativeweb/esdm/diag"
	"github.com/thenativeweb/esdm/model"
)

const ruleIDStoryDuplicateGroupName = "esdm/structure/story-duplicate-group-name"

type storyDuplicateGroupNameRule struct{}

func newStoryDuplicateGroupNameRule() *storyDuplicateGroupNameRule {
	return &storyDuplicateGroupNameRule{}
}

func (*storyDuplicateGroupNameRule) Meta() Meta {
	return Meta{
		ID:          ruleIDStoryDuplicateGroupName,
		Severity:    diag.SeverityError,
		Description: "Within a domain-story's top-level groups[] registry every group name appears at most once. Inline `groups: [name, ...]` memberships disambiguate by name; duplicate registry entries make membership references ambiguous.",
	}
}

func (*storyDuplicateGroupNameRule) Check(ctx context.Context, m *model.Model, report diag.Reporter) {
	for _, story := range sortedByName(m.Extensions.DomainStorytelling.Stories) {
		storyName, _ := story.Name().Text()
		seen := make(map[string]bool)
		for _, group := range story.Groups().Seq() {
			name, ok := group.Field("name").Text()
			if !ok {
				continue
			}
			if !seen[name] {
				seen[name] = true
				continue
			}
			report.Report(diag.Diagnostic{
				Message:  fmt.Sprintf("domain-story %q declares group %q more than once in groups[]", storyName, name),
				Location: group.Field("name").Location(),
			})
		}
	}
}
