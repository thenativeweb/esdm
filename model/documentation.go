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
//
// The index also answers the few relationship questions
// that more than one consumer asks - first of all which
// commands publish an event. Such a question has exactly
// one definition here, and the resolver, the rules, and the
// view all use it, instead of each rebuilding the relation
// from the raw scopes and drifting apart.
package model
