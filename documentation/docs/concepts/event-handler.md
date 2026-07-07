# Event Handler

An **Event Handler** is a piece of behavior that **reacts to an [Event](/concepts/event.md)** by causing a **side effect outside the model** – sending an email, calling an external API, writing to a log, pushing a notification.

An Event Handler is not an **[Aggregate](/concepts/aggregate.md)** (it owns no state and produces no Events of its own), and it is not a **[Policy](/concepts/policy.md)** (it does not derive Commands from Events). It is the *integration boundary* of the model: the place where the inside of the system meets the outside.

## Event Handler vs. Policy

The boundary between an Event Handler and a Policy is the kind of effect they cause.

A Policy stays inside the model. It reacts to an Event by emitting a Command, which lands on an Aggregate, which produces another Event. The whole chain is observable in the model.

An Event Handler exits the model. It reacts to an Event by sending an email or calling a service, and the outcome of that effect is not part of the Event sequence. If you need to react to the *result* of an external call (the email bounced, the service returned an error), you typically wrap the integration in an **[External System](/concepts/external-system.md)** and let it produce its own Events back into the model.

## Idempotency Matters

Because an Event Handler causes side effects in the world, the question of *what happens if the same Event is delivered twice?* matters more for handlers than for any other kind. An Event Handler that sends the same email twice is rarely what you want. An Event Handler declares its **delivery guarantee** – at-least-once or at-most-once – and how duplicates are handled, so the model captures the assumption rather than leaving it implicit in someone's deployment notes.

## Why Make Handlers Visible

Modeling integrations as Event Handlers makes them visible. When a handler refers to an Event that doesn't exist, the model says so. When an Event has no consumer, the model says so. Integrations are exactly the corner of the system where things rot first; making them visible is the point.
