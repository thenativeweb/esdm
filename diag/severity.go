package diag

// Severity classifies a Diagnostic by how impactful it is.
// Error-level diagnostics cause the linter to exit with a
// non-zero status; Warning, Info, and Hint do not.
type Severity int

const (
	SeverityError Severity = iota
	SeverityWarning
	SeverityInfo
	SeverityHint
)

func (s Severity) String() string {
	switch s {
	case SeverityError:
		return "error"
	case SeverityWarning:
		return "warning"
	case SeverityInfo:
		return "info"
	case SeverityHint:
		return "hint"
	default:
		return "unknown"
	}
}
