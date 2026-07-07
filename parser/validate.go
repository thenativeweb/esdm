package parser

import (
	"fmt"
	"strings"

	"github.com/santhosh-tekuri/jsonschema/v6"
	"github.com/santhosh-tekuri/jsonschema/v6/kind"
	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/thenativeweb/esdm/ast"
	"github.com/thenativeweb/esdm/diag"
	"github.com/thenativeweb/esdm/hint"
)

// defaultPrinter renders ErrorKind localized strings in
// English. The message package panics when a nil printer
// is passed to LocalizedString, so we always hand it one.
var defaultPrinter = message.NewPrinter(language.English)

// validationDiagnostics translates a ValidationError tree
// into a flat list of Diagnostics, each pinned to the
// right position inside the parsed document.
func validationDiagnostics(err *jsonschema.ValidationError, document ast.Node) []diag.Diagnostic {
	var out []diag.Diagnostic
	flatten(err, document, &out)
	return out
}

// flatten walks the ValidationError tree and emits one or
// more Diagnostics for each leaf error. Wrapper kinds
// (Schema, Group, AllOf, AnyOf, OneOf, Not) are traversed
// rather than reported: their nested causes carry the
// actionable information.
func flatten(err *jsonschema.ValidationError, document ast.Node, out *[]diag.Diagnostic) {
	switch err.ErrorKind.(type) {
	case *kind.Schema, *kind.Group, *kind.AllOf, *kind.AnyOf, *kind.OneOf, *kind.Not:
		for _, cause := range err.Causes {
			flatten(cause, document, out)
		}
		return
	}

	switch errorKind := err.ErrorKind.(type) {
	case *kind.Required:
		parent := locate(document, err.InstanceLocation)
		for _, missing := range errorKind.Missing {
			*out = append(*out, diag.Diagnostic{
				RuleID:   "esdm/structure/missing-required-field",
				Severity: diag.SeverityError,
				Message:  fmt.Sprintf("missing required field %q", missing),
				Location: parent.Location(),
			})
		}
	case *kind.AdditionalProperties:
		parent := locate(document, err.InstanceLocation)
		for _, extra := range errorKind.Properties {
			target := parent.Field(extra)
			location := target.Location()
			if location.IsZero() {
				location = parent.Location()
			}

			*out = append(*out, diag.Diagnostic{
				RuleID:   "esdm/structure/unknown-field",
				Severity: diag.SeverityError,
				Message:  fmt.Sprintf("unknown field %q", extra),
				Location: location,
			})
		}
	case *kind.Type:
		target := locate(document, err.InstanceLocation)
		*out = append(*out, diag.Diagnostic{
			RuleID:   "esdm/structure/type-mismatch",
			Severity: diag.SeverityError,
			Message:  fmt.Sprintf("expected %v, got %s", errorKind.Want, errorKind.Got),
			Location: target.Location(),
		})
	case *kind.Enum:
		target := locate(document, err.InstanceLocation)
		diagnostic := diag.Diagnostic{
			RuleID:   "esdm/structure/constraint-violation",
			Severity: diag.SeverityError,
			Message:  err.ErrorKind.LocalizedString(defaultPrinter),
			Location: target.Location(),
		}
		if suggestion := suggestEnumValue(errorKind); suggestion != "" {
			diagnostic.Related = []diag.Related{
				{
					Message:  fmt.Sprintf("did you mean %q?", suggestion),
					Location: target.Location(),
				},
			}
		}
		*out = append(*out, diagnostic)
	case *kind.FalseSchema:
		if isUnknownFieldContext(err.SchemaURL) && len(err.InstanceLocation) > 0 {
			target := locate(document, err.InstanceLocation)
			fieldName := err.InstanceLocation[len(err.InstanceLocation)-1]
			*out = append(*out, diag.Diagnostic{
				RuleID:   "esdm/structure/unknown-field",
				Severity: diag.SeverityError,
				Message:  fmt.Sprintf("unknown field %q", fieldName),
				Location: target.Location(),
			})
			return
		}

		target := locate(document, err.InstanceLocation)
		*out = append(*out, diag.Diagnostic{
			RuleID:   "esdm/structure/constraint-violation",
			Severity: diag.SeverityError,
			Message:  err.ErrorKind.LocalizedString(defaultPrinter),
			Location: target.Location(),
		})
	default:
		if len(err.Causes) > 0 {
			for _, cause := range err.Causes {
				flatten(cause, document, out)
			}
			return
		}

		target := locate(document, err.InstanceLocation)
		*out = append(*out, diag.Diagnostic{
			RuleID:   "esdm/structure/constraint-violation",
			Severity: diag.SeverityError,
			Message:  err.ErrorKind.LocalizedString(defaultPrinter),
			Location: target.Location(),
		})
	}
}

// suggestEnumValue inspects an Enum validation failure
// and returns a close string candidate when both the
// provided value and the enum alternatives are strings
// and the nearest alternative sits within the
// suggestion threshold. Returns "" when no suggestion is
// warranted.
func suggestEnumValue(enumError *kind.Enum) string {
	got, ok := enumError.Got.(string)
	if !ok {
		return ""
	}

	candidates := make([]string, 0, len(enumError.Want))
	for _, w := range enumError.Want {
		s, ok := w.(string)
		if !ok {
			return ""
		}
		candidates = append(candidates, s)
	}

	best, ok := hint.Best(got, candidates)
	if !ok {
		return ""
	}
	return best
}

// isUnknownFieldContext reports whether a FalseSchema
// error originates from a schema keyword that represents
// "no unknown properties allowed" - either
// additionalProperties or unevaluatedProperties. The
// JSON Schema error infrastructure does not expose a
// dedicated ErrorKind for unevaluatedProperties, so we
// recognize it by inspecting the failing subschema URL.
func isUnknownFieldContext(schemaURL string) bool {
	return strings.Contains(schemaURL, "/unevaluatedProperties") ||
		strings.Contains(schemaURL, "/additionalProperties")
}

// locate follows an InstanceLocation - a slice of raw,
// unescaped path segments - into the document and returns
// the resolved node, or the last reachable node on the
// way.
func locate(document ast.Node, instanceLocation []string) ast.Node {
	current := document
	for _, segment := range instanceLocation {
		var next ast.Node
		switch current.Kind() {
		case ast.KindMapping:
			next = current.Field(segment)
		case ast.KindSequence:
			var index int
			_, err := fmt.Sscanf(segment, "%d", &index)
			if err != nil {
				return current
			}
			next = current.At(index)
		default:
			return current
		}

		if !next.Exists() {
			return current
		}
		current = next
	}
	return current
}
