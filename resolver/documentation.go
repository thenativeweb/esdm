// Package resolver turns a set of parsed ESDM files into
// a model index: it groups documents by their kind, keys
// them by a scope-aware composite name, and verifies the
// cross-references between them.
//
// Indexing comes first. Documents that share a bare name
// but live in different scopes - for example two commands
// belonging to different aggregates - are distinct and
// coexist in the index; only a genuine clash within the
// same kind and scope counts as a duplicate. When names
// do clash, the first occurrence is the one that ends up
// in the index and the later ones are reported, so the
// rest of the pipeline always sees a single, deterministic
// definition per name. Documents whose apiVersion or kind
// the resolver does not recognize are skipped silently,
// because the parser has already reported them.
//
// Reference resolution comes second. Every scoped entity
// is walked again and each reference it makes - its scope
// chain, its bare-name references to sibling entities, its
// structured references to events and commands, its
// mapping endpoints, and its references to fields of its
// own state or of a referenced payload - is checked
// against the index. Bare-name references carry no scope
// of their own; they piggyback on the scope of the entity
// that makes them, which is what lets a short, local name
// resolve unambiguously.
//
// Failures come in two shapes that are kept deliberately
// apart. A reference to something that does not exist at
// all is unresolved; a reference to something that does
// exist but sits under a different parent than the
// reference claims is a parent mismatch, and its
// diagnostic points at the real definition so the two can
// be compared side by side. When an unresolved name is
// close to an existing one, the diagnostic carries a
// did-you-mean hint computed from the edit distance to the
// nearest candidate.
//
// The resolver is purely a derivation step: it produces a
// model that later phases - first and foremost the rule
// engine - can query without ever going back to the
// underlying YAML.
package resolver
