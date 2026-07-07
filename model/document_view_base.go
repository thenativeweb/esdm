package model

import "github.com/thenativeweb/esdm/ast"

// DocumentViewBase exposes the fields every ESDM document
// carries, regardless of kind. It is the common base type
// every kind-specific view (AggregateView, EventView,
// EntityView, ...) embeds; it is never used as a top-level
// document type on its own.
//
// The "Base" suffix marks the role explicitly: this type
// exists to be embedded, and "Document" - not "Entity" -
// is the right noun for "any ESDM YAML document". The
// noun Entity is reserved for the DDD modeling kind of
// the same name.
type DocumentViewBase struct {
	ast.Node
}

// APIVersion returns the apiVersion field.
func (d DocumentViewBase) APIVersion() ast.Node {
	return d.Field("apiVersion")
}

// Kind returns the kind field.
func (d DocumentViewBase) Kind() ast.Node {
	return d.Field("kind")
}

// Name returns the name field.
func (d DocumentViewBase) Name() ast.Node {
	return d.Field("name")
}

// Description returns the description field.
func (d DocumentViewBase) Description() ast.Node {
	return d.Field("description")
}

// Metadata returns the metadata field.
func (d DocumentViewBase) Metadata() ast.Node {
	return d.Field("metadata")
}

// FacadeOverrides declares the rename needed because
// APIVersion cannot spell its acronym in lowercase Go
// method style without conflicting with staticcheck's
// ST1003 rule.
func (DocumentViewBase) FacadeOverrides() FacadeOverrides {
	return FacadeOverrides{
		Rename: map[string]string{
			"APIVersion": "apiVersion",
		},
	}
}
