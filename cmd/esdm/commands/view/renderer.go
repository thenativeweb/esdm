package view

import (
	"strings"
)

// RenderOptions controls how the renderer styles a tree.
// Colors enables ANSI escape sequences for severity
// markings; ShowDetails emits the per-node Lines slice
// alongside the header.
type RenderOptions struct {
	Colors      bool
	ShowDetails bool
}

// Render walks the tree and produces the indented,
// box-drawing text representation. A synthetic root
// (Kind "model") is rendered transparently - its
// children become the top-level lines so the user sees
// real entities at the top.
func Render(root *Node, opts RenderOptions) string {
	if root == nil {
		return ""
	}
	var b strings.Builder
	if root.Kind == "model" {
		for _, c := range root.Children {
			renderNode(&b, c, "", "", opts)
		}
		return b.String()
	}
	renderNode(&b, root, "", "", opts)
	return b.String()
}

// renderNode writes one node's line followed by its
// detail lines and recurses into children. `prefix` is
// the leading text for the node's own line (including
// any tree connectors above); `childPrefix` is the
// indentation to use for everything that descends from
// the node (lines and child nodes).
func renderNode(b *strings.Builder, n *Node, prefix, childPrefix string, opts RenderOptions) {
	b.WriteString(prefix)
	b.WriteString(formatHeader(n, opts))
	b.WriteString("\n")

	if opts.ShowDetails {
		for _, l := range n.Lines {
			b.WriteString(childPrefix)
			b.WriteString("   ")
			b.WriteString(applyDim(l, opts.Colors))
			b.WriteString("\n")
		}
	}

	for i, c := range n.Children {
		isLast := i == len(n.Children)-1
		var connector, nextChildPrefix string
		if isLast {
			connector = "└─ "
			nextChildPrefix = "   "
		} else {
			connector = "├─ "
			nextChildPrefix = "│  "
		}
		renderNode(b, c, childPrefix+connector, childPrefix+nextChildPrefix, opts)
	}
}

// formatHeader assembles the single-line header for a
// node. The visual hierarchy puts the entity *name* on
// top (bold), the *kind* in cyan as the containment
// anchor on its own color spur, and the inline tags
// plus the right-side stats below them in dim text.
// Severity glyphs pick up their color from the same
// opts.Colors flag and stay outside the dim/bold/cyan
// scheme so their red/yellow signal is unaffected.
func formatHeader(n *Node, opts RenderOptions) string {
	var b strings.Builder
	b.WriteString(applyKindColor(n.Kind, opts.Colors))
	b.WriteString(" ")
	b.WriteString(applyBold(n.Name, opts.Colors))
	if len(n.Tags) > 0 {
		b.WriteString(" ")
		b.WriteString(applyDim("("+strings.Join(n.Tags, ", ")+")", opts.Colors))
	}
	if g := severityGlyph(n.Severity, opts.Colors); g != "" {
		b.WriteString(" ")
		b.WriteString(g)
	}
	if len(n.Stats) > 0 {
		b.WriteString("  ")
		b.WriteString(applyDim(strings.Join(n.Stats, " · "), opts.Colors))
	}
	return b.String()
}

// applyBold wraps text in the ANSI bold sequence when
// colors are enabled; otherwise it returns the text
// unchanged. Bold follows the same on/off rule as
// color throughout the renderer.
func applyBold(text string, isColorEnabled bool) string {
	if !isColorEnabled {
		return text
	}
	return "\x1b[1m" + text + "\x1b[0m"
}

// applyDim wraps text in the ANSI dim sequence when
// colors are enabled; otherwise it returns the text
// unchanged. Used for inline tags, right-side stats,
// and detail lines so they recede behind the entity
// name.
func applyDim(text string, isColorEnabled bool) string {
	if !isColorEnabled {
		return text
	}
	return "\x1b[2m" + text + "\x1b[0m"
}

// applyKindColor wraps the kind keyword in cyan when
// colors are enabled; otherwise it returns the text
// unchanged. Cyan was chosen because it sits in a
// different visual register from the bold-white name
// and the dim tags/stats, so the kind reads as a
// structural anchor without competing with the
// identity-bearing name.
func applyKindColor(text string, isColorEnabled bool) string {
	if !isColorEnabled {
		return text
	}
	return "\x1b[36m" + text + "\x1b[0m"
}

// severityGlyph returns the marking that goes on a
// node's header line for the given severity. ANSI
// color codes are emitted only when color is enabled;
// the glyph itself is always present so non-colored
// output (pipes, --color never) still surfaces the
// severity.
func severityGlyph(s Severity, isColorEnabled bool) string {
	switch s {
	case SeverityError:
		if isColorEnabled {
			return "\x1b[31m✗\x1b[0m"
		}
		return "✗"
	case SeverityWarning:
		if isColorEnabled {
			return "\x1b[33m⚠\x1b[0m"
		}
		return "⚠"
	default:
		return ""
	}
}
