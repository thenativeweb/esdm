package cmdutils

import (
	"fmt"
	"io"
	"os"

	"golang.org/x/term"
)

const (
	ColorAuto   = "auto"
	ColorAlways = "always"
	ColorNever  = "never"
)

// ResolveColor turns a --color flag value into a
// concrete boolean. "always" and "never" are
// unconditional; "auto" consults the writer (color only
// when it is a terminal) and the NO_COLOR environment
// variable.
func ResolveColor(mode string, out io.Writer) (bool, error) {
	switch mode {
	case ColorAlways:
		return true, nil
	case ColorNever:
		return false, nil
	case ColorAuto:
		if _, ok := os.LookupEnv("NO_COLOR"); ok {
			return false, nil
		}
		return WriterIsTerminal(out), nil
	default:
		return false, fmt.Errorf("unknown --color value %q (expected %q, %q, or %q)", mode, ColorAuto, ColorAlways, ColorNever)
	}
}

// WriterIsTerminal reports whether w is a *os.File
// attached to a terminal - the standard signal that ANSI
// color codes will render rather than leak into a pipe
// or a file.
func WriterIsTerminal(w io.Writer) bool {
	file, ok := w.(*os.File)
	if !ok {
		return false
	}
	return term.IsTerminal(int(file.Fd()))
}
