package model

import "regexp"

// LanguageTagPattern is the shape of a language tag as the
// core schema's `languageTag` definition accepts it: a two-
// or three-letter primary subtag, optionally followed by
// further subtags. The schema is the authority; a test keeps
// this copy identical to it, so that the CLI can validate a
// language flag the same way the schema validates a
// document.
var LanguageTagPattern = regexp.MustCompile(`^[a-z]{2,3}(-[A-Za-z0-9]{2,8})*$`)

// IsLanguageTag reports whether value has the shape of a
// BCP 47 language tag.
func IsLanguageTag(value string) bool {
	return LanguageTagPattern.MatchString(value)
}
