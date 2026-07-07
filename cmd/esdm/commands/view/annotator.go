package view

import (
	"github.com/thenativeweb/esdm/diag"
)

// Annotate walks the tree and attaches the matching
// linter severity to each node. A diagnostic matches a
// node when the diagnostic's location coincides with
// the node's name location, OR - for diagnostics anchored
// at a sub-element of a node (e.g. a constraint inside a
// command) - when the diagnostic's file is the node's
// file and no closer node matches.
//
// After direct matching, severity is bubbled upward so a
// node carries at least the worst severity of any
// descendant. This lets the user spot a problem at a
// high level and zoom in.
func Annotate(root *Node, diagnostics []diag.Diagnostic) {
	if root == nil {
		return
	}
	for _, d := range diagnostics {
		target := pickTarget(root, d.Location)
		if target == nil {
			continue
		}
		raise(target, fromDiagSeverity(d.Severity))
	}
	bubble(root)
}

// pickTarget finds the most specific node whose
// location matches the diagnostic's location. Exact
// match on file/line/column is preferred. If no exact
// match is found, the deepest node in the same file
// whose name is at-or-above the diagnostic's line is
// chosen - the AST anchor of the surrounding entity.
func pickTarget(root *Node, loc diag.Location) *Node {
	var exact *Node
	var bestFallback *Node
	bestLine := -1

	visit := func(n *Node) {}
	visit = func(n *Node) {
		if n.Location.File == loc.File {
			if n.Location.Line == loc.Line && n.Location.Column == loc.Column {
				exact = n
			} else if n.Location.Line > 0 && n.Location.Line <= loc.Line && n.Location.Line > bestLine {
				bestLine = n.Location.Line
				bestFallback = n
			}
		}
		for _, c := range n.Children {
			visit(c)
		}
	}
	visit(root)

	if exact != nil {
		return exact
	}
	return bestFallback
}

// raise lifts a node's severity to at least the given
// severity - warning never overrides error.
func raise(n *Node, s Severity) {
	if priority(s) > priority(n.Severity) {
		n.Severity = s
	}
}

// bubble propagates the worst severity of a subtree's
// descendants up to its root, recursively.
func bubble(n *Node) Severity {
	worst := n.Severity
	for _, c := range n.Children {
		s := bubble(c)
		if priority(s) > priority(worst) {
			worst = s
		}
	}
	n.Severity = worst
	return worst
}

func priority(s Severity) int {
	switch s {
	case SeverityError:
		return 2
	case SeverityWarning:
		return 1
	default:
		return 0
	}
}

func fromDiagSeverity(s diag.Severity) Severity {
	switch s {
	case diag.SeverityError:
		return SeverityError
	case diag.SeverityWarning:
		return SeverityWarning
	default:
		return SeverityNone
	}
}
