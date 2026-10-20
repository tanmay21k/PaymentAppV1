---
name: close-work-item
description: Close a work item once every acceptance criterion has evidence, the latest review passes, and the owner approves closure. Use only when the owner asks to close it.
argument-hint: "[GO-NNN]"
disable-model-invocation: true
---

Close: $ARGUMENTS

Read `docs/work-items/$ARGUMENTS-*.md`. Change nothing until step 4.

1. **Criteria.** Check each acceptance criterion for evidence (`file:line`, command output, or test name). Re-run the tests and show the output.
2. **Review.** Find the latest review in `docs/reviews/` for this item. Its verdict must be PASS or PASS WITH NON-BLOCKING OBSERVATIONS, or the owner must name each open finding they accept. If there is no review, say so.
3. **Stop if incomplete.** List what is missing and stop.
4. **Ask.** If everything is in place, ask: "Do you authorize closing $ARGUMENTS?" Continue only on an explicit yes.
5. **Close.** Set `Status: CLOSED` in the work item and in `INDEX.md`. Fill the Closure section. Append a short entry to `docs/learning/LEARNING-LOG.md`. Suggest the next task, and do not start it.
