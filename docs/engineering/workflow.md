# Engineering Workflow

One work item at a time. The owner (you) approves every gate marked 👤.

| Phase | What happens | Who | Status after |
| --- | --- | --- | --- |
| 1. Discovery | Describe the idea/bug/endpoint; clarify, inspect spec and JS code, compare approaches. No implementation. | `engineering-lead` or `/brainstorm-requirements` | — |
| 2. Alignment | Lead proposes behavior, scope, acceptance criteria. You challenge and revise. | `engineering-lead` | — |
| 3. Work item | Draft saved as `docs/work-items/GO-NNN-slug.md`, added to `INDEX.md`. 👤 You approve. | `/create-work-item` | Draft → Awaiting Approval → **Ready** |
| 4. Implementation | You write the code. Ask for concept help any time. 👤 You mark start. | you + `go-mentor`, `/learn-go-concept` | In Progress |
| 5. Self-verification | You check each criterion and collect evidence (commands, test output, file refs). Optional: `/verify-api-parity`. | you | — |
| 6. Review | Reviewer checks real changes against each criterion. You fix; reviewer re-checks. | `/review-implementation` → `code-reviewer` | In Review → Changes Required ↔ In Review |
| 7. Learning check | Mentor asks a few targeted questions on the key concepts. Misunderstandings are explained, then retried. | `go-mentor` | Ready for Closure |
| 8. Closure | 👤 You authorize. Evidence, decisions, concepts, debt recorded; next task suggested, not started. | `/close-work-item` | **Done** |

`Blocked` may be set at any phase, with the reason recorded.

## Conventions

- IDs: `GO-001`, `GO-002`, … (next number = highest in `INDEX.md` + 1).
- Statuses: Draft, Awaiting Approval, Ready, In Progress, In Review, Changes Required, Ready for Closure, Done, Blocked.
- The status in the work item's header and in `INDEX.md` must match. Agents update them only to a state you explicitly approved.
- Review reports are saved on request to `docs/reviews/GO-NNN-review-N.md` (create the directory when first needed).
- A passing review is necessary but not sufficient for closure; the learning check and your approval are also required.

## Using the agents

Skills run in the main conversation and are the best fit for multi-turn dialogue. Agents can be delegated to from a conversation ("use the go-mentor agent to…") or run as the whole session for back-and-forth work, e.g. `claude --agent go-mentor` (confirm the flag with `claude --help` on your version). Subagents return a single message to the caller, so for long interactive discussions prefer the skill or `--agent` form.
