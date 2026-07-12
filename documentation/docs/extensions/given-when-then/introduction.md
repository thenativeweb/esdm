# Given-When-Then

The **Given-When-Then** extension captures behavioral specifications about a single consistency unit. A Feature document groups related Scenarios under one shared subject; each Scenario follows the canonical Given/When/Then pattern.

The format originated in behavior-driven development (BDD, Dan North) for general behavioral specs. The event-sourced flavor – Given some Events, When a Command (or Event, timer, or Query) happens, Then a particular outcome – traces back to Greg Young's aggregate-test pattern, and ESDM adopts it across all four kinds of consistency unit.

## When You'd Reach for It

Use the Given-When-Then extension when you want to write down concrete behavior that domain experts can read alongside the model. *Given an Order with two line items, when the customer adds a third item, then a `LineItemAdded` Event is published.*

The Scenarios live alongside the model, in the same `.esdm.yaml` files, version-controlled with the rest of the code. They double as living documentation and as input for tests – an implementation can run the Scenarios as fixtures and verify that the system actually produces the Events the spec calls for.

## Invariants and Scenarios

A prose invariant and a Scenario play different, complementary roles. An **invariant defines the rule** – it lives on the consistency unit in the core model and always holds. A **Scenario verifies that rule by example**: a rejection points at the named invariant it exercises, through `then.rejection.invariant`.

Once a unit is covered by Given-When-Then, every one of its named invariants should be exercised by at least one Scenario – an invariant that no Scenario rejects against is an untested rule. A unit with no Given-When-Then Feature carries no such expectation; Given-When-Then stays optional.

## The Four Variants

A Feature is **about one consistency unit**, and the unit's kind shapes the Scenarios. The four variants are:

- **Aggregate Features** – Given the past Events of this Aggregate, When a Command on this Aggregate, Then the emitted Events or a rejection. The classic Greg-Young aggregate test.
- **Dynamic Consistency Boundary Features** – the writer side mirrors Aggregate Features, but Given carries full Event references because a DCB consults Events from many producers across its Bounded Context.
- **Process Manager Features** – Given carries the Event history that produced the current instance state, When is either an incoming Event or a timer ticking, Then covers the broader reactive surface (emitted Commands, set or canceled timers, the resulting state, an end marker).
- **Read Model Features** – Given is the Event history, When is a Query plus parameters, Then is either the expected query result or the expected materialized Read Model.
