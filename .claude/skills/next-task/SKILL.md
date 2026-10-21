---
name: next-task
description: Work out the next task of a work item by asking the owner questions, then save it as a task file. Use when the owner is ready to start or continue an item. Writes one task and no code.
argument-hint: "[GO-NNN]"
disable-model-invocation: true
---

Next task: $ARGUMENTS

Act as the `go-coach` agent (see `.claude/agents/go-coach.md`). Follow the root CLAUDE.md.

1. **Read.** Read `docs/work-items/GO-NNN-*/README.md` and every file in its `tasks/` folder. If the latest task is not `PASSED`, say so and ask whether to continue with it (`/check-task`) or start another.
2. **Ask first.** Ask one or two questions that let the owner reason out the next step: what the goal is, what they would try, what could go wrong. Wait for the answers. Use the README's acceptance criteria and open questions to decide what to ask. Open questions the owner has not decided come up here.
3. **Agree.** Once the owner has reasoned out the step, write `docs/work-items/GO-NNN-*/tasks/NN-slug.md` from `docs/work-items/templates/task.md` (`Status: OPEN`): the goal, the concepts, and a concrete check. Never include code or the solution. Add the task to the README's Tasks list, unticked.
4. **Go-ahead.** If this is the item's first task, set `Status: DOING` in the README and in `INDEX.md`.
5. **Stop.** Tell the owner to write the code and run `/check-task GO-NNN <n>` when they are done, or to ask questions at any time. Do not start the following task.
