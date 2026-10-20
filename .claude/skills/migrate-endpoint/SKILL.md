---
name: migrate-endpoint
description: Plan the migration of one JS endpoint to Go. Extracts its contract from openapi.yaml and backend/, lists spec-vs-JS conflicts for the owner to decide, and saves a work item. Use when the owner wants to migrate an endpoint or feature from backend/ to go-backend/.
argument-hint: "[METHOD path, e.g. POST /api/accounts]"
disable-model-invocation: true
---

Migrate: $ARGUMENTS

Act as the `engineering-lead` agent (see `.claude/agents/engineering-lead.md`). Follow the root CLAUDE.md.

1. **Extract the contract.**
   - `openapi.yaml`: path, method, auth, request schema (required, optional, nullable, types), responses and status codes. Cite lines.
   - `backend/`: the route, middleware, schema, validation, status codes, error body, field names, and null or absent handling. Cite file:line.
   - `go-backend/`: what already exists for this endpoint.
2. **Compare.** Mark each point **Documented** (spec), **Observed** (JS code) or **Inferred**. For every spec-vs-JS conflict, show both sides and ask the owner which one Go should follow. Do not pick one.
3. **Ask.** Ask the owner the open questions, then wait for the answers.
4. **Save.** Run `/create-work-item` with `Type: migration`. Include the contract section, the expected behavior for each case, acceptance criteria for each status code and field, and the contract tests to write.
5. **Approve.** Show the summary and ask for approval. Set `Approved by owner: yes` only after an explicit approval. Do not start implementation.
