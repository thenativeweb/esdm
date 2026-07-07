---
paths:
  - "documentation/**"
---

# Writing

Rules for English prose in the documentation under `documentation/`. These build on the formatting rules in `markdown.md`.

## Language and Voice

- Write in English.
- Address the reader directly as "you".
- Contractions are fine (don't, can't, it's); they keep the tone approachable.

## Terminology

- "Event Sourcing": always two words, both capitalized.
- "CQRS": always uppercase.
- "Domain-Driven Design": hyphenated, each word capitalized.
- "Event Storming": always two words, both capitalized.
- SQL keywords: uppercase and wrapped in backticks, for example `SELECT`, `INSERT`, `UPDATE`, `DELETE`, `ALTER TABLE`.
- Event type names: wrapped in backticks; kebab-case in technical or code context (`book-borrowed`), PascalCase in conceptual prose (`BookBorrowed`).

## Prose and Formatting

- Keep paragraphs short, about two to four sentences.
- An H2 heading appears every three to five paragraphs. Mix the heading style: statements, questions, topics.
- Explain a technical term briefly on first use, except very common ones such as Event Sourcing or CQRS.
- Bold the central, load-bearing statements so they stand out.
- Aim for roughly 20 percent of the text in bold.
- Write links in bold: `**[text](url)**`.
- For internal links, use absolute paths without the domain and keep the `.md` extension, for example `/getting-started/installing-esdm.md`. Place them organically in the text, not bundled at the end.
- For an inserted clause, use a spaced en-dash (U+2013), not an em-dash (U+2014). Do not over-apply this: the plain ASCII hyphen (U+002D) stays in compound words, kebab-case identifiers, CLI flags, YAML sequences, and YAML's `---` document separator; never replace those with en-dashes.
- Avoid generic closing headings such as "Summary", "Conclusion", "Where to Go From Here", or "Final Thoughts"; use a content-specific heading instead.
- Include a code example only when it adds value; do not add code that carries no weight.

## ESDM-Specific Conventions

- ESDM kinds (`aggregate`, `event`, `command`, ...) are written in lowercase kebab-case when they refer to the schema kind name, and in the capitalized conceptual form (Aggregate, Event, Command) in prose.
- Avoid the verb *fire* and the noun *firing*. For timers, use **tick / ticks / ticking** and the noun **tick** ("a tick of a named timer"). For Events, use **occur** ("when the Event occurs") or **publish** ("when the Event is published"). For rules and policies, use **run**, **flag**, or **catch**, depending on what's actually meant.
- Rule IDs (`esdm/<category>/<name>`) never appear in the documentation. Diagnostics are a runtime concern; the docs describe the model, not the linter's internals.
- The schema URL convention is `apiVersion: schema.esdm.io/<name>/v1` in every example YAML. Never the legacy `esdm.thenativeweb.io/schema/...` form.
- The ESDM file extension in prose is always `.esdm.yaml`, never the glob form `*.esdm.yaml`. The leading asterisk is shell glob syntax; use it only inside an actual shell command.
- The documentation presents ESDM as a product in its own right. "the native web" appears only where it has a factual role – as the maintainer and rightsholder on the License page, and through the Privacy Policy and Legal Notice links in the footer navigation – not in regular pages.
