---
name: verify-api-parity
description: Compare the Go implementation of an endpoint with openapi.yaml and the JS backend and produce a checklist of verified behavior, discrepancies and unverified assumptions.
argument-hint: "[endpoint or GO-NNN]"
---

Verify parity for: $ARGUMENTS

Act as the `api-contract-specialist` agent (see `.claude/agents/api-contract-specialist.md`), or delegate to it. Do not modify files.

1. Identify the operations (from the argument or the work item's sections 6-7).
2. Extract the contract from `openapi.yaml`, the behavior from `backend/`, and the behavior from `go-backend/`, each with file:line.
3. Compare: path/method, auth, request fields (required, optional, null, types, unknown fields), response body and field names/casing, status codes per case, headers, content type, error shape, ID/number/date formats, side effects.
4. If a test run or request is feasible and the owner agrees, run it and record actual output. Otherwise mark the item Unverified; never report behavior as verified without observing it.
5. Output three sections:
   - **Verified** — behavior that matches, with evidence.
   - **Discrepancies** — spec vs JS vs Go, both sides shown. Ask the owner how to resolve each; do not decide. Don't presume JS is correct over the spec.
   - **Unverified assumptions** — what remains unchecked and how to check it.
6. Recommend contract/integration tests that would lock the parity in.
