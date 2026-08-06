// Package docgen renders a resolved ESDM model as a Markdown
// directory tree. Each model element becomes one page, placed at its
// containment path so the tree mirrors the reference notation: an
// element addressed as esdm:<domain>/<...>/<name> lives at the same
// path on disk. An element that contains other elements is written as
// a directory with a README.md index page; a leaf element is written
// as <name>.md. Supplying a base URL therefore turns any reference
// into a link into the rendered tree.
//
// The command narrows to a subtree via an optional model path and
// refuses to write into a non-empty directory unless --force clears it
// first, so the output always mirrors the model exactly with no
// orphaned pages.
package docgen
