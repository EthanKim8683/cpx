---
name: github
description: Commits, branches, pull requests, and CI in cpx via the gh CLI. Use when creating PRs, checking CI, or working with git remotes — not for creating or managing issues.
---

# GitHub

Source control and CI in cpx. Use the **`gh` CLI** from the repo root when possible. Outside a git directory, pass `--repo owner/repo`.

Run `gh auth status` before the first `gh` call in a session. If unauthenticated, ask the user to run `gh auth login`.

## cpx conventions

**Use for:** commits and branches, pull requests and code review, CI checks.

**Do not use for:** creating or managing issues, or issue discussion — use the [linear](../linear/SKILL.md) skill.

## Branches and commits

```bash
git checkout -b branch-name
git add …
git commit -m "…"
git push -u origin HEAD
```

Use the branch name from the Linear issue when one is provided (`gitBranchName` on the issue).

## Pull requests

Create:

```bash
gh pr create --title "…" --body "$(cat <<'EOF'
## Summary
…

## Test plan
- [ ] …
EOF
)"
```

Inspect the open PR for the current branch:

```bash
gh pr view --json number,url,title,state
gh pr diff
```

## CI

Check PR checks:

```bash
gh pr checks
```

List and inspect workflow runs:

```bash
gh run list --limit 10
gh run view <run-id>
gh run view <run-id> --log-failed
```

When debugging a failure: `gh pr checks` → `gh run list` → `gh run view` → `--log-failed`.

## API

For fields not exposed by subcommands:

```bash
gh api repos/owner/repo/pulls/55 --jq '.title, .state'
```

Most commands accept `--json` and `--jq` for structured output.
