# Event

An **Event** is a fact that something happened. It is named in the **past tense** – `BookBorrowed`, `PaymentCaptured`, `OrderShipped` – and it carries the data that makes the fact concrete: who, what, when.

Events are **immutable**. Once published, an Event never changes. It records the state of the world at the moment it was emitted, and any later change is itself an Event. That is the foundation of Event Sourcing: state is not the truth; the sequence of Events is.

## Publisher and Consumers

Every Event is **published by exactly one consistency unit**. Most often that publisher is an **[Aggregate](/concepts/aggregate.md)**, but it can also be a **[Process Manager](/concepts/process-manager.md)**, a **[Dynamic Consistency Boundary](/concepts/dynamic-consistency-boundary.md)**, or another publishing kind. The single-publisher rule is what keeps responsibility traceable: when you ask *who emits this Event?*, the model gives one answer.

An Event may have any number of **consumers** – **[Read Models](/concepts/read-model.md)** that project it, **[Process Managers](/concepts/process-manager.md)** that react to it, **[Policies](/concepts/policy.md)** that derive Commands from it. An Event without consumers is almost always a modeling mistake: nobody is listening, and nothing downstream changes when it occurs.

## Data and Identity

An Event carries a **payload** – the structured data that describes the fact. It also carries an **identifier** – the natural key of the consistency unit it belongs to – so that downstream consumers know which entity the fact is about.

Events do **not** carry intent, decisions, or future tense. `OrderShouldShip` is not an Event; it is a Command that has not happened yet. `OrderShipped` is the Event that records the shipment as a fact. Keeping the tense rigorous is one of the simplest disciplines you can apply, and it pays off every time you read the model later.
