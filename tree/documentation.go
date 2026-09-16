// Package tree builds the containment tree of an ESDM model:
// domain -> subdomains and bounded contexts -> consistency
// units (aggregates, DCBs), free-standing events, read
// models, queries, entities, value objects, domain services
// and actors -> commands, events and features, plus the
// integration layer (process managers, event handlers,
// policies, external systems, context mappings) and the
// extension documents (domain stories under their domain,
// features under the consistency unit they are about).
//
// The tree is the one place where placement is defined.
// Every command that presents the model - the terminal
// summary of `esdm view`, the Markdown tree of `esdm
// documentation` - renders this tree instead of deriving
// placement again, so the two can never disagree about
// where an element sits.
//
// Placement follows one rule: an element sits at the
// position its own `scope` names, and relationships appear
// as annotations on the node, never as placement. An event
// is placed under the aggregate its scope names, or directly
// under the bounded context when its scope names no
// aggregate; the command that publishes it is an annotation,
// not its parent. The one exception is `context-mapping`,
// which has no scope by design because its endpoints may
// straddle domains; it is placed under every domain it
// touches. The rule keeps every document reachable through
// exactly the path its scope spells out, which is what path
// narrowing and the reference notation rely on.
//
// Narrowing walks the tree by name. A name is not unique at
// a position across kinds - an aggregate and a read model
// may share one - so the walk keeps every match at each
// depth and returns several final matches under a synthetic
// root, which renderers show transparently.
package tree
