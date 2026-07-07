# Domain Service

A **Domain Service** is a piece of **stateless domain logic** that doesn't belong to any single **[Aggregate](/concepts/aggregate.md)** because it operates over several – or because it has no clear single owner.

When a calculation, a check, or a transformation belongs to the domain (the rules of the business apply) but no single Aggregate is the natural home for it, you reach for a Domain Service. *Compute the optimal route from A to B given the current traffic.* *Decide whether two transactions are likely the same fraud pattern.* *Check whether an account is permitted to perform this Command.*

## What a Domain Service Is Not

A Domain Service is not infrastructure. *Send an email* is not a Domain Service; it is an integration concern, modeled as an **[Event Handler](/concepts/event-handler.md)** or via an **[External System](/concepts/external-system.md)**.

A Domain Service is also not a stateful coordinator. *Wait for these three Events, then emit a Command* is a **[Process Manager](/concepts/process-manager.md)**, not a Domain Service. The Domain Service is **stateless and synchronous**: in, compute, out.

## Why Model One

Modeling a Domain Service explicitly is a way of capturing that *this piece of logic exists, and it isn't owned by any Aggregate*. Without that statement, the logic tends to scatter – a copy here, a slightly different copy there – and the model loses fidelity.

A Domain Service declares the **functions** it offers (the operations it performs). A Domain Service with no functions is a Domain Service that does nothing – the kind of thing that drifts in.
