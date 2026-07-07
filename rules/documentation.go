// Package rules defines the interface every lint rule
// implements, the metadata carried by each rule, and the
// catalog of all built-in rules.
//
// Each rule carries a stable identifier and is created
// through a uniform constructor, even when it has no
// dependencies today. That uniform shape lets a rule gain
// dependencies later without changing how the catalog is
// assembled.
package rules
