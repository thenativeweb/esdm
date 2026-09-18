// Package parser turns ESDM source files into a parsed,
// position-aware AST and reports any structural
// diagnostics (YAML syntax errors, unknown apiVersions,
// schema violations).
//
// Parsing is a two-step process per file: first the
// contents are decoded into a YAML node tree so every
// downstream consumer has line and column information;
// then each document within the file is validated
// against the compiled schema its apiVersion refers to -
// the ESDM core schema or one of the embedded extension
// schemas. A document whose apiVersion matches no
// compiled schema is reported rather than validated.
// Schema validation errors are translated into
// structure/* diagnostics - missing-required-field,
// type-mismatch, unknown-field, or constraint-violation -
// depending on which JSON Schema keyword failed.
//
// One defect yields one diagnostic. A JSON Schema
// validator reports more than that, because a schema
// combines branches and a failed branch takes consequences
// with it. Two of those consequences are filtered out
// here, and both need the schema itself to decide, not
// just the error: the decoded schema document is
// therefore kept next to every compiled validator and
// addressed by the JSON pointer inside each error's
// schema URL.
//
// The first consequence comes from closing a document
// with unevaluatedProperties. A failed branch drops its
// annotations, so every field that branch would have
// evaluated counts as unevaluated and reads as an unknown
// field. Such a field is reported only when no failing
// branch declares it, which leaves genuine typos intact
// and drops the fallout of the real defect.
//
// The second comes from oneOf, where the validator
// reports the causes of every alternative. The branch the
// document aims at is identified by how much of it the
// document already exhibits: required properties that are
// present, and constants it matches. A branch that is
// itself a set of alternatives states its requirements
// one level down, so it is scored by the alternative the
// document comes closest to; otherwise it would score
// nothing and lose to any flat branch beside it. Only the
// winning branch is reported. When the document resembles
// the branches equally, all of them are reported, because
// then the document truly identifies none - and since
// several of them can then fail on the very same field, a
// diagnostic repeating one already reported is dropped.
//
// Filtering must never turn into silence. Some schema
// keywords fail without nested causes, so a translation
// that only walked causes would say nothing at all about
// a document the schema rejected. The two such keywords
// the schemas use - forbidding a field through not, and
// matching more than one alternative of a oneOf - are
// phrased in terms of the model, and any other failure
// that produced nothing reports itself. A linter that
// stays quiet about a broken document is worse than one
// that says too much; that direction is what the mutation
// sweep in the tests guards.
//
// Multiple documents per file are separated by YAML's
// canonical `---` (three ASCII hyphens, U+002D). The
// repo's en-dash convention for prose (U+2013) does
// not apply here: the separator is part of the YAML
// grammar and must remain three plain hyphens.
package parser
