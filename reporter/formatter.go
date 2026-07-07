package reporter

import (
	"io"

	"github.com/thenativeweb/esdm/diag"
)

// Formatter turns a slice of Diagnostics into textual or
// structured output on the provided writer.
type Formatter interface {
	Format(w io.Writer, diagnostics []diag.Diagnostic) error
}
