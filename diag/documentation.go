// Package diag provides the diagnostic vocabulary of the esdm
// linter: what a finding is, how severe it is, and where in a
// source document it points.
//
// Every stage of the linting pipeline - parsing, resolving, and
// the lint rules - produces findings in one uniform shape, and
// only the reporting side consumes them. This single shared
// vocabulary decouples the producers from the output formatting:
// a stage never needs to know whether its findings end up as
// human-readable terminal output or as structured data.
//
// A finding carries a severity that is fixed per rule, never per
// occurrence; only error-level findings make the linter exit with
// a non-zero status. Besides its primary source position, a
// finding can reference secondary positions, so that, for
// example, a duplicate name can point at both the offending and
// the original definition.
//
// Source positions use the zero value to express "no meaningful
// source location" instead of an optional: system-level findings
// that do not originate from a user file simply leave the
// position empty, and consumers ask the position itself whether
// it carries source information.
//
// Because rules run concurrently, anything that collects findings
// must be safe for concurrent use.
package diag
