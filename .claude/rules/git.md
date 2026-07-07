# Git

## Branches

- Use a separate branch; do not work in `main`.
- Use `kebab-case` to briefly describe what the branch is about (for example `update-golang-dependencies`).

## Commits

- Use conventional commit messages.
- Use `feat`, `fix`, or `chore` as the only allowed prefixes.
- Start the message with a capital letter and end it with a period (for example `chore: Update Go dependencies.`).
- Make one logical change per commit.

## History

- Never force push.
- Never use `git push --force-with-lease` either.
- Never amend a commit that has already been pushed (that would require a force push); before pushing, amending is fine only to fix a typo in the commit message.
- Never rebase; always integrate changes by merging.

## Pull Requests

- Open a pull request only after `make qa` passes.
- Pushing a branch and opening a pull request are outward-facing steps; propose them and wait for an explicit go, rather than pushing or opening a pull request unprompted. Commit locally as the work progresses. The exception is a skill whose defined workflow ends in pushing and opening or updating a pull request; there, that final step is authorized as part of running the skill.
- Use the conventional-commit format for the pull request title (same prefixes and capitalization as commits).
- Describe in the pull request body what was changed and why.
- When you add further changes to an existing pull request, check whether the title and the description still cover the full scope, and update them if they do not.
