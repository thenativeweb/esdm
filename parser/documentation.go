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
// Multiple documents per file are separated by YAML's
// canonical `---` (three ASCII hyphens, U+002D). The
// repo's en-dash convention for prose (U+2013) does
// not apply here: the separator is part of the YAML
// grammar and must remain three plain hyphens.
package parser
