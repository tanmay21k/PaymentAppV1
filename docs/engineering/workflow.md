# Workflow

Two ways to start work. Both end the same way.

## 1. New feature

1. `/new-feature <idea>`: agree on the requirements, then `/create-work-item` saves the work item (`Approved by owner: no`).
2. You approve the work item (`Approved by owner: yes`).
3. You implement it in `go-backend/`.
4. `/review-work-item` whenever you want a review. Fix what it finds, then ask again if you want.
5. `/close-work-item GO-NNN` when you're ready. You confirm closure.

## 2. Migrate an endpoint from JS

1. `/migrate-endpoint <METHOD> <path>`: extract the contract from `openapi.yaml` and `backend/`, list spec-vs-JS conflicts, and `/create-work-item` saves the work item.
2. You decide each conflict and approve the work item.
3. You implement it in `go-backend/`.
4. `/review-work-item` (it checks parity with the contract too), then `/close-work-item GO-NNN`.

## Rules

- One work item at a time. Statuses are **DOING** and **CLOSED**.
- Only you set `Approved by owner: yes`, and only you authorize closure.
- A review needs no approval and no self-check first. It never changes a status.
- Learning help is optional: `/learn-go-concept <concept>`, or ask the `go-mentor` agent a question.
- Quizzes on closed work: ask the `go-quizzer` agent. It records results in `docs/learning/QUIZ-PROGRESS.md`.
- IDs: `GO-NNN`. The next number is the highest in `INDEX.md` plus one.
- Review reports: `docs/reviews/GO-NNN-review-N.md`, from [templates/review.md](../reviews/templates/review.md).
- Work items: `docs/work-items/GO-NNN-slug.md`, from [templates/work-item.md](../work-items/templates/work-item.md).

## Agents

| Agent | Use it for |
| --- | --- |
| `engineering-lead` | Planning both flows and drafting work items. Writes only in `docs/work-items/`. |
| `code-reviewer` | Reviews on request. Writes only in `docs/reviews/`. Never fixes code. |
| `go-mentor` | Teaching and hints on request. Never writes the solution. |
| `go-quizzer` | Quizzes on closed work, from `QUIZ-TOPICS.md`. Writes only `QUIZ-PROGRESS.md`. |
