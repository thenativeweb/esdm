# Command

A **Command** is the expression of an intent to change the model. It is named in the **imperative** – `BorrowBook`, `CapturePayment`, `ReserveSeats` – and it targets exactly one consistency unit.

Where an **[Event](/concepts/event.md)** records what happened, a Command requests what should happen. The Command can fail: it can be rejected by an invariant, refused by a permission check, or ignored as a duplicate. A rejection comes from an **invariant** – a rule about the target's state – while a refusal comes from **authorization**, expressed through the Command's permitted issuers (`command.actors`), not as an invariant. Only when it succeeds does an Event come into being.

## Single Target

A Command has exactly one target. Most Commands target an **[Aggregate](/concepts/aggregate.md)**; some target a **[Dynamic Consistency Boundary](/concepts/dynamic-consistency-boundary.md)**. The rule that a Command targets a single unit is what keeps the consistency story tractable. A Command that wants to change two Aggregates is not really one Command; it is two, and the model is more honest if you write them that way.

## Who Issues a Command

A Command is *issued* by an **[Actor](/concepts/actor.md)**, an **[External System](/concepts/external-system.md)**, a **[Policy](/concepts/policy.md)**, or a **[Process Manager](/concepts/process-manager.md)**. The model tracks the issuers explicitly: every Command has at least one source.

A Command without an issuer is a Command nobody can send – almost always a modeling mistake.

## Data and Identity

A Command carries a **payload** – the data the consistency unit needs to honor the request – and an **identifier** for the target unit. The payload is whatever the domain requires; the identifier is the same natural key the target uses for its own identity, so the link between Command and target is unambiguous.

Commands do not carry results. The result of a Command is an Event (success) or a rejection (failure); neither is encoded in the Command itself.
