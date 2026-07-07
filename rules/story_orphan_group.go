package rules

import (
	"context"
	"fmt"

	"github.com/thenativeweb/esdm/diag"
	"github.com/thenativeweb/esdm/model"
)

const ruleIDStoryOrphanGroup = "esdm/modeling/story-orphan-group"

type storyOrphanGroupRule struct{}

func newStoryOrphanGroupRule() *storyOrphanGroupRule {
	return &storyOrphanGroupRule{}
}

func (*storyOrphanGroupRule) Meta() Meta {
	return Meta{
		ID:          ruleIDStoryOrphanGroup,
		Severity:    diag.SeverityWarning,
		Description: "Every group declared in a domain-story's top-level groups[] registry should be referenced as membership by at least one actor, work-object, or edge; an empty group is a registry entry without anyone in it.",
	}
}

func (*storyOrphanGroupRule) Check(ctx context.Context, m *model.Model, report diag.Reporter) {
	for _, story := range sortedByName(m.Extensions.DomainStorytelling.Stories) {
		members := referencedGroupsInStory(story)

		storyName, _ := story.Name().Text()
		for _, group := range story.Groups().Seq() {
			name, ok := group.Field("name").Text()
			if !ok {
				continue
			}
			if members[name] {
				continue
			}
			report.Report(diag.Diagnostic{
				Message:  fmt.Sprintf("domain-story %q declares group %q but no actor, work-object, or edge claims it as membership", storyName, name),
				Location: group.Field("name").Location(),
			})
		}
	}
}

// referencedGroupsInStory collects every group name
// claimed via an inline `groups: [name, ...]` field on
// any actor, work-object, or edge across the story.
func referencedGroupsInStory(story model.DomainStoryView) map[string]bool {
	out := make(map[string]bool)

	for _, actor := range story.Actors().Seq() {
		for _, item := range actor.Field("groups").Seq() {
			if name, ok := item.Text(); ok {
				out[name] = true
			}
		}
	}
	for _, sentence := range story.Sentences().Seq() {
		for _, wo := range sentence.Field("workObjects").Seq() {
			for _, item := range wo.Field("groups").Seq() {
				if name, ok := item.Text(); ok {
					out[name] = true
				}
			}
		}
		for _, edge := range sentence.Field("edges").Seq() {
			for _, item := range edge.Field("groups").Seq() {
				if name, ok := item.Text(); ok {
					out[name] = true
				}
			}
		}
	}
	return out
}
