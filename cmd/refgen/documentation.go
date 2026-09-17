// Command refgen writes the reference snippets that the documentation
// site embeds: the schema excerpts and the Linter Rules entries. Run it
// whenever the embedded schemas or the rule catalog change; the
// documentation sync test verifies that the committed snippets match
// what refgen would produce.
package main
