# GitHub CLI Reference

All commands assume the `gh` CLI is authenticated and run from within a repository directory.

## Identify User
```bash
gh api user -q .login
```
*If this fails, notify the user (authentication is required to manage issues).*

## Label Setup (One-time or if creation fails)
If creating the session issue fails because labels do not exist, create them first:
```bash
gh label create "session" --description "Work session tracker" --color "0E8A16"
gh label create "user:$USERNAME" --description "Sessions for $USERNAME" --color "5319E7"
```

## Search Sessions
```bash
gh issue list --label "session" --state open --limit 5
gh issue list --label "session" --label "user:$USERNAME" --state open --limit 5
```

## Create Session Issue
```bash
gh issue create \
  --label "session" \
  --label "user:$USERNAME" \
  --title "Session: $TITLE" \
  --body "$BODY"
```

## Update Session Summary
```bash
gh issue edit $NUMBER --body "$UPDATED_BODY"
```

## View Session
```bash
gh issue view $NUMBER --json body -q .body
```

## Log Event (Comment)
```bash
gh issue comment $NUMBER --body "$EVENT_DESCRIPTION"
```