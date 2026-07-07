package reporter

import (
	"encoding/json"
	"io"

	"github.com/thenativeweb/esdm/diag"
)

// JSONFormatter renders diagnostics as a JSON array.
//
// The schema is stable and intended for machine
// consumption (CI, editor integration, tests):
//
//	[
//	  {
//	    "ruleId":   "esdm/naming/event-past-tense",
//	    "severity": "warning",
//	    "message":  "...",
//	    "location": {"file": "a.esdm.yaml", "line": 3, "column": 7},
//	    "related":  [{"message": "...", "location": {...}}]
//	  },
//	  ...
//	]
type JSONFormatter struct{}

// NewJSONFormatter returns a new JSONFormatter.
func NewJSONFormatter() *JSONFormatter {
	return &JSONFormatter{}
}

type jsonLocation struct {
	File   string `json:"file"`
	Line   int    `json:"line"`
	Column int    `json:"column"`
}

type jsonRelated struct {
	Message  string       `json:"message"`
	Location jsonLocation `json:"location"`
}

type jsonDiagnostic struct {
	RuleID   string        `json:"ruleId"`
	Severity string        `json:"severity"`
	Message  string        `json:"message"`
	Location jsonLocation  `json:"location"`
	Related  []jsonRelated `json:"related,omitempty"`
}

// Format writes the diagnostics as a pretty-printed JSON
// array to w. An empty slice is rendered as [].
func (f *JSONFormatter) Format(w io.Writer, diagnostics []diag.Diagnostic) error {
	out := make([]jsonDiagnostic, 0, len(diagnostics))
	for _, d := range diagnostics {
		entry := jsonDiagnostic{
			RuleID:   d.RuleID,
			Severity: d.Severity.String(),
			Message:  d.Message,
			Location: jsonLocation{
				File:   d.Location.File,
				Line:   d.Location.Line,
				Column: d.Location.Column,
			},
		}

		if len(d.Related) > 0 {
			entry.Related = make([]jsonRelated, 0, len(d.Related))
			for _, r := range d.Related {
				entry.Related = append(entry.Related, jsonRelated{
					Message: r.Message,
					Location: jsonLocation{
						File:   r.Location.File,
						Line:   r.Location.Line,
						Column: r.Location.Column,
					},
				})
			}
		}

		out = append(out, entry)
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")

	return enc.Encode(out)
}
