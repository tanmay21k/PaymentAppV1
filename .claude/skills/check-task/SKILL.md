---
name: check-task
description: Check one task of a work item micro-step by micro-step, report Met, Not met or Unverified for each, and on a full pass write the redo checklist. Use after the owner finishes a task.
argument-hint: "[GO-NNN] [task number, e.g. 2]"
---

Check task: $ARGUMENTS

Act as the `go-coach` agent (see `.claude/agents/go-coach.md`). Follow the root CLAUDE.md. The owner writes all application code; you only check it.

1. **Find it.** Read `docs/work-items/GO-NNN-*/README.md` and `tasks/NN-*.md`. If an earlier task is not `PASSED`, name it, but do not refuse.
2. **Check each micro-step in order.** For each step, run or inspect exactly what its Check line says, and read the touched files. Do not rely on the owner's description. Use read-only commands. Do not start containers, apply migrations, install tools, or change data: those are the owner's steps. If a check needs something not running, answer **Unverified** and say exactly what to run.
3. **Report each step.** For each step give **Met**, **Not met** or **Unverified**, with evidence (command output or `file:line`). Keep it short.
4. **If every step is Met.**
   - Fill the task's `Result` (date and evidence) and set `Status: PASSED`. Tick the task in the README's Tasks list.
   - Draft the **Redo checklist** from what was actually verified: numbered steps to do this task again without AI, each with the concept behind it. Show it and save it after the owner's OK or edits.
   - Name the next task's topic only, and stop. Do not start it.
5. **If any step is Not met.** Name the step, say what differs from its expectation, and give a hint: the concept and where to look. Do not write the solution or edit application code. Leave the later steps unchecked.
6. **Never** change the item's `Status`, tick an acceptance criterion, or close the item. Acceptance criteria are checked at review and closure.
