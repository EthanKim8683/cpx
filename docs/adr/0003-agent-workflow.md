# ADR-0003: Agent Workflow

- **Status:** Accepted
- **Date:** 2026-06-28
- **Related:** [CP-10: Agent workflow guidelines](https://linear.app/ethankim8683/issue/CP-10/agent-workflow-guidelines), [GitHub #10](https://github.com/EthanKim8683/cpx/issues/10)

## Context

AI agents can be helpful when working on cpx, but there were no shared guidelines or workflow for how agents should operate in this repo. Without a record of decisions, agent work is hard to align on and refer back to.

## Decision

### Issues

Issues are short — roughly a paragraph. They describe what needs to be figured out, not a full spec. Details belong in issue comments, not the issue body.

### RFCs through issues

Design and workflow decisions are worked out through issue comments. The comment thread is the record to refer back to.

When using this workflow:

1. The user states something in conversation.
2. The agent posts a summary as a reply on the GitHub-synced comment thread.
3. The agent responds in conversation and posts that response on the same synced thread.

### Synced comments

Each issue has a system comment noting that the thread is synced with GitHub. **Reply to that comment** — do not post a new top-level comment. Top-level comments on Linear do not sync to GitHub.

### Summaries

A summary is a short paragraph written from the user's perspective — as if they wrote it themselves. It captures what they said, nothing more.

- Paragraph, not bullet points
- No heading (e.g. no "Summary" title)
- First person ("I want…", not "The user wants…")
- Only what the user actually said — no additions or extrapolation

### Responses

After posting the summary, post a separate reply with the agent's response — acknowledgment, questions, or pushback.

## Alternatives

- **Long issue descriptions as specs** — rejected; issues stay short, details live in comments.
- **Top-level Linear comments** — rejected; they do not sync to GitHub.
- **Bullet-point summaries with headings** — rejected; summaries are a paragraph in the user's voice.
- **RFC decisions only in chat** — rejected; issue comments are the shared record.

## Consequences

- **Linear–GitHub sync dependency** — the RFC workflow requires replying to the synced comment thread.
- **Issue comments as canonical record** — refer to the synced thread for decisions, not chat history alone.
- **Open topics** — what to read first, code conventions, and commit/PR norms are not covered here yet.

## References

- [CP-10: Agent workflow guidelines](https://linear.app/ethankim8683/issue/CP-10/agent-workflow-guidelines)
- [GitHub #10](https://github.com/EthanKim8683/cpx/issues/10)
