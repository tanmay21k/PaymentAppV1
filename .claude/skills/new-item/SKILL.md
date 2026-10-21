---
name: new-item
description: Start a new work item by asking questions until the requirements are agreed, then save its README. Use for a new feature or for migrating one JS endpoint to Go. Writes no tasks.
argument-hint: "[feature idea | METHOD path, e.g. POST /api/v1/user/signup]"
disable-model-invocation: true
---

New item: $ARGUMENTS

Act as the `go-coach` agent (see `.claude/agents/go-coach.md`). Follow the root CLAUDE.md.

1. **Ground.** Read the relevant `go-backend/` code. For a migration, also read the endpoint in `openapi.yaml` and in `backend/` (route, middleware, schema, validation, status codes), citing file:line, and note what already exists in Go.
2. **Ask.** Ask one or two questions at a time that change the design, and wait for the answers. For a migration, list each spec-vs-JS conflict with both sides and ask which one Go follows. Do not pick for the owner.
3. **Converge.** Propose one recommended approach (one alternative at most), the scope, the non-goals, and draft acceptance criteria. Label each behaviour **[Confirmed]** or **[Assumption]**. Ask the owner to challenge it.
4. **Save.** Once the owner agrees on the main points, assign the next ID (the highest in `docs/work-items/INDEX.md` plus one), fill `docs/work-items/templates/README.md`, and save it as `docs/work-items/GO-NNN-slug/README.md` with `Status: TODO`. Add the `INDEX.md` row. Acceptance criteria must be observable and say how they are verified.
5. **Stop.** Do not write any task, and do not start implementation. Tell the owner to run `/next-task GO-NNN` when they want to begin.
