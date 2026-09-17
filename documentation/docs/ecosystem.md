# Ecosystem

ESDM's schema is public, and the `esdm` binary is not the only program that reads it. This page lists tools that other people have built on ESDM, so you can find them without searching, and so their authors have a place to be found. It is a directory, not a list of recommendations: we ran each tool against a current model when we added it, but we do not maintain these tools, and we cannot vouch for how they develop.

## How This Page Works

### What Gets Listed

A tool qualifies for this page when it

- is open source under a recognized license,
- reads ESDM models as they are, following the current schema, rather than requiring a format of its own,
- is publicly usable without registration or payment, and
- is actively maintained.

Each entry names the maintainer, the license, the parts of ESDM the tool reads – the core schema and the extensions –, how to run it, where the source lives, and the ESDM version we last verified it with. That last line is a promise to you and a reminder to us: when a new ESDM version ships, we re-check the entries and update the version, or note what no longer works.

### Proposing a Tool

Have you built something on ESDM? Open a pull request that adds an entry in the shape below, or **[write to us](mailto:hello@thenativeweb.io)** and we'll take it from there. The **[contributing guide](https://github.com/thenativeweb/esdm/blob/main/CONTRIBUTING.md)** explains how pull requests to the documentation work. A typo or a broken link in an existing entry needs no issue – fix it right away.

## Available Tools

### ESDM Visualizer

The ESDM Visualizer renders a directory of `.esdm.yaml` files as an interactive board. Aggregates, Commands, Events, Read Models, Actors, and Bounded Contexts appear as connected nodes, and the arrows come from the model itself – a Command that names its Actors, a Command that publishes an Event, a Read Model that projects one. Domain Stories and Given-When-Then Features get boards of their own. It is useful for walking a team through a model without reading YAML, and for spotting what nothing connects to. The story of how it came about is told in **[Someone built a visualizer for ESDM, and it wasn't us](https://docs.eventsourcingdb.io/blog/2026/08/06/someone-built-a-visualizer-for-esdm--and-it-wasnt-us/)**.

- **Maintainer:** Impierce Technologies
- **License:** Apache-2.0
- **Reads:** core, domain-storytelling, given-when-then
- **Run:** `docker run -p 3000:3000 -v .:/data impierce/esdm-visualizer` in the model directory, then open `http://localhost:3000`
- **Source:** **[github.com/impierce/esdm-visualizer](https://github.com/impierce/esdm-visualizer)**
- **Verified with ESDM:** 0.15.0
