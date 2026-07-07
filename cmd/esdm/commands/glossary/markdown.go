package glossary

import (
	"fmt"
	"strings"
)

// Render turns a resolved glossary into Markdown. The output
// always starts with a single `# Glossary` heading - an
// empty glossary renders as just that heading plus a
// trailing newline - followed by one `##` section per
// bounded context and one `###` entry per term. Every
// heading is separated from the following paragraph by a
// blank line. Each discouraged alternative becomes its own
// paragraph: an italicized "Avoid the term ..." sentence
// followed, when present, by the reason as a plain sentence.
// The output always ends with exactly one trailing newline.
func Render(g *Glossary) string {
	var b strings.Builder
	b.WriteString("# Glossary\n")

	for _, section := range g.Sections {
		fmt.Fprintf(&b, "\n## %s\n", section.BoundedContext)
		for _, term := range section.Terms {
			fmt.Fprintf(&b, "\n### %s\n\n", term.Term)
			b.WriteString(term.Definition)
			b.WriteString("\n")

			for _, avoid := range term.Avoid {
				fmt.Fprintf(&b, "\n_Avoid the term %q._", avoid.Term)
				if avoid.Reason != "" {
					b.WriteString(" ")
					b.WriteString(avoid.Reason)
				}
				b.WriteString("\n")
			}
		}
	}

	return b.String()
}
