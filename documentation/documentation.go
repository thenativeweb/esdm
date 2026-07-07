// Package documentation holds the published ESDM documentation
// site together with the tests that keep it honest as the
// language evolves.
//
// Most of the site is prose, but two parts of it mirror the
// embedded schemas and can silently fall out of step with them:
// the reference pages, which describe the language one kind at a
// time, and the schema excerpts those pages embed as generated
// snippets. When a schema gains a field, an enum value, or a
// constant, either surface can be left behind without anyone
// noticing.
//
// The tests in this package are the safety net against that
// drift. They read the embedded schemas as the single source of
// truth and check, first, that every field, enum value, and
// constant a schema defines is actually mentioned on the matching
// reference page, and second, that the embedded snippets are
// exactly what regenerating them from the current schemas would
// produce. A schema change that is not reflected in the
// documentation therefore surfaces as a failing test rather than
// as a stale page.
package documentation
