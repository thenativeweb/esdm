// Package addschema provides the `esdm add-schema`
// subcommand. It writes the embedded core schema and all
// embedded extension schemas into a `schemas/` directory
// at the current working directory's root, intended for
// editor support (YAML language servers resolving
// `# yaml-language-server: $schema=...` references against
// the local copy).
//
// The command refuses to run if `schemas/` already exists,
// pointing the user at `update-schema` instead. This
// preserves the invariant that the local `schemas/`
// directory is always under the binary's sole control: a
// fresh install (`add-schema`) is distinct from a refresh
// (`update-schema`), and the user has to declare which one
// they intend.
package addschema
