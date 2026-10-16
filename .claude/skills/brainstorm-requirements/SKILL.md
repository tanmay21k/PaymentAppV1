---
name: brainstorm-requirements
description: Explore an idea, bug or migration need and reach agreed requirements before any work item or code exists. Use at the start of a new task.
argument-hint: "[idea, endpoint, or problem]"
---

Explore: $ARGUMENTS

Act as the `engineering-lead` agent (see `.claude/agents/engineering-lead.md`). Follow the root CLAUDE.md rules.

1. **Ground yourself.** Read the relevant `openapi.yaml` operations, the JS implementation in `backend/`, and the current Go code in `go-backend/`. Summarize what you found with file:line references, separating verified facts from assumptions.
2. **Ask clarifying questions.** Ask the few that matter most (2-5) and wait for answers. Cover: goal, scope, constraints, behavior on errors, compatibility with the frontend, what to leave out.
3. **Surface conflicts.** If the spec and JS disagree, present both and ask the owner to decide.
4. **Compare approaches** (2-3 when real alternatives exist), with trade-offs, a recommendation and reasoning. Don't presume a framework or library before discussing it.
5. **List** risks, edge cases, dependencies, likely affected files, Go concepts needed, and a testing approach.
6. **Converge.** Propose a requirements summary: behavior, scope, non-goals, draft acceptance criteria, with each item marked [Confirmed] or [Assumption]. Ask the owner to challenge it.

Do not produce a final implementation plan while important ambiguities remain open. Do not write code or files. When the owner confirms the requirements, suggest `/create-work-item` — do not run it.
