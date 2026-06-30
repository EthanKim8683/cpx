# Agent Guidelines: Issue Tracking

AI agents must keep GitHub Issues updated whenever working on or discussing tasks tied to a known issue.

---

## 1. Single Source of Truth: GitHub Issues

All issue discussions, progress updates, and technical notes must be written to **GitHub Issues** via the GitHub CLI (`gh`).

- **Linear is read-only**: GitHub comments sync automatically to Linear, but Linear comments do not reliably sync back to GitHub.
- **Open-source transparency**: Keeping context on GitHub ensures external contributors have full visibility without needing third-party tool access.

---

## 2. When to Comment (Significance Triggers)

Post a comment on the corresponding GitHub issue whenever significant information emerges:

- **Design & Architecture**: Architectural choices, trade-offs, or changes in strategy.
- **Context & Discoveries**: Technical findings, upstream compiler details, code references, or useful resources.
- **Progress, Roadblocks & Next Steps**: Status updates, milestones reached, encountered blockers, or changes in work ordering and prerequisites (e.g., *"Must implement X before Y"*).

---

## 3. How to Update Issues

Use the GitHub CLI (`gh`) to post clear, free-form markdown comments:

```bash
gh issue comment <issue-number> --body "<comment-text>"
```

> [!NOTE]
> If working from a Linear issue ID, locate the corresponding GitHub issue number first before commenting.
