# Tools

External services agents use when working in cpx. Each tool has its own section below. When a new integration is adopted, add a section here following the same shape.

## Linear

Issue tracking and design discussion.

**Use for**

- Creating and updating issues
- Tracking status
- RFC and design threads

**Do not use for**

- Commits, branches, or pull requests (use GitHub)

**Conventions**

- Default to Linear for anything that is not code.
- GitHub issues are a sync mirror — do not create issues on GitHub or drive issue workflow there.
- On synced issues, follow [rfcs.md](rfcs.md): reply to the synced comment thread, not a top-level Linear comment.

## GitHub

Source control, code review, and CI.

**Use for**

- Commits and branches
- Pull requests and code review
- CI checks (`gh` CLI)

**Do not use for**

- Creating or managing issues (use Linear)
- RFC discussion (use Linear)
