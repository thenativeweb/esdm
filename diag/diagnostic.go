package diag

// Related attaches a secondary source position to a
// Diagnostic. It is typically used when a finding refers
// to multiple relevant sites, e.g. "duplicate name defined
// here and there".
type Related struct {
	Message  string
	Location Location
}

// Diagnostic is the single uniform type that travels
// through the esdm linter pipeline. It is produced by the
// parser, resolver, and rules, and consumed by the
// reporter's formatters.
//
// DocumentationURL is the address of the finding's entry
// in the published documentation. Producers leave it
// empty; the runner fills it in for every diagnostic, so
// that formatters can point the reader at the explanation
// without knowing how the documentation is laid out.
type Diagnostic struct {
	RuleID           string
	Severity         Severity
	Message          string
	Location         Location
	Related          []Related
	DocumentationURL string
}

// Reporter collects Diagnostics during a linter run.
// Implementations must be safe for concurrent use.
type Reporter interface {
	Report(Diagnostic)
}
