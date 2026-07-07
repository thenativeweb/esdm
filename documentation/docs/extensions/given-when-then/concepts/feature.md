# Feature

A **Feature** is a group of Scenarios about the same consistency unit. Picking one consistency unit per Feature is the structural rule that keeps Scenarios focused: when you find yourself wanting to mix Scenarios about an Aggregate and a Process Manager, you have two Features, not one.

A Feature has four variants, one per kind of consistency unit it can target: Aggregate, Dynamic Consistency Boundary, Process Manager, or Read Model. The variant determines the shape of every **[Scenario](/extensions/given-when-then/concepts/scenario.md)** inside the Feature.

## Why Events Are Referenced Differently per Variant

Aggregate Features identify Events by **bare name** because the Feature already fixes the producing Aggregate. Dynamic Consistency Boundary, Process Manager, and Read Model Features carry **full Event references** because their preceding Events may originate in other Bounded Contexts.

A Feature is therefore not just a grouping mechanism – it is a frame that determines what kind of facts the Scenarios inside it can talk about, and how.
