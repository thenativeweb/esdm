package view

// Node is one entry in the render tree the view command
// produces. Each node carries the data needed by the
// renderer in one place: a kind/name pair (the heading),
// optional inline tags (e.g. "(core)"), optional
// right-side stats (e.g. a "4 cmd" and "6 evt" count
// pair joined by a middle dot), and optional detail
// lines that the renderer emits only when the
// --with-details flag is set. Children form the
// containment hierarchy. Severity is attached by the
// annotator so the renderer can mark each node with the
// matching glyph.
type Node struct {
	Kind     string
	Name     string
	Tags     []string
	Stats    []string
	Lines    []string
	Children []*Node
	Severity Severity
	Key      string

	// Location is the source location of the node's
	// `name` field. The annotator matches diagnostic
	// locations against it to attach the right severity.
	Location SourceLocation
}

// SourceLocation captures where a node was defined in
// the source files. Mirrors diag.Location structurally
// but kept package-local so the view package does not
// drag the diag package into its core data shapes.
type SourceLocation struct {
	File   string
	Line   int
	Column int
}

// Severity classifies how the renderer marks a node
// against any diagnostic that targeted it. The empty
// value is the unmarked default; warning and error
// match the diag.Severity values produced by the
// linter.
type Severity string

const (
	SeverityNone    Severity = ""
	SeverityWarning Severity = "warning"
	SeverityError   Severity = "error"
)
