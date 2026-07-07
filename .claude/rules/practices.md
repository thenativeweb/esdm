# Development Practices

## Goal-Driven Execution

- Before starting, restate the task as a concrete, verifiable goal: what must be true once it is done.
- Derive the success criteria up front (a passing test, an observable behavior, a green `make qa`) and work towards them, instead of executing step by step without a finish line.

## Test-Driven Development

- Follow test-driven development without exception, using the red-green-refactor cycle:
  - Red: write a failing test first.
  - Green: make it pass with the simplest change.
  - Refactor: improve the code while the tests stay green.
- The refactor step needs no new test: improving existing code (naming, structure, readability) is allowed without adding a test, as long as the existing tests stay green.

## Debugging and Bug Fixing

- Fix the root cause, not the symptom. Identify and address the underlying cause instead of patching the visible symptom.
- For a bug fix, the red step of the cycle above is a test that reproduces the bug: write that failing test first, then make it pass with the fix.
- When debugging, assume nothing: suspect your own code first (not the compiler, libraries, or OS) and prove assumptions by observation instead of guessing.

## Design

- Favor orthogonality and loose coupling: keep components independent so that a change in one place does not ripple into unrelated ones.

## Simplicity

- Write the minimum code that solves the problem at hand; nothing speculative.
  - No features beyond what was asked, no abstraction only a hypothetical future would need, no configurability nobody requested.
  - Before adding an abstraction, ask whether a senior engineer would call it overcomplicated. If so, simplify.

## Scope Discipline

- Follow the boy scout rule, but only within the scope of your task: leave the logical unit you are already changing (the function, method, or type) cleaner than you found it. This clean-up is the refactor step of the cycle above.
- Do not fix issues outside that scope (other functions, other files) on your own; point them out in your reply and suggest a fix instead.
- Exception: trivial, risk-free issues right next to your change (a typo in a comment, stray whitespace) may be fixed directly; mention them in your reply.
