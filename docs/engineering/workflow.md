# Engineering Workflow

One work item at a time. The owner (you) approves every gate marked 👤.

| Phase | What happens | Who | Status after |
| --- | --- | --- | --- |
| 1. Discovery | Describe the idea/bug/endpoint; clarify, inspect spec and JS code, compare approaches. No implementation. | `engineering-lead` or `/brainstorm-requirements` | — |
| 2. Alignment | Lead proposes behavior, scope, acceptance criteria. You challenge and revise. | `engineering-lead` | — |
| 3. Work item | Item saved as `docs/work-items/GO-NNN-slug.md`, added to `INDEX.md`. 👤 You approve (header: `Approved by owner: yes`). | `/create-work-item` | **DOING** |
| 4. Implementation | You write the code. Ask for concept help any time. 👤 You mark start. | you + `go-mentor`, `/learn-go-concept` | **DOING** |
| 5. Self-verification | You check each criterion and collect evidence (commands, test output, file refs). Optional: `/verify-api-parity`. | you | — |
| 6. Review | Reviewer checks real changes against each criterion. You fix; reviewer re-checks. | `/review-implementation` → `code-reviewer` | **IN REVIEW** (back to DOING while changes are fixed) |
| 7. Learning check | Mentor asks a few targeted questions on the key concepts. Misunderstandings are explained, then retried. | `go-mentor` | IN REVIEW |
| 8. Closure | 👤 You authorize. Evidence, decisions, concepts, debt recorded; next task suggested, not started. | `/close-work-item` | **CLOSED** |

Only three statuses exist: **DOING**, **IN REVIEW**, **CLOSED**. Note a blocker in the work item's text; the status stays DOING.

## Conventions

- IDs: `GO-001`, `GO-002`, … (next number = highest in `INDEX.md` + 1).
- Statuses: DOING (includes drafting, approval and changes), IN REVIEW (includes ready for closure), CLOSED.
- The status in the work item's header and in `INDEX.md` must match. Agents update them only to a state you explicitly approved.
- Review reports are saved on request to `docs/reviews/GO-NNN-review-N.md` (create the directory when first needed).
- A passing review is necessary but not sufficient for closure; the learning check and your approval are also required.

## Using the agents

Skills run in the main conversation and are the best fit for multi-turn dialogue. Agents can be delegated to from a conversation ("use the go-mentor agent to…") or run as the whole session for back-and-forth work, e.g. `claude --agent go-mentor` (confirm the flag with `claude --help` on your version). Subagents return a single message to the caller, so for long interactive discussions prefer the skill or `--agent` form.
