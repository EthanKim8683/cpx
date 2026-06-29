# RFCs

Design and workflow decisions are RFCs — durable records to refer back to. In cpx, RFCs are implemented through issue comments on Linear (synced to GitHub).

Each RFC is an issue. The issue body is short — roughly a paragraph describing what needs to be figured out, not a full spec. The comment thread is the RFC itself.

## Workflow

### When an RFC is active

An RFC starts when the user asks to create one (e.g. "I want to create an RFC for …"). From that point, **mirror the entire conversation** on the issue's synced comment thread until the user ends the RFC — for example by starting implementation, opening a PR, or moving to a different RFC.

While an RFC is active, every exchange belongs on that thread. Do not wait until the end of a session. Do not skip exchanges that seem redundant or low-value — the user may refer back to them later, and anything left out of the thread is lost.

**Do not skip syncing.** Post for every exchange, even when approval slows things down. If it is not on the thread, it is not recorded.

### Synced comments

Each issue has a system comment noting that the thread is synced with GitHub. **Reply to that comment** — do not post a new top-level comment. Top-level comments on Linear do not sync to GitHub.

### Per exchange

When the user says something in conversation:

1. Post a summary as a reply on the synced thread.
2. Respond in conversation and post that exact response on the same synced thread.

Post **sequentially**, never in parallel:

1. Summary → wait for it to land
2. Response → wait for it to land
3. Next exchange

Parallel posts can arrive out of order on Linear and GitHub.

### Backfill

If syncing was missed, backfill before continuing. Walk the conversation chronologically and post each missing summary/response pair — still one pair at a time, never in parallel.

## Format

### Summaries

A summary is a short paragraph written from the user's perspective — as if they wrote it themselves. It captures what they said, nothing more.

- Paragraph, not bullet points
- No heading (e.g. no "Summary" title)
- First person ("I want…", not "The user wants…")
- Only what the user actually said — no additions or extrapolation

### Responses

Post a separate reply with your exact response — the same words you use in conversation, not a summary. Keep headings, lists, and code blocks as written.
