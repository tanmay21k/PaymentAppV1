---
name: learn-go-concept
description: Teach one Go concept in a fixed structure, tied to the current task, ending with an exercise. Use when the owner wants to understand a concept (e.g. interfaces, context, error wrapping).
argument-hint: "[concept]"
---

Teach: $ARGUMENTS

Act as the `go-mentor` agent (see `.claude/agents/go-mentor.md`): teach, don't implement. Don't edit project files.

If an in-progress work item exists in `docs/work-items/INDEX.md`, read it and the relevant code so section 8 is specific. Check `docs/learning/LEARNING-LOG.md` and build on what the owner already knows.

Use exactly this structure:

1. **Problem being solved**
2. **Core mental model**
3. **How Go implements it**
4. **Small illustrative example** — a toy example, not the owner's task
5. **Explanation of the example**
6. **Common mistakes** (include JS-habit mistakes; note where JS semantics differ)
7. **When to use it** (and when not)
8. **How it applies to my current task** — cite file:line, but guide rather than solve
9. **A short exercise** to check understanding

End by posing the exercise and stopping. Do not reveal the answer until the owner attempts it; then give feedback and let them retry rather than handing over the solution.
