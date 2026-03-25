# AGENTS.md
Repository-specific guidance for agentic coding tools working in
`/Users/tong/Development/k9s-multi-cluster-plugin`.

This file is intentionally based on the repository as it exists today.
Do not invent tooling, architecture, commands, or language-specific rules that are not present in the working tree.

## Current Repository State
- The repository now includes a Go module and implementation source files.
- The CLI is structured around Cobra commands.
- Tests exist for command behavior and generator logic.
- There is still no repository-specific lint configuration.
- No Cursor rules were found in `.cursor/rules/`, `.cursorrules`, or `.github/copilot-instructions.md`.

## Primary Rule
Work from observed repository facts.

- Inspect the current tree before assuming a stack.
- Prefer repository-local conventions over generic defaults.
- Update this file when real tooling or rules are added.
- Keep instructions accurate to the current branch contents.

## Product Intent
The main source of truth right now is `README.md`.

Based on `README.md` and the current Go implementation, the intended project direction is:
- a CLI for generating K9s plugin configuration
- multiple clusters supported from shared templates
- inputs include kubeconfig, template YAML, and override data
- output is generated K9s plugin YAML for the active cluster
- template placeholders are intended to use gomplate-style / Go-template-style syntax

Treat `README.md` as product intent and the Go/Cobra code as the current implementation direction.

Current template-direction guidance from `README.md`:
- keep templates close to normal K9s YAML
- use gomplate-style expressions for substitutions and defaults
- prefer template-defined defaults over hard-coding unmatched-cluster fallback in code

## Commands
### Build
Configured via Go module.

Agent guidance:
- Primary build command: `go build ./...`
- CLI entrypoint build: `go build -o ./bin/k9s-multi-cluster-plugin .`

### Lint
Not configured yet.

No detected repository-specific lint configuration such as golangci-lint.

Agent guidance:
- Do not claim lint passes unless a real lint tool exists and was run.
- Use `gofmt -w .` for Go formatting until a repository-specific lint tool is added.

### Test
Configured via Go test.

Agent guidance:
- Full test command: `go test ./...`
- Fixture-only, intentionally non-runnable test states are acceptable until a
  real feature implementation and test runner exist.
- Do not add fake failing tests or placeholder executables merely to simulate
  progress.

### Run A Single Test
Configured via Go test selection flags.

Example single-test commands:
- One package: `go test ./internal/generator`
- One named test: `go test ./internal/generator -run TestLoadActiveCluster`

### Install Or Bootstrap
Configured via Go modules.

- Install dependencies: `go mod download`
- Reproducible environment: open the repository in `.devcontainer/` and run
  build/test commands there when host Go tooling is unavailable.

## Code Style Guidelines
Follow the rules below, but let future repo config and established source files override these defaults.

### General Style
- Prefer clarity over cleverness.
- Keep changes minimal and tied to the requested task.
- Avoid speculative abstractions in an early-stage repository.
- Preserve existing wording and examples unless there is a reason to improve them.
- Keep documentation direct, concrete, and example-driven.

### Formatting
Observed conventions from `README.md`:
- Markdown uses ATX headings with concise titles.
- Prose is plain and direct.
- Examples are shown in fenced code blocks.
- YAML examples use 2-space indentation.
- YAML keys are lowercase unless an external schema requires otherwise.
- Template examples currently use gomplate-style / Go-template-style delimiters.

Agent guidance:
- Match the 2-space YAML indentation style in docs and examples.
- Keep Markdown headings short and content easy to scan.
- If editing template examples, preserve gomplate-style expression syntax unless the product direction changes.

### Imports
Go is now present, but there is no repository-specific import customization beyond standard Go conventions.

Agent guidance:
- Follow standard Go import grouping and formatting.
- If the repo later adds formatter or linter rules, those take precedence.
- Do not add custom import grouping rules without evidence in repo config.

### Types
Go is the current implementation language.

Agent guidance:
- Prefer explicit, readable types at public boundaries.
- Avoid unnecessary type indirection in early code.

### Naming
Observed convention:
- The repository name uses kebab-case: `k9s-multi-cluster-plugin`.

Agent guidance:
- Prefer descriptive names over short or clever names.
- Use names that match the surrounding language or file-format conventions.
- For documentation and standalone config files, kebab-case is a safe default.
- For YAML keys, prefer clear lowercase names unless a target schema dictates another form.

### Error Handling
The current Go implementation establishes these error-handling expectations:

Agent guidance:
- Surface actionable errors.
- Include enough context to identify which input failed.
- Avoid swallowing errors silently.
- Prefer deterministic failure over partially applied hidden behavior.
- Prefer returning errors with input path or resource context, as seen in the
  generator package.

### Comments And Docs
- Add comments only when intent is not obvious from the code itself.
- Keep comments factual and durable.
- Update docs when behavior changes.
- Prefer examples when documenting config formats or CLI behavior.

## Workflow Guidance
Before making changes:
- Inspect the repository tree.
- Read `README.md`.
- Check for newly added manifests and instruction files.
- Confirm whether this file also needs an update.

When adding a new stack or tool:
- Record the real build command.
- Record the real lint command.
- Record the real full-test command.
- Record at least one real single-test command.
- Add stack-specific style rules only if they are backed by config or code.
- Make integration fixtures runnable only when the corresponding feature and
  real test harness are implemented.

For feature work:
- Prefer using a separate git worktree for each feature so implementation work
  stays isolated from other in-progress changes.

When updating this file:
- Remove statements that are no longer true.
- Prefer exact commands over placeholders.
- Include Cursor or Copilot rules if those files are later added.
- Keep the file factual and repository-specific.

## Editor And Assistant Rules
At the time this file was written, `.cursor/rules/`, `.cursorrules`, and `.github/copilot-instructions.md` do not exist.
If any of those files are added later, fold their repository-specific guidance into this document so agents have a single source of truth.

## What Not To Do
- Do not assume the repo is Go, Node.js, Python, Rust, or any other stack.
- Do not create fictional commands in documentation.
- Do not claim formatting or linting standards that are not configured.
- Do not replace repository evidence with generic agent preferences.
