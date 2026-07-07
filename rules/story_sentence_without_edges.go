package rules

import (
	"context"
	"fmt"

	"github.com/thenativeweb/esdm/diag"
	"github.com/thenativeweb/esdm/model"
)

const ruleIDStorySentenceWithoutEdges = "esdm/structure/story-sentence-without-edges"

type storySentenceWithoutEdgesRule struct{}

func newStorySentenceWithoutEdgesRule() *storySentenceWithoutEdgesRule {
	return &storySentenceWithoutEdgesRule{}
}

func (*storySentenceWithoutEdgesRule) Meta() Meta {
	return Meta{
		ID:          ruleIDStorySentenceWithoutEdges,
		Severity:    diag.SeverityWarning,
		Description: "Every sentence in a domain-story must carry at least one edge; mirrors the storytelling schema's required + minItems: 1 constraint on edges as defense in depth so an accidental schema relaxation still surfaces the modeling issue.",
	}
}

func (*storySentenceWithoutEdgesRule) Check(ctx context.Context, m *model.Model, report diag.Reporter) {
	for _, story := range sortedByName(m.Extensions.DomainStorytelling.Stories) {
		name, _ := story.Name().Text()
		for _, sentence := range story.Sentences().Seq() {
			if len(sentence.Field("edges").Seq()) > 0 {
				continue
			}
			seq, _ := sentence.Field("sequenceNumber").Int()
			report.Report(diag.Diagnostic{
				Message:  fmt.Sprintf("domain-story %q sentence %d has no edges", name, seq),
				Location: sentence.Location(),
			})
		}
	}
}
