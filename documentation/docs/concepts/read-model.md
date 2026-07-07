# Read Model

A **Read Model** is a view of the domain that is optimized for **reading** rather than for enforcing rules. It is the answer to *what does the world look like right now?* – derived from the Events that have been published, shaped to fit the queries it serves.

A Read Model is **not** an **[Aggregate](/concepts/aggregate.md)**. It does not enforce invariants, it does not accept **[Commands](/concepts/command.md)**, and it does not own behavior. It owns **shape and projection**: which Events to listen to, how to fold them, what the resulting structure looks like.

## Projections and Queries

A Read Model declares two things. Its **projections** describe the Events it consumes and how those Events update its data. Its **queries** describe the questions the Read Model can answer for a caller.

A Read Model with no projections has no source of truth and would never have data. A Read Model with no queries is a write-only sink – data goes in, but nothing reads it. Both ends matter: the projections keep the Read Model honest, the queries make it useful.

## Many Read Models per Aggregate

A single Aggregate can feed any number of Read Models. The customer-list Read Model and the customer-detail Read Model and the customer-by-region Read Model can all consume the same `CustomerRegistered` and `CustomerMoved` Events, each shaping them differently for its purpose.

This is one of the ergonomic wins of Event Sourcing. The write side stays focused on invariants, the read side stays focused on queries, and they communicate through the immutable record of Events.

## Paradigm

A Read Model can be labeled with its **paradigm** – tabular, document, graph, key-value, time-series, vector, and so on. The label is descriptive rather than enforced, but it is a fact about the Read Model that helps later readers understand what kind of database it expects to live in.
