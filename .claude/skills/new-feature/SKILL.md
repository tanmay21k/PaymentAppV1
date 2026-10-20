---
name: new-feature
description: Take a new feature for the Go backend from idea to a saved work item that the owner can approve. Use when the owner wants behavior that does not come from the JS backend.
argument-hint: "[feature idea]"
disable-model-invocation: true
---

Feature: $ARGUMENTS

Act as the `engineering-lead` agent (see `.claude/agents/engineering-lead.md`). Follow the root CLAUDE.md.

1. **Ground.** Read `openapi.yaml` if an API is involved, and the relevant `go-backend/` code. Cite file:line.
2. **Ask.** Ask 2-4 questions that change the design, then wait for the answers.
3. **Propose.** Give one recommended approach (one alternative at most), the scope, the non-goals, and draft acceptance criteria. Label each behavior **[Confirmed]** or **[Assumption]**.
4. **Save.** Once the owner agrees on the main points, run `/create-work-item` to save the work item (`Status: DOING`, `Approved by owner: no`) and its `INDEX.md` row.
5. **Approve.** Show the summary and ask for approval. Set `Approved by owner: yes` only after an explicit approval. Do not start implementation.
