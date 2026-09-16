# Event

An **Event** is a fact that something happened. It is named in the **past tense** – `BookBorrowed`, `PaymentCaptured`, `OrderShipped` – and it carries the data that makes the fact concrete: who, what, when. These are the spoken forms; the **[Naming](#naming)** section below shows how they map to the model.

Events are **immutable**. Once published, an Event never changes. It records the state of the world at the moment it was emitted, and any later change is itself an Event. That is the foundation of Event Sourcing: state is not the truth; the sequence of Events is.

## Publisher and Consumers

Every Event is **published by exactly one consistency unit**. Most often that publisher is an **[Aggregate](/concepts/aggregate.md)**, but it can also be a **[Process Manager](/concepts/process-manager.md)**, a **[Dynamic Consistency Boundary](/concepts/dynamic-consistency-boundary.md)**, or another publishing kind. The single-publisher rule is what keeps responsibility traceable: when you ask *who emits this Event?*, the model gives one answer.

An Event may have any number of **consumers** – **[Read Models](/concepts/read-model.md)** that project it, **[Process Managers](/concepts/process-manager.md)** that react to it, **[Policies](/concepts/policy.md)** that derive Commands from it. An Event without consumers is almost always a modeling mistake: nobody is listening, and nothing downstream changes when it occurs.

## Data and Identity

An Event carries a **payload** – the structured data that describes the fact. It also carries an **identifier** – the natural key of the consistency unit it belongs to – so that downstream consumers know which entity the fact is about.

Events do **not** carry intent, decisions, or future tense. `OrderShouldShip` is not an Event; it is a Command that has not happened yet. `OrderShipped` is the Event that records the shipment as a fact. Keeping the tense rigorous is one of the simplest disciplines you can apply, and it pays off every time you read the model later.

## Naming

In the model, an Event carries the **bare name**: the past-tense verb, in kebab-case, without the Aggregate it belongs to. The Aggregate `book` publishes `reserved`, not `book-reserved`. The Event's scope already names the Aggregate, and every other form of the name is derived from that pair.

Spoken, the two combine: "book reserved". That is the form these pages use when they write `BookReserved`. Wherever an Event leaves the model, the Aggregate is composed in at the boundary that needs it. As a CloudEvents type it becomes `io.eventsourcingdb.library.book-reserved`; in code it is `Reserved` when the language nests it inside the Aggregate, and `BookReserved` when it stands as a class of its own. A generator produces each of these forms from scope plus name, deterministically and the same way every time.

That is why the name must not repeat the Aggregate. `book-reserved` under `book` derives to `book-book-reserved`, and the linter warns about it with `esdm/modeling/event-name-with-aggregate-prefix`. The same convention holds for **[Commands](/concepts/command.md)**: `reserve` under `book`, spoken "reserve book", rendered as `ReserveBook` – never `reserve-book` in the model.

A **free-standing Event** has no Aggregate to compose in, so its name has to stand on its own: choose one that says what happened without a container to lean on. The linter rule above does not apply to it.
