# Linter Rules

The Domain Storytelling extension brings its own rules. They run with every `esdm lint` and throw only where a model contains Domain Stories. They check that a story is consistent in itself: that every declared Actor and Work Object is drawn by at least one edge, that every Group has members and every membership names a declared Group, that no name is declared twice where the story needs it to be unique, and that a story has Sentences and every Sentence has edges.

These rules use the `structure` and `modeling` categories like the core rules; their names start with `story-`. Severities and the format of the IDs are explained on the core **[Linter Rules](/reference/linter-rules.md)** page.

--8<-- "reference/linter-rules/domain-storytelling.md"
