# cpx

A competitive programming utility that handles bundling files for submission to online judges, scraping problems and contests, automatic submission, and scaffolding problem/contest directories.

## Setup

1. Install [direnv](https://direnv.net/) and [Task](https://taskfile.dev/).

2. Clone and allow direnv:
   ```sh
   git clone https://github.com/EthanKim8683/cpx.git
   cd cpx
   direnv allow
   ```

3. Run an AI agent on the repo. The root `AGENTS.md` will guide it through
   setting up environment variables and generating compiler configs.

4. Verify the build:
   ```sh
   go build ./...
   go test -race ./...
   ```

## Development

See [`docs/README.md`](docs/README.md) for architecture, conventions, and how to contribute.
