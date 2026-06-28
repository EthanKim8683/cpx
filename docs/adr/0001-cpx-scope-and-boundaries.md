# ADR-0001: CPX scope and boundaries

- **Status:** Accepted
- **Date:** 2026-06-28
- **Related:** [CP-5](https://linear.app/ethankim8683/issue/CP-5/cpx-vision-goals-and-feature-discussion), [GitHub #1](https://github.com/EthanKim8683/cpx/issues/1)

## Context

CPX is a personal tool for competitive programming. It exists to reduce repetitive manual work in CP workflows and to support a practice style where more time goes to thinking and less to mechanical implementation.

Several capability areas have been identified — bundling, scraping, scaffolding, and AI-assisted practice — but the project needs a single, explicit statement of what CPX is in scope for, what is out of scope, and how those areas relate. Without that, later technical ADRs lack a shared foundation.

## Decision

### What CPX is

CPX is a CLI-oriented tool (and supporting library code) that speeds up common competitive programming workflows: preparing submissions, fetching problem/contest material, setting up workspace structure, and — during **practice only** — integrating AI to compress implementation time so thinking gets more reps.

### Goals

1. Reduce repetitive manual labor in competitive programming.
2. Speed up workflows around solving problems and participating in contests.
3. Use AI tooling **during practice** to shift time from implementation toward thinking, so high-level understanding and decision-making can move faster.

### Practice vs. contests — hard boundary

**AI tools are for practice only. Using them during contests is cheating and is explicitly out of scope.**

Training happens outside of contests — Project Euler, topic study, ICPC prep, upsolving — so that thinking skills transfer to contest conditions where work is unaided. CPX must make this boundary obvious and enforceable (e.g. practice mode vs. contest mode).

### Core problem being addressed

The tool is not only about saving keystrokes. The recurring pattern in practice is:

- Disproportionate time spent **implementing** vs. **thinking**
- A "thinking brain" that lags because it is under-trained at speed
- ICPC study sessions that drift into unproductive stretches with little real thinking

**Hypothesis:** compressing implementation time during practice yields more cycles for navigation — seeing the big picture, making high-level decisions, and moving through problem space faster. That training should show up in contests **without AI assistance there**.

This shows up most clearly in:

- **Project Euler / math-heavy problems** — weakness in high-level decision-making and problem navigation, not only missing knowledge
- **Topic study and ICPC prep** — getting stuck in implementation loops instead of the thinking loop

### In-scope feature areas

These are planned directions. Priority and detailed design are open; each area will get its own technical ADR(s) when concrete implementation decisions are made.

| Area | Description |
|------|-------------|
| **Bundling** | Turn a source file and its dependencies into a single source file suitable for submission on online judges (Codeforces, AtCoder, etc.). |
| **Scraping** | Scrape problems and contests; submit solutions to online judges. |
| **Scaffolding** | Create directory structure for solving individual problems or participating in contests. |
| **AI-assisted practice** *(emerging)* | Integrate AI into the **practice** loop only: fast enough implementation that thinking does not stall; more reps on high-level understanding and navigation; clear separation from contest workflows. |

### Out of scope

- AI assistance during live contests or any context where it would constitute cheating
- General-purpose IDE or editor replacement
- Features unrelated to competitive programming workflow

### Documentation split

- **[CP-5 / GitHub #1](https://github.com/EthanKim8683/cpx/issues/1)** — living discussion: vision, priorities, open questions
- **ADR-0001 (this document)** — settled product scope and boundaries
- **Future ADRs (0002+)** — one decision per technical choice within each feature area

## Consequences

### Positive

- Clear product boundary for all future work
- Explicit ethical line on AI use (practice only)
- Feature areas can evolve independently via follow-on ADRs without re-litigating scope
- CP-5 remains the place for ongoing discussion; this ADR stays stable unless scope materially changes

### Negative / open

- Priority order across bundling, scraping, scaffolding, and AI is not yet decided
- Which online judges to support first is not yet decided
- CLI vs. library surface area is not yet decided
- Concrete AI integration (when to invoke, context, enforcement of practice mode) requires future ADR(s)
- How to measure thinking-speed improvement is not yet defined

## References

- [CP-5: CPX vision, goals, and feature discussion](https://linear.app/ethankim8683/issue/CP-5/cpx-vision-goals-and-feature-discussion)
- [GitHub #1](https://github.com/EthanKim8683/cpx/issues/1)
