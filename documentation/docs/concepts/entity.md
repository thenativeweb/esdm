# Entity

An **Entity** is an identity-bearing modeling element inside a **[Bounded Context](/concepts/bounded-context.md)** that has **no consistency container of its own**. It names the *what* behind an identifier – `student`, `course`, `book`, `seat` – so the model can talk about the thing that carries the ID, not just the ID itself.

An Entity is not a consistency unit. It does not receive Commands, it does not own state in the sense an **[Aggregate](/concepts/aggregate.md)** does, and it does not enforce invariants across a boundary. It is referenced from elsewhere – typically the `identifiedBy` field of a **[Dynamic Consistency Boundary](/concepts/dynamic-consistency-boundary.md)**, or the payload of a Command – and that is the work the Entity exists to do.

## When to Reach for an Entity

Reach for an Entity when the model has a thing that needs a name and an identity, but no Aggregate to hold it. The typical case is a model built around a **[Dynamic Consistency Boundary](/concepts/dynamic-consistency-boundary.md)**: a DCB spans a set of identifiers per decision, but it does not own them. Each identifier still refers to a thing in the model – a Student, a Course – and that thing needs a name in the vocabulary. The Entity is that name.

When the same thing has its own lifecycle and invariants over evolving state, model it as an Aggregate instead. An Aggregate already carries identity together with that state, and a separate Entity for the identity alone would only duplicate what the Aggregate provides.

## What an Entity Is Not

An Entity is not a **[Value Object](/concepts/value-object.md)**: two Entities with the same data are still two different Entities, because identity matters more than data. An Entity is not an **[Aggregate](/concepts/aggregate.md)**: it has no `state`, accepts no Commands, and is not a unit of consistency. An Entity is not an **[Actor](/concepts/actor.md)**: Actors initiate Commands, Entities are the things Commands and DCBs refer to. An Entity is not a place to put behavior – computations that span several Entities belong on a **[Domain Service](/concepts/domain-service.md)** or on the Aggregate or DCB that hosts the decision.

## Anatomy

An Entity declares a `schema` describing the shape of a single instance, in the same paradigm-agnostic JSON Schema form a Value Object uses. It declares an `identifiedBy` strategy that says where the identity comes from – either a field inside the schema (`source: schema`) or a fixed string when the Entity exists exactly once (`source: static`). It can declare `invariants` over the instance shape – checksums, cross-field constraints, format rules – the same way a Value Object can, because both lack a lifecycle. There is no `state`, no `consults`, no `publishes`: an Entity is a description, not a mechanism.

## Where Identity Comes From

An Entity describes *what* its identity looks like; it does not mint identifiers. When a new Student appears in the model, the identifier arrives through a Command – issued by an Actor, accepted by an Aggregate or DCB – not by the Entity itself. The Entity is the modeling anchor that lets the Command and the DCB say "this Student" and have everyone in the model know what *Student* refers to.
