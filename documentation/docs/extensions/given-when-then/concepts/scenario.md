# Scenario

A **Scenario** is a single behavioral assertion inside a **[Feature](/extensions/given-when-then/concepts/feature.md)**. It carries three blocks: Given, When, and Then.

**Given** lists the preceding Events that put the system into the state the Scenario starts from. The list is in append order – earlier entries precede later ones in the Event log. The empty list is legal and means "no preceding history" (initial Commands, Process Manager start, fresh Read Model).

**When** is the trigger that drives the Scenario. The trigger's shape depends on the Feature variant. Aggregate and Dynamic Consistency Boundary Scenarios take a Command plus its data (and an optional Actor). Process Manager Scenarios take either an incoming Event or a tick of a named timer. Read Model Scenarios take a Query plus its parameters.

**Then** is the expected outcome. Aggregate and Dynamic Consistency Boundary Scenarios expect either a list of emitted Events or a rejection. Process Manager Scenarios may expect any combination of emitted Commands, set or canceled timers, an updated instance state, or an end marker. Read Model Scenarios expect either the query result or the materialized Read Model content.

## Cross-Reference Integrity

The Commands, Events, timers, Actors, Queries, and Read Models a Scenario mentions all refer back to artifacts the rest of the model declares. A Command in a When block must exist as a Command on the targeted unit; an Event name in Given must exist as an Event somewhere; a timer in When must be a timer the Process Manager declares; an Actor in When must be among the Command's permitted issuers.

The model captures these references as facts that hold across the whole project. Structural validity – does the YAML parse, does the Scenario have the right shape – is one concern; cross-reference integrity is another, and keeping them separate makes both easier to evolve.
