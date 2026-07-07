package rules

import (
	"context"
	"fmt"

	"github.com/thenativeweb/esdm/ast"
	"github.com/thenativeweb/esdm/diag"
	"github.com/thenativeweb/esdm/model"
)

const ruleIDStoryUnknownGroupMembership = "esdm/structure/story-unknown-group-membership"

type storyUnknownGroupMembershipRule struct{}

func newStoryUnknownGroupMembershipRule() *storyUnknownGroupMembershipRule {
	return &storyUnknownGroupMembershipRule{}
}

func (*storyUnknownGroupMembershipRule) Meta() Meta {
	return Meta{
		ID:          ruleIDStoryUnknownGroupMembership,
		Severity:    diag.SeverityError,
		Description: "Every group name claimed via an inline `groups: [name, ...]` membership on an actor, work-object or edge must exist in the story's top-level groups[] registry. Membership in an undeclared group is a stale reference.",
	}
}

func (*storyUnknownGroupMembershipRule) Check(ctx context.Context, m *model.Model, report diag.Reporter) {
	for _, story := range sortedByName(m.Extensions.DomainStorytelling.Stories) {
		declared := declaredGroupsInStory(story)
		storyName, _ := story.Name().Text()

		for _, actor := range story.Actors().Seq() {
			actorName, _ := actor.Field("name").Text()
			for _, item := range actor.Field("groups").Seq() {
				groupName, ok := item.Text()
				if !ok || declared[groupName] {
					continue
				}
				report.Report(diag.Diagnostic{
					Message:  fmt.Sprintf("domain-story %q actor %q claims membership in undeclared group %q", storyName, actorName, groupName),
					Location: item.Location(),
				})
			}
		}

		for _, sentence := range story.Sentences().Seq() {
			for _, wo := range sentence.Field("workObjects").Seq() {
				woName, _ := wo.Field("name").Text()
				for _, item := range wo.Field("groups").Seq() {
					groupName, ok := item.Text()
					if !ok || declared[groupName] {
						continue
					}
					report.Report(diag.Diagnostic{
						Message:  fmt.Sprintf("domain-story %q work-object %q claims membership in undeclared group %q", storyName, woName, groupName),
						Location: item.Location(),
					})
				}
			}
			for _, edge := range sentence.Field("edges").Seq() {
				for _, item := range edge.Field("groups").Seq() {
					groupName, ok := item.Text()
					if !ok || declared[groupName] {
						continue
					}
					report.Report(diag.Diagnostic{
						Message:  fmt.Sprintf("domain-story %q edge %s claims membership in undeclared group %q", storyName, edgeLabel(edge), groupName),
						Location: item.Location(),
					})
				}
			}
		}
	}
}

// declaredGroupsInStory returns the set of group names
// in the story's top-level groups[] registry.
func declaredGroupsInStory(story model.DomainStoryView) map[string]bool {
	out := make(map[string]bool)
	for _, group := range story.Groups().Seq() {
		if name, ok := group.Field("name").Text(); ok {
			out[name] = true
		}
	}
	return out
}

// edgeLabel renders a short, human-friendly identifier
// for an edge in diagnostic messages - its label if
// present, otherwise a "from -> to" sketch.
func edgeLabel(edge ast.Node) string {
	if label, ok := edge.Field("label").Text(); ok && label != "" {
		return fmt.Sprintf("%q", label)
	}
	from := nodeRefLabel(edge.Field("from"))
	to := nodeRefLabel(edge.Field("to"))
	return fmt.Sprintf("%s → %s", from, to)
}

func nodeRefLabel(ref ast.Node) string {
	if name, ok := ref.Field("actor").Text(); ok {
		return "actor:" + name
	}
	if name, ok := ref.Field("workObject").Text(); ok {
		return "workObject:" + name
	}
	return "?"
}
