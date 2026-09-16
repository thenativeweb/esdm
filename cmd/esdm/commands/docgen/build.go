package docgen

import (
	"sort"

	"github.com/thenativeweb/esdm/diag"
	"github.com/thenativeweb/esdm/model"
	"github.com/thenativeweb/esdm/modelpath"
	"github.com/thenativeweb/esdm/tree"
)

// Page is one rendered Markdown file: its slash-separated
// path inside the output directory and its content.
type Page struct {
	Path    string
	Content string
}

// generator holds what rendering a page needs to know about
// every other page: where each node's page lives, how each
// node is referenced, which node a model document belongs
// to, and which nodes are being written at all.
type generator struct {
	m          *model.Model
	root       *tree.Node
	pathOf     map[*tree.Node]string
	refOf      map[*tree.Node]string
	byLocation map[diag.Location]*tree.Node
	selected   map[*tree.Node]bool
}

// Build renders the pages for the model's tree, narrowed to
// the region the path selects. The empty path renders every
// element plus the root index; a path renders the matching
// subtree with the full paths its pages would have in the
// complete tree, so that references keep resolving.
func Build(m *model.Model, p modelpath.Path) ([]Page, error) {
	root := tree.Build(m, true)
	g := &generator{
		m:          m,
		root:       root,
		pathOf:     make(map[*tree.Node]string),
		refOf:      make(map[*tree.Node]string),
		byLocation: make(map[diag.Location]*tree.Node),
		selected:   make(map[*tree.Node]bool),
	}
	for _, domain := range root.Children {
		g.index(domain, nil)
	}

	target, err := tree.Narrow(root, p)
	if err != nil {
		return nil, err
	}
	g.select_(target)

	var pages []Page
	seen := make(map[string]bool)
	if target == root {
		pages = append(pages, Page{Path: "README.md", Content: g.renderRoot()})
	}
	g.walk(root, func(n *tree.Node) {
		if !g.selected[n] || seen[g.pathOf[n]] {
			return
		}
		seen[g.pathOf[n]] = true
		pages = append(pages, Page{Path: g.pathOf[n], Content: g.renderPage(n)})
	})
	sort.Slice(pages, func(i, j int) bool {
		return pages[i].Path < pages[j].Path
	})
	return pages, nil
}

// index assigns every node its segments, and from them its
// page path and reference. A context mapping has no domain,
// so its segments start over at itself regardless of the
// domain the tree lists it under; the same mapping listed
// under two domains therefore maps to one page.
func (g *generator) index(n *tree.Node, parent []segment) {
	var segments []segment
	if n.Kind == "context-mapping" {
		segments = []segment{{kind: n.Kind, name: n.Name}}
	} else {
		segments = append(append([]segment{}, parent...), segment{kind: n.Kind, name: n.Name})
	}
	g.pathOf[n] = pagePath(segments)
	g.refOf[n] = reference(segments)
	if n.Document.Exists() {
		g.byLocation[n.Document.Location()] = n
	}
	for _, child := range n.Children {
		g.index(child, segments)
	}
}

// select_ marks the nodes to write: every node beneath the
// narrowed target. A synthetic root - the whole model, or
// several same-named matches - contributes its children.
func (g *generator) select_(target *tree.Node) {
	if target.Kind == "model" {
		for _, child := range target.Children {
			g.walk(child, func(n *tree.Node) { g.selected[n] = true })
		}
		return
	}
	g.walk(target, func(n *tree.Node) { g.selected[n] = true })
}

func (g *generator) walk(n *tree.Node, visit func(*tree.Node)) {
	if n.Kind != "model" {
		visit(n)
	}
	for _, child := range n.Children {
		g.walk(child, visit)
	}
}

// nodeFor returns the tree node of a model document, if the
// tree placed it.
func (g *generator) nodeFor(document model.DocumentViewBase) (*tree.Node, bool) {
	n, ok := g.byLocation[document.Location()]
	return n, ok
}
