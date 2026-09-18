package parser

import (
	"strconv"
	"strings"
)

// instanceDescendingKeywords are the schema keywords
// whose subschemas describe a different value than the
// schema they sit in. A walk through the schema stops at
// them, because from there on the declarations are about
// something other than the value the walk started at.
// `$defs` and `definitions` are listed for a related
// reason: a definition is reached by a `$ref` from
// elsewhere, so the pointer leading into it says nothing
// about the value at hand.
var instanceDescendingKeywords = map[string]struct{}{
	"properties":            {},
	"patternProperties":     {},
	"additionalProperties":  {},
	"unevaluatedProperties": {},
	"propertyNames":         {},
	"items":                 {},
	"prefixItems":           {},
	"additionalItems":       {},
	"unevaluatedItems":      {},
	"contains":              {},
	"$defs":                 {},
	"definitions":           {},
}

// maximumReferenceDepth bounds how far a chain of `$ref`s
// is followed. A schema that pointed a `$ref` back at
// itself would otherwise loop forever. The embedded
// schemas are verified at build time, so the bound is a
// safeguard against a malformed schema, not a supported
// nesting depth.
const maximumReferenceDepth = 16

// pointerSegments splits a JSON pointer into its
// unescaped segments. The empty pointer addresses the
// document root and yields no segments.
func pointerSegments(pointer string) []string {
	if pointer == "" {
		return nil
	}

	raw := strings.Split(strings.TrimPrefix(pointer, "/"), "/")

	segments := make([]string, 0, len(raw))
	for _, segment := range raw {
		segment = strings.ReplaceAll(segment, "~1", "/")
		segment = strings.ReplaceAll(segment, "~0", "~")
		segments = append(segments, segment)
	}

	return segments
}

// step follows one JSON pointer segment into a decoded
// schema value: a key of a mapping, or an index of a
// sequence.
func step(value any, segment string) (any, bool) {
	switch typed := value.(type) {
	case map[string]any:
		next, isPresent := typed[segment]
		return next, isPresent
	case []any:
		index, err := strconv.Atoi(segment)
		if err != nil {
			return nil, false
		}
		if index < 0 || index >= len(typed) {
			return nil, false
		}

		return typed[index], true
	default:
		return nil, false
	}
}

// valueAt resolves a JSON pointer against a decoded
// schema document.
func valueAt(document any, pointer string) (any, bool) {
	current := document
	for _, segment := range pointerSegments(pointer) {
		next, isPresent := step(current, segment)
		if !isPresent {
			return nil, false
		}

		current = next
	}

	return current, true
}

// dereference follows a chain of local `$ref`s to the
// subschema that carries the actual constraints, so
// callers see what a branch declares rather than the
// reference to it. A `$ref` into another document is left
// alone: the linter compiles one document per apiVersion,
// so there is nothing to follow it into.
func dereference(document any, subschema map[string]any) map[string]any {
	for range maximumReferenceDepth {
		reference, isString := subschema["$ref"].(string)
		if !isString {
			return subschema
		}

		pointer, isLocal := strings.CutPrefix(reference, "#")
		if !isLocal {
			return subschema
		}

		value, isPresent := valueAt(document, pointer)
		if !isPresent {
			return subschema
		}

		referenced, isMapping := value.(map[string]any)
		if !isMapping {
			return subschema
		}

		subschema = referenced
	}

	return subschema
}

// subschemaAt resolves a JSON pointer to the subschema it
// addresses, with local `$ref`s followed.
func subschemaAt(document any, pointer string) (map[string]any, bool) {
	value, isPresent := valueAt(document, pointer)
	if !isPresent {
		return nil, false
	}

	subschema, isMapping := value.(map[string]any)
	if !isMapping {
		return nil, false
	}

	return dereference(document, subschema), true
}

// declaredProperties returns the property names a
// subschema declares directly.
func declaredProperties(subschema map[string]any) []string {
	properties, isMapping := subschema["properties"].(map[string]any)
	if !isMapping {
		return nil
	}

	names := make([]string, 0, len(properties))
	for name := range properties {
		names = append(names, name)
	}

	return names
}

// requiredProperties returns the property names a
// subschema requires.
func requiredProperties(subschema map[string]any) []string {
	required, isSequence := subschema["required"].([]any)
	if !isSequence {
		return nil
	}

	names := make([]string, 0, len(required))
	for _, entry := range required {
		name, isString := entry.(string)
		if !isString {
			continue
		}

		names = append(names, name)
	}

	return names
}

// constrainedConstants returns the properties a subschema
// pins to one literal string, keyed by property name.
// These are what makes a variant recognizable: a `type`,
// `source`, or `kind` field whose value names the branch
// it belongs to.
func constrainedConstants(subschema map[string]any) map[string]string {
	properties, isMapping := subschema["properties"].(map[string]any)
	if !isMapping {
		return nil
	}

	constants := make(map[string]string)
	for name, declaration := range properties {
		property, isMapping := declaration.(map[string]any)
		if !isMapping {
			continue
		}

		value, isString := property["const"].(string)
		if !isString {
			continue
		}

		constants[name] = value
	}

	return constants
}

// propertiesBelow walks the schema keywords from base
// down to target and collects every property name the
// subschemas along the way declare - in other words, the
// properties the failing branch at target would have
// evaluated had it succeeded. The base itself is left
// out: its own `properties` are evaluated before
// `unevaluatedProperties` runs, so they can never end up
// unevaluated.
func propertiesBelow(document any, base string, target string) map[string]struct{} {
	rest, isBelow := strings.CutPrefix(target, base+"/")
	if !isBelow {
		return nil
	}

	current, isPresent := valueAt(document, base)
	if !isPresent {
		return nil
	}

	names := make(map[string]struct{})
	for _, segment := range pointerSegments("/" + rest) {
		_, isDescending := instanceDescendingKeywords[segment]
		if isDescending {
			break
		}

		next, isPresent := step(current, segment)
		if !isPresent {
			break
		}
		current = next

		subschema, isMapping := current.(map[string]any)
		if !isMapping {
			continue
		}

		for _, name := range declaredProperties(dereference(document, subschema)) {
			names[name] = struct{}{}
		}
	}

	return names
}
