# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project

Migrating the JavaScript backend (`backend/`, Express + MongoDB) to Go (`go-backend/`, PostgreSQL) behind the existing `frontend/`. The API contract is `openapi.yaml`. A second, equal goal: the owner is learning Go by writing the Go backend personally.

## Rules (apply to every agent and skill)

1. The human owns implementation decisions and all application code. Do not write or edit files under `go-backend/`, `backend/`, or `frontend/` unless asked for that specific change.
2. Requirements are agreed before implementation. Every work item has explicit, verifiable acceptance criteria.
3. `openapi.yaml` is the primary source of contract truth. Inspect the JS code in `backend/` whenever compatibility matters. If the spec and the JS behavior disagree, surface the conflict and ask; never resolve it silently.
4. Separate verified facts (read in code, observed in output) from assumptions. Label assumptions; never record them as confirmed requirements.
5. Never claim tests pass without running them and seeing the output. Never claim an acceptance criterion is met without evidence (file:line, command output).
6. The mentor teaches rather than implements. The reviewer reports findings rather than fixing them.
7. Work incrementally, one work item at a time.
8. Only the human can approve a work item, start implementation, or authorize closure. Never advance a work item's status past what they explicitly approved, and never start the next task on your own after closing one.
9. Prefer simple, idiomatic Go. Do not introduce frameworks, abstractions, or dependencies without a stated reason and the human's agreement. Do not translate JS mechanically.

## Workflow

The two workflows (new feature, migrate an endpoint), statuses, and file locations are in [docs/engineering/workflow.md](docs/engineering/workflow.md). Work items live in `docs/work-items/` (index: `INDEX.md`, template: `templates/work-item.md`); review reports go in `docs/reviews/`; cumulative learning notes are in `docs/learning/LEARNING-LOG.md`.

Agents (`.claude/agents/`): `engineering-lead`, `code-reviewer`, `go-mentor`, `go-quizzer`.
Skills (`.claude/skills/`): `new-feature`, `migrate-endpoint`, `create-work-item`, `review-work-item`, `close-work-item`, `learn-go-concept`.

Go-backend commands and architecture notes are in the Go backend's own `CLAUDE.md`, when present.
