# Workflow

You work like Boot.dev: the coach asks questions, you reason out each step, you write the code, and a check tells you whether it worked. Each work item is a folder with a sequence of tasks, so you can reuse the sequences later as a reference when you code without AI.

The sequence diagram is in [WORKFLOW.md](../../WORKFLOW.md).

## Steps

1. `/new-item <idea | METHOD path>`: the coach asks questions, one or two at a time, until the requirements are agreed. For a migration it lists the spec-vs-JS conflicts for you to decide. It saves `GO-NNN-slug/README.md` (status `TODO`). No tasks yet.
2. `/next-task GO-NNN`: the coach asks questions so you reason out the next step, then writes that one task file (goal, concepts, check). The first task moves the item to `DOING`.
3. You write the code in `go-backend/`. Ask the coach whenever you are unsure: it gives hints, never the solution.
4. `/check-task GO-NNN <n>`: the coach runs the task's check and answers Met, Not met or Unverified, with evidence. On Met it drafts the **redo checklist** for you to confirm.
5. Repeat steps 2 to 4 until the acceptance criteria are covered.
6. `/review-work-item` whenever you want a review.
7. `/close-work-item GO-NNN` when every task has passed. You confirm closure.

## Layout

```
docs/work-items/
├── INDEX.md
├── templates/            (README.md for items, task.md for tasks)
└── GO-NNN-slug/
    ├── README.md         (problem, scope, acceptance criteria, open questions, task list, closure)
    └── tasks/
        ├── 01-slug.md    (goal, concepts, check, result, redo checklist)
        └── 02-slug.md
```

## Rules

<<<<<<< Updated upstream
- One work item at a time. Statuses are **DOING** and **CLOSED**.
- Only you set `Approved by owner: yes`, and only you authorize closure.
- A review needs no approval and no self-check first. It never changes a status.
- Learning help is optional: `/learn-go-concept <concept>`, or ask the `go-mentor` agent a question.
- Quizzes on closed work: ask the `go-quizzer` agent. It records results in `docs/learning/QUIZ-PROGRESS.md`.
=======
- One work item at a time. Statuses are **TODO** (created, no task started), **DOING** (first task started) and **CLOSED**.
- Only you start an item (by starting its first task) and authorize closure. There is no approval field.
- Tasks are written one at a time. Task N+1 exists only after task N has passed.
- The redo checklists are the learning record: numbered steps to repeat each task without AI, each with its concept.
- A review needs no start and no self-check first. It never changes a status.
- Learning help: `/learn-go-concept <concept>`, or ask the coach. Quizzes on closed work: ask the `go-quizzer` agent; it records results in `docs/learning/QUIZ-PROGRESS.md`. `docs/learning/LEARNING-LOG.md` is kept as history and is no longer appended.
>>>>>>> Stashed changes
- IDs: `GO-NNN`. The next number is the highest in `INDEX.md` plus one.
- Review reports: `docs/reviews/GO-NNN-review-N.md`, from [templates/review.md](../reviews/templates/review.md).
- Templates: [item README](../work-items/templates/README.md) and [task](../work-items/templates/task.md).

## Agents

| Agent | Use it for |
| --- | --- |
| `go-coach` | Asking questions to plan an item, guiding each task, checking it, and writing the redo checklist. Writes only in `docs/work-items/`. Never writes application code. |
| `code-reviewer` | Reviews on request. Writes only in `docs/reviews/`. Never fixes code. |
| `go-quizzer` | Quizzes on closed work, from `QUIZ-TOPICS.md` and the redo checklists. Writes only `QUIZ-PROGRESS.md`. |
