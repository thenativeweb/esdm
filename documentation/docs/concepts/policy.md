# Policy

A **Policy** is a piece of **decision logic** that reacts to **[Events](/concepts/event.md)** and produces **[Commands](/concepts/command.md)**. It encodes the *whenever this happens, do that* rules that hold across the system without belonging to any single **[Aggregate](/concepts/aggregate.md)**.

A Policy is stateless. It does not own data, it does not have an identity, and it is not a consistency unit. It listens to Events, applies a rule, and emits Commands. Whatever happens after that is the business of the Aggregates the Commands target.

## Policy vs. Process Manager

Policies and **[Process Managers](/concepts/process-manager.md)** are easy to confuse, because both react to Events and emit Commands. The distinction is **state**.

A Policy has none. It evaluates a rule against a single Event in isolation: *if this Event has this shape, then emit this Command*. There is no memory of previous Events, no waiting for a second trigger, no timeout.

A Process Manager has state. It tracks where it is in a longer-running flow, waits for multiple Events, can time out, and produces Commands as the flow progresses. If the rule you're modeling fits in a single sentence about a single Event, it's a Policy. If it needs to remember anything between Events, it's a Process Manager.

## When to Reach for a Policy

Use a Policy when the rule is global and stateless. *Whenever an Order is placed, reserve the inventory.* *Whenever a returned Book is overdue, charge the late fee.* *Whenever an Order ships, debit the customer's account.*

Use it sparingly. A Policy that grows to several conditions, that needs context from earlier Events, or that needs to defer its decision is a Policy that has outgrown its shape – move it to a Process Manager before it becomes hard to reason about.
