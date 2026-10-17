---
name: review-implementation
description: Review the owner's changes against an approved work item and produce a findings report with a verdict. Use after self-verification. Never edits code.
argument-hint: "[GO-NNN]"
---

Review work item: $ARGUMENTS

Act as the `code-reviewer` agent (see `.claude/agents/code-reviewer.md`), or delegate to it. Do not modify any file.

1. Locate `docs/work-items/$ARGUMENTS-*.md` (ask which item if none given or several are DOING/IN REVIEW). Confirm `Approved by owner: yes` in the header; if not, stop and say so.
2. Ask the owner for their self-verification evidence if they haven't given it (checklist, test output, commands). Treat it as a claim to check, not as proof.
3. Inspect the real changes: `git status`, `git diff`, and the touched files. Run the build and tests yourself and report actual output; say explicitly what you could not run.
4. Check every acceptance criterion individually: **Met / Not met / Unverified**, each with evidence.
5. Review the full list in the code-reviewer agent: correctness, OpenAPI compliance, JS compatibility, idiomatic Go, structure, errors, security, concurrency, resources, test quality, maintainability, logging, complexity, edge cases.
6. Report findings in the required format (severity, location, issue, why it matters, concrete scenario, conceptual fix, how to verify), distinguishing defects from style opinions, explaining important ones so the owner learns. Don't pad; if it's sound, say so.
7. Give one verdict: PASS, PASS WITH NON-BLOCKING OBSERVATIONS, or CHANGES REQUIRED. Recommend the status (IN REVIEW, or DOING if changes are required) but let the owner confirm before it is changed.
8. Offer to save the report to `docs/reviews/GO-NNN-review-N.md`; save only if the owner agrees. Remind them that closure also needs the learning check and their approval.

Never fix the findings. The owner makes the changes, then asks for a re-review.
