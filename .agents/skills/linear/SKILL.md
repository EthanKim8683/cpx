---
name: linear
description: Create, read, and update Linear issues and comments via the Linear MCP server. Use when working with Linear issues, issue status, or synced issue threads — not for commits, branches, or pull requests.
---

# Linear

Issue tracking in cpx. Use the **Linear MCP server** (`https://mcp.linear.app/mcp`). Read each tool's schema before calling it.

If MCP calls fail with an auth error, ask the user to connect Linear in their agent's MCP settings, then retry.

## cpx conventions

**Use for:** creating and updating issues, tracking status, recording significant updates on issue threads.

**Do not use for:** commits, branches, or pull requests — use the [github](../github/SKILL.md) skill.

- Default to Linear for anything that is not code.
- GitHub issues are a sync mirror — do not create issues on GitHub or drive issue workflow there.
- On synced issues, follow [docs/ai/issues.md](../../docs/ai/issues.md): reply to the synced system comment (`parentId`), not a top-level comment.

Default team: **Competitive Programming**. Issue IDs look like `CP-123`.

Issue bodies stay short. Put decisions, findings, and progress in comments.

## How to call tools

1. Read the tool schema in the MCP descriptor before calling.
2. Pass markdown in `description` and `body` as literal text — real newlines, not `\n` escape sequences.
3. Read before write: `get_issue` / `list_issues` / `list_comments`, then `save_issue` / `save_comment`.
4. Post issue comments **one at a time** — parallel replies arrive out of order.

## Issues

### Read

| Tool | When |
| --- | --- |
| `get_issue` | One issue by ID (`CP-123`). Set `includeRelations` when you need blockers/duplicates. |
| `list_issues` | Search or filter. Use `team`, `state`, `assignee` (`me`), `query`, `project`, `label`. |

### Create

`save_issue` — omit `id`. Required: `title`, `team` (`Competitive Programming`).

Keep `description` to a line or two. Optional: `state`, `priority` (0=None … 4=Low), `project`, `labels`, `assignee` (`me`).

### Update

`save_issue` — pass `id` (`CP-123` or UUID). Only include fields to change: `state`, `title`, `description`, `priority`, `assignee`, etc.

Use `list_issue_statuses` when you need valid state names for the team.

## Comments

Synced issues have a system comment linked to GitHub. **Reply to that comment**, not the issue top-level.

1. `list_comments` with `issueId` — find the synced system comment (top-level, links to GitHub).
2. `save_comment` with `parentId` set to that comment's ID and your `body`.

For a new top-level thread (rare in cpx): `save_comment` with `issueId` and `body`, no `parentId`.

To edit an existing comment: `save_comment` with `id` and updated `body`.

## Other tools

Use when needed — not part of the default workflow:

| Tool | Purpose |
| --- | --- |
| `list_projects` / `get_project` / `save_project` | Projects |
| `list_documents` / `get_document` / `save_document` | Linear docs |
| `search_documentation` | Search Linear help docs |
| `list_issue_labels` / `create_issue_label` | Labels |
| `create_attachment` / `get_attachment` | File attachments on issues |
