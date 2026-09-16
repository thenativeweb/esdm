package tree

import (
	"fmt"
	"sort"
	"strings"

	"github.com/thenativeweb/esdm/ast"
)

// scopeText pulls scope.<field> as a string, returning
// "" when the field is missing or not a scalar.
func scopeText(scope ast.Node, field string) string {
	v, _ := scope.Field(field).Text()
	return v
}

// nameLocation converts a view's `name` AST location
// into the package-local SourceLocation that nodes
// carry. Used by the builder so each rendered node knows
// where in the source it lives.
func nameLocation(view interface{ Name() ast.Node }) SourceLocation {
	loc := view.Name().Location()
	return SourceLocation{File: loc.File, Line: loc.Line, Column: loc.Column}
}

// sortedByBareName returns the values of a Model map
// sorted by their bare entity name. Used for the
// view command's deterministic top-down rendering.
func sortedByBareName[V interface{ Name() ast.Node }](m map[string]V) []V {
	out := make([]V, 0, len(m))
	for _, v := range m {
		out = append(out, v)
	}
	sort.Slice(out, func(i, j int) bool {
		ni, _ := out[i].Name().Text()
		nj, _ := out[j].Name().Text()
		return ni < nj
	})
	return out
}

// joinComma renders a list of names with comma + space
// separators, suitable for inline annotation slots.
func joinComma(values []string) string {
	return strings.Join(values, ", ")
}

// plural renders a count alongside a noun. The noun is
// always written in its given form - stats use compact
// abbreviations (cmd, evt, rm) where adding an English
// "s" would be awkward, so plural is kept invariant on
// purpose. Returns "" for a zero count so callers can
// chain stats without explicit "if" branches.
func plural(count int, noun string) string {
	if count == 0 {
		return ""
	}
	return fmt.Sprintf("%d %s", count, noun)
}

// appendStat joins two stat fragments with the centered
// dot separator used throughout the renderer.
func appendStat(left, right string) string {
	if left == "" {
		return right
	}
	if right == "" {
		return left
	}
	return left + " · " + right
}

// schemaSummary turns a JSON-Schema-shaped node into a
// short, parenthesized list of property names so the
// renderer can show "data: {playerId, courseId, ...}"
// for compact inline use without dumping full schemas.
func schemaSummary(schemaNode ast.Node) string {
	properties := schemaNode.Field("properties")
	if !properties.Exists() {
		return ""
	}
	var names []string
	for _, entry := range properties.Entries() {
		if k, ok := entry.Key.Text(); ok {
			names = append(names, k)
		}
	}
	if len(names) == 0 {
		return ""
	}
	sort.Strings(names)
	return "{" + strings.Join(names, ", ") + "}"
}

// narrow walks the render tree along the given path
// segments and returns the matching subtree. The path's
// first segment must match a domain; subsequent segments
// follow the natural containment hierarchy. An unknown
// segment returns an error.
//
// A segment is matched by name alone, and a name is not
// unique at a position today: an entity and an actor may
// both be called `customer`, an aggregate and a DCB may
// both be called `order`. The walk therefore keeps every
// match at each depth and continues below all of them,
// so that no sibling is hidden by the accident of which
// kind the builder appends first. Several final matches
// are returned under a synthetic root, which the renderer
// shows transparently - the collision becomes visible in
// the output instead of silently resolving to one side.
func narrow(root *Node, segments []string) (*Node, error) {
	frontier := []*Node{root}
	for i, seg := range segments {
		var matches []*Node
		for _, parent := range frontier {
			for _, child := range parent.Children {
				if child.Name == seg {
					matches = append(matches, child)
				}
			}
		}
		if len(matches) == 0 {
			matched := strings.Join(segments[:i], "/")
			if matched == "" {
				return nil, fmt.Errorf("no element %q under model root", seg)
			}
			return nil, fmt.Errorf("no element %q under %q", seg, matched)
		}
		frontier = matches
	}
	if len(frontier) == 1 {
		return frontier[0], nil
	}
	return &Node{Kind: "model", Children: frontier}, nil
}
