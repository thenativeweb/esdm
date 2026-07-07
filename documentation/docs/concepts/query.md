# Query

A **Query** is a question a caller asks of a **[Read Model](/concepts/read-model.md)**. It declares the parameters the caller provides, the **result** the Read Model returns, and which Read Model serves it.

A Query is purely a read operation. It does not change state, it does not produce **[Events](/concepts/event.md)**, and it does not have side effects. *Get the open invoices for customer X*. *List the books available in branch B*. *Show the leaderboard for the current month*.

## Why Queries Are Modeled

You could argue that Queries are an implementation detail – the Read Model is there, just read it. Queries are modeled explicitly because the **shape of a Query is part of the contract** between consumers and the Read Model. When a Query is written down, the Read Model knows what it has to support, and the caller knows what it can expect.

In practice, this matters most when several callers consume the same Read Model. The Query is the public face; the projections are the implementation. The Read Model can change how it builds its data without changing the Queries, and the callers don't notice.

## Query and Read Model

Every Query references the Read Model it serves. The link goes both ways: every Query points to a Read Model that exists, and every Read Model that has Queries also has projections to feed them. A Query without a Read Model is a question without a source; a Read Model without Queries is a projection without consumers.
