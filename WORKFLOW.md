# Work flow

The question-led flow for the Go backend migration. The coach asks, you reason it out and write the code, and a check tells you whether it worked. Details are in [docs/engineering/workflow.md](docs/engineering/workflow.md).

```mermaid
sequenceDiagram
    actor Owner
    participant Coach as go-coach
    participant Code as go-backend
    participant Docs as docs/work-items
    participant Reviewer as code-reviewer

    Note over Owner,Docs: 1. Plan the item
    Owner->>Coach: /new-item (idea or METHOD path)
    loop until the requirements are agreed
        Coach->>Owner: one or two questions
        Owner->>Coach: answers and decisions
    end
    Coach->>Docs: save GO-NNN-slug/README.md (TODO) and the INDEX row

    Note over Owner,Docs: 2. Work one task at a time
    loop until the acceptance criteria are covered
        Owner->>Coach: /next-task GO-NNN
        Coach->>Owner: questions to reason out the next step
        Owner->>Coach: answers
        Coach->>Docs: write tasks/NN-slug.md (goal, concepts, check)
        Note right of Docs: the first task sets the item to DOING
        Owner->>Code: write the code
        opt when stuck
            Owner->>Coach: question
            Coach->>Owner: hint, never the solution
        end
        Owner->>Coach: /check-task GO-NNN n
        Coach->>Code: run the check (read-only)
        alt Met
            Coach->>Owner: draft redo checklist
            Owner->>Coach: confirm or edit
            Coach->>Docs: task PASSED, redo checklist saved
        else Not met or Unverified
            Coach->>Owner: hint, or what to run
        end
    end

    Note over Owner,Docs: 3. Review and close
    opt any time
        Owner->>Reviewer: /review-work-item
        Reviewer->>Docs: docs/reviews/GO-NNN-review-N.md
    end
    Owner->>Coach: /close-work-item GO-NNN
    Coach->>Owner: every task passed? authorize closing?
    Owner->>Coach: yes
    Coach->>Docs: CLOSED in the README and INDEX
```

## Statuses

| Status | Meaning |
| --- | --- |
| TODO | The item is saved and no task has started. |
| DOING | The first task has started. |
| CLOSED | Every task passed and you authorized closing. |
