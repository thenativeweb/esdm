# Extensions

ESDM's core vocabulary is documented under **[Concepts](/concepts/overview.md)** and **[Reference](/reference/core-schema/overview.md)**. Extensions add additional vocabulary on top: independent schemas that describe artifacts the core deliberately leaves out, primarily discovery outputs and behavioral specifications.

Each extension is **self-contained**. It bundles its own concepts, its own reference, and its own examples under the same chapter, so you can explore one extension without first crossing into the rest of the model. The asymmetry between Core (which dominates the top-level navigation) and Extensions (gathered into a single chapter) is deliberate – Core is the language; Extensions are modules.

## Available Extensions

ESDM ships two extensions today:

- **[Domain Storytelling](/extensions/domain-storytelling/introduction.md)** captures discovery stories – Actors, Work Objects, and the activities that connect them – in the format defined by Stefan Hofer and Henning Schwentner.
- **[Given-When-Then](/extensions/given-when-then/introduction.md)** captures behavioral scenarios about a single consistency unit (Aggregate, Dynamic Consistency Boundary, Process Manager, or Read Model) in the canonical Given/When/Then format. Each Scenario lists preceding Events, a triggering Command, Event, timer, or Query, and the expected outcome.

## How an Extension Plugs In

Each extension has its own identifier. A document validates against exactly one schema – the core schema for core documents, the matching extension schema for extension documents. Extensions never inject kinds into the core, and core documents never validate against an extension. That symmetry is what lets us add new extensions over time without ever touching the core.

Extensions are developed and maintained by the ESDM team. They are not user-extensible: writing a new extension requires changes to the binary itself.
