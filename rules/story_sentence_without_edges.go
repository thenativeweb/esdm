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
		Extension:   "domain-storytelling",
		Severity:    diag.SeverityWarning,
		Description: "Every Sentence of a Domain Story must draw at least one edge; a Sentence without edges tells nothing. The schema already requires this; the rule keeps the requirement in place independently of the schema.",
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
