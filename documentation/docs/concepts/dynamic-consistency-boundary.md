# Dynamic Consistency Boundary

A **Dynamic Consistency Boundary**, or **DCB**, is a consistency unit whose scope is determined at runtime rather than fixed at modeling time. It is the right shape when an invariant has to hold over a set of **[Entities](/concepts/entity.md)** that varies by **[Command](/concepts/command.md)**.

An **[Aggregate](/concepts/aggregate.md)** has a fixed boundary: one Aggregate, one identifier, one transactional scope. A DCB is the variant where the boundary is computed from the Command itself. *Reserve seats 5, 7, and 9 for tonight's screening, atomically* – the seats involved depend on the request, and the consistency that must hold spans only those seats.

## When a DCB Is the Right Shape

DCBs solve a specific kind of problem: an invariant that crosses entities, but only the entities mentioned in a particular Command. If you find yourself wanting an Aggregate that "owns the world" so it can enforce a rule across many things, a DCB is usually the better answer. The Aggregates stay narrow, and the cross-cutting consistency lives in a unit explicitly designed for it.

The trade-off is that a DCB is a more advanced concept, and not every runtime supports it natively. Modeling a DCB is fine even if the implementation has to fall back to a coarser-grained Aggregate – the model captures the intent, and the runtime catches up.

## Identity, State, Consults

A DCB is identified by a field that names the **set** the boundary covers. It declares which Events and Aggregates it **consults** to enforce the invariant, so the cross-cutting consistency story is internally coherent.

Like an Aggregate, a DCB receives Commands and produces Events. The Events are still single-publisher in the rest of the model; the DCB is just a different shape of publisher.
