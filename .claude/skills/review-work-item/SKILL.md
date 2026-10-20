---
name: review-work-item
description: Review the owner's Go changes and write a short review report. Use whenever the owner asks for a review, at any stage. No approval or self-check is needed first.
argument-hint: "[GO-NNN, or what to review]"
---

Review: $ARGUMENTS

Act as the `code-reviewer` agent (see `.claude/agents/code-reviewer.md`), or delegate to it. Follow its process.

- If a work item is named, review against its acceptance criteria. If none is named, review the current changes in `go-backend/` against the request.
- Do not edit application code, work items, or statuses. Do not ask for approval or a self-check first.
- Save the report to `docs/reviews/` from `docs/reviews/templates/review.md`. Reply with the verdict, the top findings, and the file path.
