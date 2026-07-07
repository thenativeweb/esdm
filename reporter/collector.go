package reporter

import (
	"sort"
	"sync"

	"github.com/thenativeweb/esdm/diag"
)

// Collector is a thread-safe implementation of
// diag.Reporter that accumulates Diagnostics during a
// linter run.
type Collector struct {
	mu          sync.Mutex
	diagnostics []diag.Diagnostic
}

// NewCollector returns an empty Collector ready to accept
// Diagnostics from concurrent producers.
func NewCollector() *Collector {
	return &Collector{}
}

// Report appends a Diagnostic to the Collector. It is safe
// to call concurrently.
func (c *Collector) Report(d diag.Diagnostic) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.diagnostics = append(c.diagnostics, d)
}

// All returns all collected Diagnostics in deterministic
// order: first by file name, then by line, then by column,
// finally by rule ID. The returned slice is a copy - the
// caller may modify it freely.
func (c *Collector) All() []diag.Diagnostic {
	c.mu.Lock()
	defer c.mu.Unlock()

	out := make([]diag.Diagnostic, len(c.diagnostics))
	copy(out, c.diagnostics)

	sort.SliceStable(out, func(i, j int) bool {
		a, b := out[i], out[j]

		if a.Location.File != b.Location.File {
			return a.Location.File < b.Location.File
		}
		if a.Location.Line != b.Location.Line {
			return a.Location.Line < b.Location.Line
		}
		if a.Location.Column != b.Location.Column {
			return a.Location.Column < b.Location.Column
		}

		return a.RuleID < b.RuleID
	})

	return out
}

// HasErrors reports whether any collected Diagnostic has
// Severity SeverityError. Callers use this to decide the
// process exit code.
func (c *Collector) HasErrors() bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, d := range c.diagnostics {
		if d.Severity == diag.SeverityError {
			return true
		}
	}

	return false
}
