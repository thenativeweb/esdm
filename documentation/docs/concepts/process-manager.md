# Process Manager

A **Process Manager** is a long-running coordinator. It listens to **[Events](/concepts/event.md)**, holds **state** between them, and emits **[Commands](/concepts/command.md)** as a flow progresses through several steps.

Where a **[Policy](/concepts/policy.md)** is a stateless reaction to a single Event, a Process Manager is what you reach for when the flow needs memory: *wait for both Events to arrive, then emit the next Command*. *Start a timer when this Event happens, and emit a different Command if the timer ticks before another Event arrives*.

## State, Correlation, Lifecycle

A Process Manager has three things every model needs to declare. It has **state** – the data it tracks across the flow. It has a **correlation field** – the identifier that ties incoming Events to the right Process Manager instance. And it has a **lifecycle**: the conditions under which a new instance starts, and the conditions under which an existing instance ends.

All three are explicit, not implicit. Without a starting condition, a Process Manager has no way to come into being. Without an ending condition, it would never finish. Without a correlation field, Events from different flows would land on the same instance.

## Timers

A Process Manager can also schedule **timers**. A timer is a future moment at which the Process Manager wants to react – *if no confirmation arrives within 24 hours, send a reminder*. A timer is a first-class part of the Process Manager's lifecycle: declared with what triggers it, the duration after which it ticks, and the Commands it emits.

## Where Process Managers Belong

A Process Manager lives inside a **[Bounded Context](/concepts/bounded-context.md)**, just like an **[Aggregate](/concepts/aggregate.md)** does. It is a consistency unit in its own right, although the consistency it enforces is **temporal** rather than transactional: it ensures the right Commands happen in the right order, not that two pieces of state change atomically.
