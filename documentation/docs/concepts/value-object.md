# Value Object

A **Value Object** is a typed, named piece of data that has **no identity of its own**. Two Value Objects with the same fields are the same Value Object; you don't compare them by reference, you compare them by value.

In ESDM, Value Objects are the building blocks of Command and Event payloads, and of the state inside an **[Aggregate](/concepts/aggregate.md)**. A `MoneyAmount` is a Value Object. A `PostalAddress` is a Value Object. A `Score` is a Value Object. They have structure, they have rules, but they do not have independent existence.

## Why Model Them Explicitly

You can write a Command or an Event payload directly with primitive fields – a string here, a number there – and most models do that for the simplest cases. The reason to lift a structure into a Value Object is the moment you find yourself writing the same shape twice, or the moment a piece of data gains rules of its own.

A `MoneyAmount` that lives only as `currency: string, amount: number` is fine until you realize you also need to enforce that `amount` is non-negative and that `currency` is a known code. At that point, naming the shape is cheaper than re-stating the rules at every use site. The Value Object is where the rules live.

## What a Value Object Is Not

A Value Object is not a row in a database. It is not a Read Model. It is not an Aggregate. It is a **structural definition** that other parts of the model can refer to by name. The same Value Object can appear inside multiple Commands, multiple Events, and inside Aggregate state, and each appearance carries the same meaning.
