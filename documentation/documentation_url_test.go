package documentation_test

import (
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thenativeweb/esdm/rules"
)

// TestDocumentationURLsPointAtExistingPages makes sure the URL a
// diagnostic carries resolves to a page of this documentation. The
// rules package derives the page path from a rule's extension
// without knowing the site layout, and the pages are maintained by
// hand, so this test is what ties the two together: a renamed or
// moved page fails here before a user follows a dead link.
func TestDocumentationURLsPointAtExistingPages(t *testing.T) {
	t.Parallel()

	byPage := map[string]string{}
	for _, rule := range rules.Catalog() {
		meta := rule.Meta()
		byPage[pagePathOf(t, meta.DocumentationURL())] = meta.ID
	}

	// The pipeline diagnostics carry a bare ID and resolve to the
	// core page; include one so the core page is covered even if
	// every catalog rule ever moved to an extension.
	pipelineMeta := rules.Meta{ID: "esdm/structure/unresolved-reference"}
	byPage[pagePathOf(t, pipelineMeta.DocumentationURL())] = pipelineMeta.ID

	for page, ruleID := range byPage {
		t.Run(page, func(t *testing.T) {
			t.Parallel()

			_, err := os.Stat(filepath.Join("docs", filepath.FromSlash(page)))
			assert.NoErrorf(t, err, "rule %s links to a page that does not exist", ruleID)
		})
	}
}

// pagePathOf turns a documentation URL into the Markdown source
// path relative to docs/: the site serves directory URLs, so
// /reference/linter-rules/ is built from reference/linter-rules.md.
func pagePathOf(t *testing.T, documentationURL string) string {
	t.Helper()

	parsed, err := url.Parse(documentationURL)
	require.NoError(t, err)
	require.Equal(t, "www.esdm.io", parsed.Host)
	require.NotEmpty(t, parsed.Fragment, "documentation URL %s has no anchor", documentationURL)

	return strings.TrimSuffix(strings.TrimPrefix(parsed.Path, "/"), "/") + ".md"
}
