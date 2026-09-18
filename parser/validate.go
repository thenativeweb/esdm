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
func validationDiagnostics(err *jsonschema.ValidationError, document ast.Node, schema *compiledSchema) []diag.Diagnostic {
	context := newValidationContext(err, document, schema)

	var out []diag.Diagnostic
	context.flatten(err, &out)

	return out
}

// validationContext carries what translating one
// ValidationError tree needs beyond the error itself: the
// document the errors are about, the schema they were
// produced against, and every schema position the tree
// points at. The last one is what tells a defect apart
// from the follow-on errors it drags along.
type validationContext struct {
	document ast.Node
	schema   *compiledSchema
	failures []schemaFailure
}

// schemaFailure reduces one node of the ValidationError
// tree to the two coordinates the follow-on analysis
// works with: where the failure sits in the schema, and
// which part of the document it is about.
type schemaFailure struct {
	pointer          string
	instanceLocation []string
}

func newValidationContext(err *jsonschema.ValidationError, document ast.Node, schema *compiledSchema) *validationContext {
	context := &validationContext{
		document: document,
		schema:   schema,
	}
	context.collectFailures(err)

	return context
}

func (c *validationContext) collectFailures(err *jsonschema.ValidationError) {
	pointer, isOwn := c.pointerOf(err.SchemaURL)
	if isOwn {
		c.failures = append(c.failures, schemaFailure{
			pointer:          pointer,
			instanceLocation: err.InstanceLocation,
		})
	}

	for _, cause := range err.Causes {
		c.collectFailures(cause)
	}
}

// pointerOf reduces a schema URL to the JSON pointer it
// addresses inside the schema document. It reports false
// for a URL belonging to another document, which the
// analysis cannot reason about.
func (c *validationContext) pointerOf(schemaURL string) (string, bool) {
	return strings.CutPrefix(schemaURL, c.schema.BaseURL+"#")
}

// flatten walks the ValidationError tree and emits one or
// more Diagnostics for each leaf error. Wrapper kinds
// (Schema, Group, AllOf, AnyOf, OneOf, Not) are traversed
// rather than reported: their nested causes carry the
// actionable information. A OneOf is traversed
// selectively, and an unevaluated property is dropped
// when it turns out to be follow-on noise.
func (c *validationContext) flatten(err *jsonschema.ValidationError, out *[]diag.Diagnostic) {
	switch err.ErrorKind.(type) {
	case *kind.Schema, *kind.Group, *kind.AllOf, *kind.AnyOf, *kind.Not:
		c.flattenCauses(err.Causes, out)
		return
	case *kind.OneOf:
		causes := err.Causes

		distinguished := c.distinguishedCause(err)
		if distinguished != nil {
			causes = []*jsonschema.ValidationError{distinguished}
		}

		c.flattenCauses(causes, out)
		return
	}

	switch errorKind := err.ErrorKind.(type) {
	case *kind.Required:
		parent := locate(c.document, err.InstanceLocation)
		for _, missing := range errorKind.Missing {
			*out = append(*out, diag.Diagnostic{
				RuleID:   "esdm/structure/missing-required-field",
				Severity: diag.SeverityError,
				Message:  fmt.Sprintf("missing required field %q", missing),
				Location: parent.Location(),
			})
		}
	case *kind.AdditionalProperties:
		parent := locate(c.document, err.InstanceLocation)
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
		target := locate(c.document, err.InstanceLocation)
		*out = append(*out, diag.Diagnostic{
			RuleID:   "esdm/structure/type-mismatch",
			Severity: diag.SeverityError,
			Message:  fmt.Sprintf("expected %v, got %s", errorKind.Want, errorKind.Got),
			Location: target.Location(),
		})
	case *kind.Enum:
		target := locate(c.document, err.InstanceLocation)
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
			if c.isFollowOnUnknownField(err) {
				return
			}

			target := locate(c.document, err.InstanceLocation)
			fieldName := err.InstanceLocation[len(err.InstanceLocation)-1]
			*out = append(*out, diag.Diagnostic{
				RuleID:   "esdm/structure/unknown-field",
				Severity: diag.SeverityError,
				Message:  fmt.Sprintf("unknown field %q", fieldName),
				Location: target.Location(),
			})
			return
		}

		target := locate(c.document, err.InstanceLocation)
		*out = append(*out, diag.Diagnostic{
			RuleID:   "esdm/structure/constraint-violation",
			Severity: diag.SeverityError,
			Message:  err.ErrorKind.LocalizedString(defaultPrinter),
			Location: target.Location(),
		})
	default:
		if len(err.Causes) > 0 {
			c.flattenCauses(err.Causes, out)
			return
		}

		target := locate(c.document, err.InstanceLocation)
		*out = append(*out, diag.Diagnostic{
			RuleID:   "esdm/structure/constraint-violation",
			Severity: diag.SeverityError,
			Message:  err.ErrorKind.LocalizedString(defaultPrinter),
			Location: target.Location(),
		})
	}
}

func (c *validationContext) flattenCauses(causes []*jsonschema.ValidationError, out *[]diag.Diagnostic) {
	for _, cause := range causes {
		c.flatten(cause, out)
	}
}

// distinguishedCause picks the cause of the one `oneOf`
// branch the document aims at, or nil when the document
// identifies none of them. Without this, a document that
// gets one field of one branch wrong is reported against
// every branch, so most of the findings describe
// alternatives the document never intended.
//
// A branch is identified by how much of it the document
// already exhibits: every required property that is
// present and every discriminating constant it matches
// counts for one. A tie means the document resembles the
// branches equally - reporting all of them is then the
// accurate answer, not a fallback.
func (c *validationContext) distinguishedCause(err *jsonschema.ValidationError) *jsonschema.ValidationError {
	base, isOwn := c.pointerOf(err.SchemaURL)
	if !isOwn {
		return nil
	}

	instance := locate(c.document, err.InstanceLocation)
	if instance.Kind() != ast.KindMapping {
		return nil
	}

	branchPrefix := base + "/oneOf/"

	var distinguished *jsonschema.ValidationError
	highestAffinity := 0
	isTied := false

	for _, cause := range err.Causes {
		pointer, isOwn := c.pointerOf(cause.SchemaURL)
		if !isOwn {
			return nil
		}

		rest, isBranch := strings.CutPrefix(pointer, branchPrefix)
		if !isBranch {
			return nil
		}

		index, _, _ := strings.Cut(rest, "/")
		branch, isResolved := subschemaAt(c.schema.Document, branchPrefix+index)
		if !isResolved {
			return nil
		}

		affinity := branchAffinity(branch, instance)
		switch {
		case affinity > highestAffinity:
			distinguished = cause
			highestAffinity = affinity
			isTied = false
		case affinity == highestAffinity:
			isTied = true
		}
	}

	if isTied || highestAffinity == 0 {
		return nil
	}

	return distinguished
}

// branchAffinity counts how much of a `oneOf` branch the
// instance already exhibits. Constants weigh in alongside
// required properties because a constant is what still
// identifies a variant when its required fields are
// incomplete, which is exactly the situation that
// produced the error.
func branchAffinity(branch map[string]any, instance ast.Node) int {
	affinity := 0

	for _, name := range requiredProperties(branch) {
		if instance.HasField(name) {
			affinity++
		}
	}

	for name, want := range constrainedConstants(branch) {
		got, isText := instance.Field(name).Text()
		if isText && got == want {
			affinity++
		}
	}

	return affinity
}

// isFollowOnUnknownField reports whether an unevaluated
// property is fallout of another failure rather than a
// defect of its own.
//
// The core schema dispatches kinds through `allOf` with
// `if`/`then` and closes the document with
// `unevaluatedProperties: false`. When a branch fails for
// any reason, its annotations are dropped, so every
// property that branch would have evaluated counts as
// unevaluated and looks like an unknown field. The
// property is fallout exactly when some failure below the
// same schema position sits in a branch that declares it;
// a genuine typo appears in no branch and is still
// reported.
func (c *validationContext) isFollowOnUnknownField(err *jsonschema.ValidationError) bool {
	pointer, isOwn := c.pointerOf(err.SchemaURL)
	if !isOwn {
		return false
	}

	// Only `unevaluatedProperties` loses annotations this
	// way. An `additionalProperties` failure states that a
	// property is unknown right where it stands, with no
	// branch involved that could have evaluated it.
	base, isUnevaluated := strings.CutSuffix(pointer, "/unevaluatedProperties")
	if !isUnevaluated {
		return false
	}

	property := err.InstanceLocation[len(err.InstanceLocation)-1]
	owner := err.InstanceLocation[:len(err.InstanceLocation)-1]

	for _, failure := range c.failures {
		if !isWithin(owner, failure.instanceLocation) {
			continue
		}

		declared := propertiesBelow(c.schema.Document, base, failure.pointer)
		_, isDeclared := declared[property]
		if isDeclared {
			return true
		}
	}

	return false
}

// isWithin reports whether location addresses the value
// at prefix, or something nested inside it.
func isWithin(prefix []string, location []string) bool {
	if len(location) < len(prefix) {
		return false
	}

	for i, segment := range prefix {
		if location[i] != segment {
			return false
		}
	}

	return true
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
