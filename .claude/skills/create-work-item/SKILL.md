---
name: create-work-item
description: Turn agreed requirements into a structured work item (GO-NNN) with acceptance criteria, saved under docs/work-items/. Use only after requirements are agreed.
argument-hint: "[short title]"
disable-model-invocation: true
---

Create a work item for: $ARGUMENTS

1. **Check readiness.** Requirements must have been agreed in this conversation (e.g. via `/brainstorm-requirements`). If important ambiguities remain, list them and stop; do not draft.
2. **Assign an ID.** Read `docs/work-items/INDEX.md`; use the highest GO-NNN + 1 (start at GO-001). Slug the title in kebab-case.
3. **Fill the template** `docs/work-items/templates/work-item.md` completely, all 15 sections plus the empty closure record. Rules:
   - Re-read the cited `openapi.yaml` operations and `backend/` code so references (file:line) are accurate.
   - Acceptance criteria must each be observable and say how they are verified (a command, a request/response, a test name). Reject vague wording like "works correctly".
   - Tag behaviors [Confirmed] or [Assumption]; unresolved assumptions go under risks/open questions.
   - Make the item small and independently verifiable; split it if it isn't.
   - List the Go concepts the owner will need.
4. **Save** to `docs/work-items/GO-NNN-slug.md` with status `DOING`, and add its row to `INDEX.md`.
5. **Stop and ask for approval.** Show the owner the item summary and ask them to approve or request changes. Only after their explicit "approved" record `Approved by owner: yes` in the file header. The status stays `DOING`, and the index row must match. Do not start implementation or any other task.
