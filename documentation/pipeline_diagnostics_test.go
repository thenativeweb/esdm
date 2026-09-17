package documentation_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/thenativeweb/esdm/rules"
)

// pipelinePackages are the packages that emit diagnostics
// without going through the rule catalog: the parser (schema
// validation), the resolver (reference resolution), and the
// runner (infrastructure findings). Their IDs are plain string
// literals in the code, so the test reads them from the sources.
var pipelinePackages = []string{"parser", "resolver", "runner"}

var diagnosticIDPattern = regexp.MustCompile(`^esdm/[a-z]+/[a-z-]+$`)

// TestLinterRulesPageMentionsEveryPipelineDiagnostic makes sure the
// core Linter Rules page has an entry for every diagnostic the
// pipeline can emit outside the rule catalog. The catalog rules
// are covered by the generated snippets; these diagnostics have no
// Meta() to generate from, so their entries are hand-written and
// this test is what keeps them complete. The entry must use the
// same heading shape as the generated ones, anchor included, so
// that every ID the linter prints can be linked the same way.
func TestLinterRulesPageMentionsEveryPipelineDiagnostic(t *testing.T) {
	t.Parallel()

	catalogIDs := map[string]bool{}
	for _, rule := range rules.Catalog() {
		catalogIDs[rule.Meta().ID] = true
	}

	ids := collectDiagnosticIDs(t)
	require.NotEmpty(t, ids, "no diagnostic IDs found in the pipeline packages")

	markdown := readReference(t, filepath.Join("docs", "reference", "linter-rules.md"))

	for _, id := range ids {
		if catalogIDs[id] {
			continue
		}
		t.Run(id, func(t *testing.T) {
			t.Parallel()

			meta := rules.Meta{ID: id}
			heading := "### `" + id + "` { #" + meta.Anchor() + " }"
			if !strings.Contains(markdown, heading) {
				t.Errorf("reference/linter-rules.md has no entry %q", heading)
			}
		})
	}
}

// collectDiagnosticIDs returns every string literal of the form
// esdm/<category>/<name> found in the non-test sources of the
// pipeline packages, sorted and without duplicates. The files are
// parsed one by one; build tags play no role here, because the
// pipeline packages have none.
func collectDiagnosticIDs(t *testing.T) []string {
	t.Helper()

	found := map[string]struct{}{}
	fileSet := token.NewFileSet()

	for _, packageName := range pipelinePackages {
		directory := filepath.Join("..", packageName)
		entries, err := os.ReadDir(directory)
		require.NoError(t, err)

		for _, entry := range entries {
			name := entry.Name()
			if entry.IsDir() || !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
				continue
			}

			file, err := parser.ParseFile(fileSet, filepath.Join(directory, name), nil, 0)
			require.NoError(t, err)

			ast.Inspect(file, func(node ast.Node) bool {
				literal, isLiteral := node.(*ast.BasicLit)
				if !isLiteral || literal.Kind != token.STRING {
					return true
				}
				value, err := strconv.Unquote(literal.Value)
				if err != nil {
					return true
				}
				if diagnosticIDPattern.MatchString(value) {
					found[value] = struct{}{}
				}
				return true
			})
		}
	}

	return sortedSet(found)
}
