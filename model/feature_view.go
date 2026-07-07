package model

import "github.com/thenativeweb/esdm/ast"

// FeatureView is the typed view over an ESDM document
// whose kind is "feature", defined by the given-when-then
// extension schema.
type FeatureView struct {
	DocumentViewBase
}

// Scope returns the scope field. The scope is one of four
// variants distinguished by which discriminator field is
// present - aggregate, dynamicConsistencyBoundary,
// processManager, or readModel.
func (f FeatureView) Scope() ast.Node {
	return f.Field("scope")
}

// Scenarios returns the scenarios array.
func (f FeatureView) Scenarios() ast.Node {
	return f.Field("scenarios")
}

// ScenarioView is a typed view over a single entry inside
// a feature's scenarios array. Unlike DocumentViewBase-based
// views, scenarios are not standalone ESDM documents - a
// scenario only exists inside its enclosing feature - so
// ScenarioView wraps an ast.Node directly without the
// document-level apiVersion/kind/name infrastructure.
type ScenarioView struct {
	ast.Node
}

// Name returns the scenario's bare name.
func (s ScenarioView) Name() ast.Node {
	return s.Field("name")
}

// Description returns the scenario's description prose.
func (s ScenarioView) Description() ast.Node {
	return s.Field("description")
}

// Given returns the list of preceding events.
func (s ScenarioView) Given() ast.Node {
	return s.Field("given")
}

// When returns the trigger - the form depends on the
// enclosing feature's scope variant.
func (s ScenarioView) When() ast.Node {
	return s.Field("when")
}

// Then returns the expected outcome - the form depends on
// the enclosing feature's scope variant.
func (s ScenarioView) Then() ast.Node {
	return s.Field("then")
}
