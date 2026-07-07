---
name: saving-sessions
description: Initialises and maintains a living issue summary and event log comments for the session.
version: 1.8.0
languages: all
---

# Saving Sessions

## Overview
All conversation sessions must be persisted in a tracking issue (e.g., GitHub issue).
1. **Issue Body**: A living, high-level summary of the overall workspace state, configurations, and final decisions. It synthesizes the progress made and must be kept up-to-date as the single source of truth.
2. **Issue Comments**: An append-only log of significant events, containing detailed context, findings, error traces, and step-by-step outcomes as they occur.
3. **Dual Update**: When a significant event occurs, immediately post a detailed comment, then rewrite the issue body summary to synthesize the new state.

## Initialization (Run First)
At the very start of the conversation, initialize the session tracker:
1. **Identify Username**: Retrieve your username (`gh api user -q .login`). If authentication fails, notify the user.
2. **Create**: Create the session issue using labels `session` and `user:{username}` (prefix title with `Session: `). If creation fails because the labels do not exist, create them first (see [github.md](references/github.md)). Do this before modifying any code.

## Defining "Significant" Events
A significant event is when **something new or unexpected occurs** (especially anything contrary to the initial plan). Examples (non-exhaustive):
- **Unexpected Discoveries**: System behavior, undocumented rules, or API constraints.
- **Architectural Decisions**: Decisions made (e.g., choosing X over Y, changing plans).
- **Blockers**: Environmental, permission, or API limitations halting progress.
- **Milestones**: Critical accomplishments (e.g., tests passing, feature complete).

## The Sync Loop
Immediately upon encountering a significant event, execute these two steps:
1. **Comment**: Post the event description as a new comment on the session issue (specific, objective, no timestamps).
2. **Rewrite Summary**: Read the current session body, update the living summary to reflect the new state, and edit the body.

## Sub-Issues
Sub-issues are regular issues created during the session to track independent work. Link any created sub-issues inside the parent session issue's `## Sub-Issues` section as a reference to establish context.

## Quick Reference
- **GitHub CLI Commands**: See [github.md](references/github.md)
- **Formatting Templates**: See [templates.md](references/templates.md)

## Common Mistakes
- **Mistake**: Making workspace changes before initializing the session tracker.
  - *Fix*: Make username lookup and issue creation your absolute first steps.
- **Mistake**: Batching comments or delaying body updates.
  - *Fix*: Run the comment + rewrite loop immediately as each significant event occurs.
- **Mistake**: Logging minor step-by-step actions instead of conceptual events.
  - *Fix*: Log only new/unexpected findings, decisions, or milestones.
- **Mistake**: Writing incomplete summaries assuming the conversation will close later.
  - *Fix*: Treat the current moment as the final state. Make the summary complete and self-contained.