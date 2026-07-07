package rules

import (
	"context"
	"fmt"

	"github.com/thenativeweb/esdm/diag"
	"github.com/thenativeweb/esdm/model"
)

const ruleIDStoryOrphanWorkObject = "esdm/modeling/story-orphan-work-object"

type storyOrphanWorkObjectRule struct{}

func newStoryOrphanWorkObjectRule() *storyOrphanWorkObjectRule {
	return &storyOrphanWorkObjectRule{}
}

func (*storyOrphanWorkObjectRule) Meta() Meta {
	return Meta{
		ID:          ruleIDStoryOrphanWorkObject,
		Severity:    diag.SeverityWarning,
		Description: "A work-object declared inside a sentence's workObjects[] list should be referenced from at least one edge in that same sentence; a declaration without any drawn presence is decorative. Implicit work objects - referenced from edges only, never listed in workObjects[] - are intentional and not flagged.",
	}
}

func (*storyOrphanWorkObjectRule) Check(ctx context.Context, m *model.Model, report diag.Reporter) {
	for _, story := range sortedByName(m.Extensions.DomainStorytelling.Stories) {
		storyName, _ := story.Name().Text()
		for _, sentence := range story.Sentences().Seq() {
			referenced := make(map[string]bool)
			for _, edge := range sentence.Field("edges").Seq() {
				for _, slot := range []string{"from", "to"} {
					if name, ok := edge.Field(slot).Field("workObject").Text(); ok {
						referenced[name] = true
					}
				}
			}

			seq, _ := sentence.Field("sequenceNumber").Int()
			for _, wo := range sentence.Field("workObjects").Seq() {
				name, ok := wo.Field("name").Text()
				if !ok {
					continue
				}
				if referenced[name] {
					continue
				}
				report.Report(diag.Diagnostic{
					Message:  fmt.Sprintf("domain-story %q sentence %d declares work-object %q but no edge in that sentence references it", storyName, seq, name),
					Location: wo.Field("name").Location(),
				})
			}
		}
	}
}
