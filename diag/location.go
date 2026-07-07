package diag

// Location points at a specific position in a source file.
// Line and Column are 1-based. Column is measured in bytes.
// A zero Location (empty File, Line and Column == 0) means
// "no meaningful source location" - used for internal or
// system-level diagnostics that do not originate from a user
// file.
type Location struct {
	File   string
	Line   int
	Column int
}

// IsZero reports whether the Location carries no meaningful
// source information.
func (l Location) IsZero() bool {
	return l.File == "" && l.Line == 0 && l.Column == 0
}
