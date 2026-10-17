---
name: close-work-item
description: Run the final verification gate for a work item and, only with the owner's explicit approval, mark it CLOSED and record evidence, decisions and learning. Never closes a task just because code exists.
argument-hint: "[GO-NNN]"
disable-model-invocation: true
---

Close work item: $ARGUMENTS

Read `docs/work-items/$ARGUMENTS-*.md`, then verify each item below. Report each as Done/Not done with evidence; do not change any status yet.

1. **Every acceptance criterion has evidence** (file:line, command output, test name). Re-run the tests now and show the actual results.
2. **Relevant tests pass** (observed in this session).
3. **OpenAPI contract respected** — run `/verify-api-parity` logic for the item's operations if not already done; any discrepancy must be resolved or an owner-approved intentional difference recorded.
4. **Review** — the latest review outcome is PASS or PASS WITH NON-BLOCKING OBSERVATIONS, and every finding is fixed or explicitly accepted by the owner (list them).
5. **Documentation** updated where needed (README, `openapi.yaml`, Go backend notes).
6. **Learning check done** — the `go-mentor` learning verification was completed and the owner can explain the key concepts and design decisions. If not, stop and run it.

If anything is not done, list what remains, set nothing, and stop.

If everything is done, summarize and ask: "Do you authorize closing $ARGUMENTS?" Proceed only after an explicit yes. Then:

- Set status `CLOSED` in the work item and `INDEX.md`.
- Fill the closure record: verification evidence, design decisions, Go concepts learned, remaining technical debt.
- Append an entry to `docs/learning/LEARNING-LOG.md` (concepts, understood well, needs practice, mistakes and lessons, rationale, optional follow-ups).
- Name the next logical task as a suggestion only. Do not start it, draft it, or assign exercises without the owner's agreement.
