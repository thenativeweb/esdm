# Linter Rules

The Given-When-Then extension brings its own rules. They run with every `esdm lint` and throw only where a model contains Features. They check that a Feature and its Scenarios agree with the model they describe: that every Command, Event, Query, Timer, Actor, and Invariant a Scenario names is declared where the Feature's scope says it is, and that the shape of each Scenario matches the variant of its Feature.

All of these rules carry the `gwt` category. Severities and the format of the IDs are explained on the core **[Linter Rules](/reference/linter-rules.md)** page.

--8<-- "reference/linter-rules/given-when-then.md"
