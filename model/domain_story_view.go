package model

import "github.com/thenativeweb/esdm/ast"

// DomainStoryView is the typed view over an ESDM
// document whose kind is "domain-story", defined by the
// domain-storytelling extension schema.
type DomainStoryView struct {
	DocumentViewBase
}

// Scope returns the scope field (scopeDomain -
// {domain}).
func (d DomainStoryView) Scope() ast.Node {
	return d.Field("scope")
}

// PointInTime returns the pointInTime field (as-is or
// to-be).
func (d DomainStoryView) PointInTime() ast.Node {
	return d.Field("pointInTime")
}

// Granularity returns the granularity field
// (coarse-grained or fine-grained).
func (d DomainStoryView) Granularity() ast.Node {
	return d.Field("granularity")
}

// DomainPurity returns the domainPurity field (pure or
// digitalized).
func (d DomainStoryView) DomainPurity() ast.Node {
	return d.Field("domainPurity")
}

// Groups returns the groups registry.
func (d DomainStoryView) Groups() ast.Node {
	return d.Field("groups")
}

// Actors returns the story-global actors list.
func (d DomainStoryView) Actors() ast.Node {
	return d.Field("actors")
}

// Sentences returns the ordered sentences that make up
// the story.
func (d DomainStoryView) Sentences() ast.Node {
	return d.Field("sentences")
}
