package reporter

import (
	"fmt"
	"io"

	"github.com/thenativeweb/esdm/diag"
)

// HumanFormatter renders diagnostics as multi-line blocks
// separated by blank lines, optionally using ANSI color
// codes to distinguish severity at a glance:
//
//	error: esdm/structure/unresolved-reference
//	  at /tmp/model.esdm.yaml:29:14
//	  unresolved aggregate "ordr"
//	  note: did you mean "order"? (/tmp/model.esdm.yaml:13:7)
//
//	warning: esdm/modeling/event-name-aggregate-prefix
//	  at /tmp/model.esdm.yaml:47:7
//	  event name "invoice-issued" does not start with its aggregate's name "order"
//
// The Colors field toggles ANSI escape codes. Callers
// decide based on their own TTY detection (or a
// --color=auto|always|never flag) and hand the decision
// here explicitly; the formatter itself stays deterministic
// and does no I/O introspection.
type HumanFormatter struct {
	Colors bool
}

// NewHumanFormatter returns a HumanFormatter with colors
// disabled by default. Turn them on via the Colors field.
func NewHumanFormatter() *HumanFormatter {
	return &HumanFormatter{}
}

// ANSI escape sequences used for the colored variant.
const (
	ansiReset  = "\x1b[0m"
	ansiBold   = "\x1b[1m"
	ansiDim    = "\x1b[2m"
	ansiRed    = "\x1b[31m"
	ansiYellow = "\x1b[33m"
	ansiCyan   = "\x1b[36m"
)

// Format writes each diagnostic as a multi-line block to
// w, followed by a blank line to visually separate it
// from the next.
func (f *HumanFormatter) Format(w io.Writer, diagnostics []diag.Diagnostic) error {
	for i, d := range diagnostics {
		if i > 0 {
			_, err := fmt.Fprintln(w)
			if err != nil {
				return err
			}
		}

		err := f.writeDiagnostic(w, d)
		if err != nil {
			return err
		}
	}

	return nil
}

func (f *HumanFormatter) writeDiagnostic(w io.Writer, d diag.Diagnostic) error {
	label, color := severityLabelAndColor(d.Severity)

	_, err := fmt.Fprintf(w, "%s: %s\n", f.paint(label, color, ansiBold), f.paint(d.RuleID, ansiBold))
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(w, "  %s %s\n", f.paint("at", ansiDim), f.paint(locationPrefix(d.Location), ansiDim))
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(w, "  %s\n", d.Message)
	if err != nil {
		return err
	}

	for _, r := range d.Related {
		_, err = fmt.Fprintf(w, "  %s %s %s\n",
			f.paint("note:", ansiCyan),
			r.Message,
			f.paint("("+locationPrefix(r.Location)+")", ansiDim),
		)
		if err != nil {
			return err
		}
	}

	return nil
}

// paint wraps s in the given ANSI escape codes when
// Colors is enabled. When disabled, it returns s
// unchanged so the colored and plain branches share the
// same layout.
func (f *HumanFormatter) paint(s string, codes ...string) string {
	if !f.Colors || len(codes) == 0 {
		return s
	}

	prefix := ""
	for _, c := range codes {
		prefix += c
	}

	return prefix + s + ansiReset
}

// severityLabelAndColor maps a severity to its display
// word and the ANSI color code used for it.
func severityLabelAndColor(s diag.Severity) (label, color string) {
	switch s {
	case diag.SeverityError:
		return "error", ansiRed
	case diag.SeverityWarning:
		return "warning", ansiYellow
	case diag.SeverityInfo:
		return "info", ansiCyan
	case diag.SeverityHint:
		return "hint", ansiDim
	default:
		return s.String(), ""
	}
}

// locationPrefix renders a Location as the canonical
// "file:line:col" prefix, or "<internal>" for zero
// locations.
func locationPrefix(loc diag.Location) string {
	if loc.IsZero() {
		return "<internal>"
	}
	return fmt.Sprintf("%s:%d:%d", loc.File, loc.Line, loc.Column)
}
