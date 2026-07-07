// Package updateschema provides the `esdm update-schema`
// subcommand. It refreshes the local `schemas/` directory
// to match the embedded schema set: if any embedded
// schema has a higher SemVer revision than the local
// counterpart, or the local inventory drifts from the
// embedded inventory, the directory is wiped and rewritten
// in full. If every embedded revision equals its local
// counterpart and the bytes match exactly, the command
// reports that there is nothing to update.
//
// Local revisions strictly higher than the embedded ones
// are rejected: the binary cannot downgrade a project's
// schemas, and the user is told to use a newer binary
// instead. The wipe-and-rewrite policy is intentional -
// the `schemas/` directory is under the binary's sole
// control by convention, and the linter rejects any
// drift, so partial updates would only paper over
// inconsistencies.
package updateschema
