# AI-Assisted Workflow Guide

> **Deprecated (2026-10-21).** This guide describes the old flow: the `engineering-lead` and `go-mentor` agents, skills such as `/create-work-item` and `/brainstorm-requirements`, nine statuses, and an approval step. It no longer matches the repository. The current flow is in [docs/engineering/workflow.md](../docs/engineering/workflow.md), with a sequence diagram in [WORKFLOW.md](../WORKFLOW.md). This file is kept for reference only.

This guide explains how to use the Claude Code agents and skills set up for the Go backend, and which one to reach for in common situations.

The setup lives at the repository root, not in this folder:

- `CLAUDE.md` — project rules that apply to every agent and skill
- `.claude/agents/` — the four agents
- `.claude/skills/` — the six skills
- `docs/engineering/workflow.md` — the 8-phase workflow
- `docs/work-items/` — work items, the index, and the template
- `docs/learning/LEARNING-LOG.md` — your cumulative Go learning record

Start Claude Code from the repository root (`/home/tanmayk/PaymentAppV1`) so these files load.

---

## The core idea

The AI never writes the implementation for you. The work is split across four roles:

| Role | Who | What it does |
| --- | --- | --- |
| Plan | `engineering-lead` agent, `brainstorm-requirements` and `create-work-item` skills | Clarifies the requirement and writes the work item |
| Teach | `go-mentor` agent, `learn-go-concept` skill | Explains concepts and gives hints |
| Check the contract | `api-contract-specialist` agent, `verify-api-parity` skill | Compares the spec, the JavaScript backend and your Go code |
| Judge | `code-reviewer` agent, `review-implementation` and `close-work-item` skills | Reviews your diff against the acceptance criteria and gates closure |

You write the code. You also approve every status change. Nothing moves to DOING, IN REVIEW or CLOSED unless you say so.

---

## Skills and agents at a glance

### Skills (type as slash commands)

| Command | When to use it |
| --- | --- |
| `/brainstorm-requirements <idea>` | You have an idea, bug, or endpoint to migrate and no agreed requirements yet |
| `/create-work-item <title>` | Requirements are agreed; you want a formal work item (`GO-NNN`) |
| `/learn-go-concept <concept>` | You want to understand one Go concept in depth |
| `/verify-api-parity <endpoint or GO-NNN>` | You want a checklist of how your Go behavior matches the spec and JS code |
| `/review-implementation GO-NNN` | Your implementation is ready for review |
| `/close-work-item GO-NNN` | Review passed and you want to close the task |

`/create-work-item` and `/close-work-item` can only be run by you, not invoked automatically by Claude.

### Agents (ask by name)

Ask for an agent in plain language, for example: *"Use the go-mentor agent to help me understand why…"*

To run an agent for an entire session, you can try `claude --agent go-mentor`. Check `claude --help` to confirm that flag works in your installed version.

| Agent | Use it for | Can it change files? |
| --- | --- | --- |
| `engineering-lead` | Planning, trade-offs, drafting work items | Only `docs/work-items/` |
| `go-mentor` | Concepts, hints, debugging reasoning, the learning check | No |
| `code-reviewer` | Reviewing your changes | No |
| `api-contract-specialist` | Spec vs JS vs Go questions, contract tests | No |

**Skills or agents?** Skills are better for guided, multi-step conversations. An agent returns a single message when delegated, so for a long back-and-forth, use the skill or start the session with `--agent`.

---

## Situations and what to do

### Situation 1: You are starting a new piece of work

Example: "I want to migrate `POST /api/v1/user/signup`" or "I want to start the transfer endpoint."

1. Run `/brainstorm-requirements <what you want to build>`.
   - Claude reads the OpenAPI spec and the JS code first, then asks you questions.
   - Answer them. If Claude reports a conflict between the spec and the JS code, decide which one wins. Claude will not decide for you.
2. Review the requirements summary. Mark anything you're not sure of as an assumption, and ask Claude to challenge it.
3. When you're satisfied, run `/create-work-item <short title>`.
   - Claude assigns the next ID (for example `GO-001`) and saves `docs/work-items/GO-001-slug.md`.
4. Read the work item carefully. Approve it explicitly ("approved") or ask for changes. Only then does it move to DOING.

**Do not** start coding before the work item is approved (`Approved by owner: yes`).

Agent alternative: use `engineering-lead` if you prefer a freer discussion before formalizing anything.

---

### Situation 2: You are stuck and want to understand a concept

Example: "I don't understand why this function takes a `context.Context`" or "What is the difference between returning an error and panicking?"

Two options, depending on what you need:

- **You want a structured lesson on one concept**: run `/learn-go-concept <concept>`. You'll get the mental model, an example, common mistakes, when to use it, how it applies to your task, and an exercise.
  - The exercise answer is withheld until you attempt it. Try it before asking for the answer.
- **You have a specific question about your code**: ask `go-mentor` directly. Tell it:
  1. Which work item you're on (`GO-NNN`)
  2. The file and line you're working on
  3. What you've tried and what you expected

The mentor will:
- Ask what you understand first
- Give you a hint before an explanation
- Give you a small, unrelated example if you remain stuck
- Ask you to apply it yourself

It will not write the solution. If you explicitly want code, it will ask whether you want a conceptual example, a small snippet, or a full implementation, and it will only give the last if you confirm.

**Tip:** mention the concept you're blocked on by name (for example "error wrapping with `%w`"). That lets the mentor use the right material.

---

### Situation 3: You're stuck on a bug or a failing test

1. Run the failing test yourself and keep the output:
   ```bash
   cd go-backend
   go test ./... -run <TestName> -v
   ```
2. Ask `go-mentor` for reasoning help. Include the test output, what you expected, and what you think is happening.
3. Do not ask for the fix. Ask why the behavior occurs.

If the failure is about the API behavior (status codes, response shape), use `api-contract-specialist` or `/verify-api-parity` instead, because the spec is the reference.

---

### Situation 4: You want to know what the API must do

Example: "What does the spec say about `GET /api/v1/user/bulk`?" or "Does the JS version return `404` or `400` for an unknown recipient?"

1. For the contract alone, ask `api-contract-specialist` for the operation's request and response schemas, status codes and auth requirements.
2. For a comparison with your Go code, run `/verify-api-parity <endpoint>`.
   - It produces three lists: **Verified**, **Discrepancies**, and **Unverified assumptions**.
3. If there are discrepancies, decide how each is resolved. Record intentional differences in the work item. Do not let the agent choose the behavior for you.

---

### Situation 5: You've finished implementing a work item

Work through these steps in order.

**Step 1: Self-verify.** Before asking for a review, check every acceptance criterion in the work item and collect evidence. At minimum, run:

```bash
cd go-backend
go build ./...
go vet ./...
go test ./...
```

Keep the output. Each acceptance criterion should map to a test, a command, or a manual request you can point to.

**Step 2: Review.** Run `/review-implementation GO-NNN` and give it your self-verification evidence.

The reviewer will return one of three verdicts:

| Verdict | What to do next |
| --- | --- |
| **PASS** | Go to Step 3 |
| **PASS WITH NON-BLOCKING OBSERVATIONS** | Decide on each observation: fix it now, or record it as accepted debt. Then go to Step 3. |
| **CHANGES REQUIRED** | Fix the findings yourself. Then re-run `/review-implementation GO-NNN`. Repeat until the verdict is PASS or PASS WITH NON-BLOCKING OBSERVATIONS. |

The reviewer explains each finding, but it won't fix anything. If you're unsure how to address a finding, ask `go-mentor` about the concept behind it.

**Step 3: Learning check.** Ask `go-mentor` to run the learning verification for `GO-NNN`. It asks a few targeted questions about the concepts you used. If you get something wrong, it explains and lets you retry.

**Step 4: Close the work item.** Run `/close-work-item GO-NNN`.
- The skill checks evidence, tests, contract parity, review findings, documentation and the learning check.
- If anything is missing, it lists what remains and stops.
- If everything is complete, it asks you to authorize closure. Only after you say yes does it mark the item `CLOSED`, fill in the closure record, and append to `docs/learning/LEARNING-LOG.md`.

**Step 5: Decide what's next.** The skill suggests the next logical task. It does not start it. Start a new cycle with Situation 1 when you're ready.

---

### Situation 6: The spec and the JS code disagree

Example: the spec says a field is required, but the JS code accepts it as optional.

1. Do not pick one silently. Ask `api-contract-specialist` to describe the discrepancy with file and line references from both sides.
2. Decide which behavior is correct. Options include following the spec, keeping the JS behavior, or proposing a spec change.
3. Record your decision in the work item under "Intentional differences" (or in the acceptance criteria), with your reason.

---

### Situation 7: You want a quick review of a small change

If the change is small and does not need a work item (for example, fixing a typo or updating a comment), you can skip the full workflow. Ask `code-reviewer` for a quick look.

Be aware that the reviewer checks the diff against the acceptance criteria, so a change with no work item has no acceptance criteria to check against. Treat any change that adds behavior as a work item.

---

### Situation 8: You don't know what to build next

Run `/brainstorm-requirements` with a general prompt, for example "what should I migrate after signup?" The engineering lead will look at the spec and the existing Go code and suggest options with trade-offs. It will not start work for you.

---

## Rules to remember

- **Requirements before code.** No implementation until the work item is approved.
- **You approve.** Only you can approve drafts, start implementation, and close tasks.
- **Verify before claiming.** Run the tests and read the output yourself. Don't accept "tests pass" without seeing the result.
- **Spec conflicts go to you.** The agents raise them and ask you to decide.
- **One task at a time.** Finish, review, and close a work item before starting another.

---

## Quick reference: Go commands

Run these from the `go-backend/` directory:

```bash
go build ./...                                # compile everything
go vet ./...                                  # static checks
go test ./...                                 # all tests
go test ./internal/auth -run TestName -v      # one test, verbose
```

Database commands depend on the migration and Docker setup in `go-backend/`.
