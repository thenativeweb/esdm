// Package model provides the typed facades that the rest
// of the esdm linter uses to navigate a parsed ESDM
// document. A facade is a thin, 1:1 mirror of a schema
// kind: each method corresponds to exactly one schema
// field and returns a position-aware syntax node from the
// ast package.
//
// The package does not reify a materialized domain model.
// Views are lightweight wrappers around those syntax nodes,
// and the model index resolves entity names to their views.
// Rules query the index; the index never copies or
// transforms data.
package model
