// Package schema embeds the ESDM core schema and any
// extension schemas at compile time, so the linter binary
// is self-contained and the `add-schema` /
// `update-schema` commands can materialize them into a
// project's local `schemas/` directory without relying on
// an installation path or a network round-trip.
//
// The authoritative source for all ESDM document
// structure lives here: the `lint` command reads from it
// to validate user documents and to verify that any local
// `schemas/` directory matches the embedded set
// byte-for-byte; `add-schema` and `update-schema` read
// from it to seed and refresh that local copy. Every
// consumer shares the same bytes.
package schema
