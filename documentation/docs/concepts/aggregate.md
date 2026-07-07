# Aggregate

An **Aggregate** is the unit of behavior and consistency inside a **[Bounded Context](/concepts/bounded-context.md)**. It owns a slice of the model that must change together, and it is the only place where a transactional rule about that slice can be enforced.

In ESDM, an Aggregate is the canonical kind of *consistency unit*: it accepts **[Commands](/concepts/command.md)**, applies the rules of the domain, and produces **[Events](/concepts/event.md)** that record what happened. Every Event in the model is published by exactly one Aggregate (or by another consistency unit), and every Command targets exactly one.

## Identity, State, and Invariants

An Aggregate is identified by a single field – its natural key. Every Command and Event referencing the Aggregate carries the same identifier, so the model can trace a fact back to the Aggregate that produced it.

Beyond identity, an Aggregate has **state** – the data it carries between Commands – and **invariants** – the rules that must hold over that state. Invariants are the reason Aggregates exist. They are the place where "you cannot withdraw from an empty account" or "a borrowed book cannot be borrowed again before it is returned" lives, and they are checked at the moment a Command is handled.

## When Not to Use an Aggregate

Not every consistency unit is an Aggregate. When the rule that must hold spans multiple Aggregates, a **[Dynamic Consistency Boundary](/concepts/dynamic-consistency-boundary.md)** is the right shape. When the behavior reacts to Events over time without owning state of its own, a **[Process Manager](/concepts/process-manager.md)** is the right shape. When the goal is to expose data for queries rather than to enforce a rule, a **[Read Model](/concepts/read-model.md)** is the right shape.

The Aggregate is the workhorse, but it is not the universal hammer. ESDM has all four shapes because each carries a different invariant, and conflating them turns into the kind of model that is hard to reason about a year later.
