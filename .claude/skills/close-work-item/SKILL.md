---
name: close-work-item
description: Close a work item once every task has passed, every acceptance criterion has evidence, and the owner says close. Use only when the owner asks to close it.
argument-hint: "[GO-NNN]"
disable-model-invocation: true
---

Close: $ARGUMENTS

Read `docs/work-items/$ARGUMENTS-*/README.md` and every file in its `tasks/` folder. Change nothing until step 4.

1. **Tasks.** Every task must be `PASSED` and have a redo checklist. List any that are not.
2. **Criteria.** Check each acceptance criterion for evidence (`file:line`, command output, or test name). Re-run the tests and show the output.
3. **Review (optional).** If a review exists in `docs/reviews/` for this item, report its verdict and any open findings. A review is not required to close.
4. **Stop if incomplete.** List what is missing and stop.
5. **Ask.** If everything is in place, ask: "Do you authorize closing $ARGUMENTS?" Continue only on an explicit yes.
6. **Close.** Set `Status: CLOSED` in the README and in `INDEX.md`. Fill the README's Closure section briefly: evidence, decisions, debt and follow-ups. The redo checklists are the learning record, so do not add a learning-log entry. Suggest the next task, and do not start it.
