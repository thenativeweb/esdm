// Package docgen implements the `esdm documentation` command.
// It renders the containment tree of an ESDM model, built by
// the tree package, as a directory of Markdown pages written
// to an output directory: one page per element, containers
// as directories with a README.md index, leaves as files,
// cross-linked with relative links.
//
// The package is named docgen rather than documentation
// because the Go toolchain excludes every non-test file of a
// package literally named documentation from the build.
//
// A page's path follows the element's reference: each
// segment of the reference becomes a directory or file
// named kind_name, with an underscore where the reference
// has an equals sign, and a container ends in README.md. A
// context mapping has no domain and sits at the root. The
// output is neutral Markdown: no site configuration, no
// theme, so GitHub renders it directly and a static-site
// generator can pick it up.
//
// Every page states the element's kind and reference, its
// stats line as `esdm view` shows it, its description, its
// details, and its children grouped by kind. Relationships
// are written out and linked - a command publishes events
// and is issued by actors, an event is published by
// commands - rather than drawn with the view's arrows. A
// bounded context renders its ubiquitous language with the
// translations of each term, a context mapping its endpoints
// and term pairs. A link to an element outside the written
// tree falls back to the element's reference, so no link is
// ever broken.
//
// The optional path argument narrows the output to a subtree
// while keeping full paths, so pages keep the addresses
// their references imply. A non-empty output directory is
// refused unless --force clears it first, so the tree always
// mirrors the model with no orphaned pages.
package docgen
