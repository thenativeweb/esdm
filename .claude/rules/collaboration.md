# Collaboration

## User-Facing Design Requires Explicit Sign-Off

- User-facing design choices always require explicit sign-off: never decide them autonomously, not even under auto mode. This covers CLI command and flag names, argument shapes (positional vs. flag), default values, default behavior, output formatting and visual layout, the wording of diagnostics, and the wording of error messages shown for ill-formed user input.
- The final say lies with the maintainers; for external contributions, raise open design questions in the issue or pull request before implementing.
- A "go ahead" or similar greenlight only authorizes what was explicitly discussed beforehand. Open design questions still require explicit sign-off regardless of the surrounding workflow mode.
- When an unspecified user-facing aspect surfaces during implementation, stop and ask. The friction of one extra clarification is far smaller than the friction of reworking an interface afterwards. If unsure whether a topic was clarified, default to asking.
