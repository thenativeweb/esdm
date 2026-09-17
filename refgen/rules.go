package refgen

import (
	"fmt"
	"path"
	"sort"
	"strings"

	"github.com/thenativeweb/esdm/rules"
)

// categoryOrder fixes the order of the category sections on a
// Linter Rules page. Structure comes first because its rules say
// what must hold before the model means anything; the modeling
// rules build on that. The remaining categories are listed for
// completeness so that a rule in one of them gets a section
// instead of being dropped.
var categoryOrder = []string{"structure", "naming", "modeling", "linguistic", "system", "gwt"}

// categoryLabels maps a category to its section heading.
var categoryLabels = map[string]string{
	"structure":  "Structure",
	"naming":     "Naming",
	"modeling":   "Modeling",
	"linguistic": "Linguistic",
	"system":     "System",
	"gwt":        "Given-When-Then",
}

// generateRules emits one Markdown snippet per rule source: the
// core rules and the rules of each extension that has any. A
// rule's Extension decides which snippet it lands on; its category
// decides the section within that snippet.
func generateRules(out map[string][]byte) error {
	byExtension := map[string][]rules.Meta{}
	for _, rule := range rules.Catalog() {
		meta := rule.Meta()
		byExtension[meta.Extension] = append(byExtension[meta.Extension], meta)
	}

	for extension, metas := range byExtension {
		text, err := renderRules(metas)
		if err != nil {
			return fmt.Errorf("%s: %w", rulesSnippetPath(extension), err)
		}
		out[rulesSnippetPath(extension)] = []byte(text)
	}
	return nil
}

// rulesSnippetPath places every rule snippet in one directory,
// reference/linter-rules, named after its source. The snippets do
// not follow the layout of the schema snippets on purpose: the core
// page lives at docs/reference/linter-rules.md, and both docs and
// snippets are snippet base paths, so a core snippet at the same
// relative path would make the page include itself.
func rulesSnippetPath(extension string) string {
	if extension == "" {
		return path.Join("reference", "linter-rules", "core.md")
	}
	return path.Join("reference", "linter-rules", extension+".md")
}

// renderRules writes the entries of one Linter Rules page. Every
// rule becomes a level-three heading carrying the full ID and an
// explicit anchor, followed by its severity and its description.
// The anchor is explicit because the heading text contains slashes
// and backticks, which the site's automatic slugs would strip into
// an unreadable fragment; the explicit form is what diagnostics
// can link to. The category headings are omitted when the page
// has a single category, because a lone heading would only repeat
// what the page title already says.
func renderRules(metas []rules.Meta) (string, error) {
	byCategory := map[string][]rules.Meta{}
	for _, meta := range metas {
		category := meta.Category()
		_, isKnown := categoryLabels[category]
		if !isKnown {
			return "", fmt.Errorf("rule %s has a category without a section label", meta.ID)
		}
		byCategory[category] = append(byCategory[category], meta)
	}

	var builder strings.Builder
	hasMultipleCategories := len(byCategory) > 1

	for _, category := range categoryOrder {
		categoryMetas := byCategory[category]
		if len(categoryMetas) == 0 {
			continue
		}
		sort.Slice(categoryMetas, func(i, j int) bool {
			return categoryMetas[i].ID < categoryMetas[j].ID
		})

		if hasMultipleCategories {
			fmt.Fprintf(&builder, "## %s\n\n", categoryLabels[category])
		}
		for _, meta := range categoryMetas {
			fmt.Fprintf(&builder, "### `%s` { #%s }\n\nSeverity: `%s`\n\n%s\n\n", meta.ID, meta.Anchor(), meta.Severity, meta.Description)
		}
	}

	// Every entry ends with a blank line so entries stay apart;
	// the final one is trimmed so the file ends with exactly one
	// newline, as .editorconfig requires.
	return strings.TrimSuffix(builder.String(), "\n"), nil
}
