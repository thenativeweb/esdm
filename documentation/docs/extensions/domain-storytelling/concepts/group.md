# Group

A **Group** is an arbitrary cluster of elements – **[Actors](/extensions/domain-storytelling/concepts/actor.md)**, **[Work Objects](/extensions/domain-storytelling/concepts/work-object.md)**, or edges inside a **[Sentence](/extensions/domain-storytelling/concepts/sentence.md)** – that the story wants to mark together. Groups express things like "these Actors belong to the customer side of the flow", "these Work Objects all relate to a specific subdomain", or "these edges form one logical sub-process within a larger Sentence".

Groups are declared once at the story level, where they carry their own description and annotation, and then referenced from the elements that belong to them.

## Why Decoupled Declaration

The decoupling between top-level declaration and element-level membership is deliberate. It mirrors how a diagram is read: the Group in the legend says *what the cluster means*, and the membership tags on individual elements say *who is in it*.
