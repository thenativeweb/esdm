package rules

import (
	"context"
	"fmt"

	"github.com/thenativeweb/esdm/diag"
	"github.com/thenativeweb/esdm/model"
)

const ruleIDStoryOrphanActor = "esdm/modeling/story-orphan-actor"

type storyOrphanActorRule struct{}

func newStoryOrphanActorRule() *storyOrphanActorRule {
	return &storyOrphanActorRule{}
}

func (*storyOrphanActorRule) Meta() Meta {
	return Meta{
		ID:          ruleIDStoryOrphanActor,
		Severity:    diag.SeverityWarning,
		Description: "An actor declared in a domain-story's actors[] list should be referenced from at least one edge in the story; an actor that nothing draws is decorative. Implicit actors - referenced from edges only, never listed in actors[] - are intentional and not flagged.",
	}
}

func (*storyOrphanActorRule) Check(ctx context.Context, m *model.Model, report diag.Reporter) {
	for _, story := range sortedByName(m.Extensions.DomainStorytelling.Stories) {
		referenced := referencedActorsInStory(story)

		storyName, _ := story.Name().Text()
		for _, actor := range story.Actors().Seq() {
			name, ok := actor.Field("name").Text()
			if !ok {
				continue
			}
			if referenced[name] {
				continue
			}
			report.Report(diag.Diagnostic{
				Message:  fmt.Sprintf("domain-story %q declares actor %q but no edge references it", storyName, name),
				Location: actor.Field("name").Location(),
			})
		}
	}
}

// referencedActorsInStory walks every edge in every
// sentence of a domain-story and returns the set of actor
// names referenced from the `from` or `to` slot.
func referencedActorsInStory(story model.DomainStoryView) map[string]bool {
	out := make(map[string]bool)
	for _, sentence := range story.Sentences().Seq() {
		for _, edge := range sentence.Field("edges").Seq() {
			for _, slot := range []string{"from", "to"} {
				if name, ok := edge.Field(slot).Field("actor").Text(); ok {
					out[name] = true
				}
			}
		}
	}
	return out
}
