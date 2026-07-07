package model

// FacadeOverrides lets a facade declare deviations from the
// default reflection convention used by the schema-facade
// sync test. In 95% of cases a facade method Foo maps to
// the schema field "foo" (first letter lowercased) without
// any override being necessary.
//
// Use Rename when the Go method name cannot match the
// schema field name mechanically - typically because the
// field starts with a reserved character (e.g. "$schema")
// or contains an acronym that would be mangled by Go's
// naming conventions (e.g. APIVersion -> apiVersion).
//
// Use UnusedSchemaFields to explicitly opt out of exposing
// a schema field on the facade, for instance because it
// is purely internal or metadata that the linter does not
// reason about.
type FacadeOverrides struct {
	// Rename maps Go method names to schema field names
	// when they differ from the default convention.
	Rename map[string]string

	// UnusedSchemaFields lists schema fields that are
	// intentionally not exposed as facade methods.
	UnusedSchemaFields []string
}
