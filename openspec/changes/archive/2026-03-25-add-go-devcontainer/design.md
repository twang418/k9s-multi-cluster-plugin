## Context

The repository has begun moving from documentation-only planning into a real Go/Cobra implementation, but the current shell environment still does not guarantee that `go` is installed. That makes build and test verification unreliable and weakens the meaning of “feature complete.”

The generator change already assumes a Go-based workflow and updated `AGENTS.md` with Go build/test commands. A devcontainer should become the reproducible baseline so these commands run consistently for maintainers, contributors, and coding agents.

## Goals / Non-Goals

**Goals:**
- Add a reproducible devcontainer with a working Go toolchain for this repository.
- Make it possible to run `go mod download`, `go build ./...`, and `go test ./...` inside the container.
- Keep the environment minimal but sufficient for Go CLI development and verification.
- Establish devcontainer-based build/test verification as the completion path for the generator change and future Go feature work.

**Non-Goals:**
- Implement new application features unrelated to environment setup.
- Add a full CI pipeline in this change.
- Add every optional developer convenience tool on the first pass.
- Replace later feature-specific testing changes such as the executable end-to-end suite.

## Decisions

- Use a repository-local `.devcontainer/` configuration as the first reproducible environment.
  - Rationale: this is the most direct way to make the repo self-describing for local development and agent execution.
  - Alternative considered: documenting host-machine prerequisites only. Rejected because that keeps setup machine-dependent.

- Pin a Go-capable container image rather than relying on host Go installation.
  - Rationale: the immediate problem is missing or inconsistent Go availability.
  - Alternative considered: downloading Go manually in scripts. Rejected because the container should already contain the toolchain.

- Keep the first devcontainer minimal and focused on build/test verification.
  - Rationale: the main need is deterministic execution of Go commands, not a heavily customized IDE environment.
  - Alternative considered: loading many extra tools and extensions immediately. Rejected because it adds noise before the core workflow is stable.

- Make the generator change depend on devcontainer-based verification for final completion.
  - Rationale: this aligns the user’s new requirement with the actual definition of done.
  - Alternative considered: leaving the generator change marked complete without reproducible verification. Rejected because it weakens trust in the result.

## Risks / Trade-offs

- [Devcontainer image choice becomes stale] -> Pin a reasonable Go baseline and revisit version updates as part of normal maintenance.
- [Container startup feels heavier than host-native development] -> Keep the first container minimal and focused on required tooling.
- [Feature work remains marked complete too early] -> Reopen generator verification tasks until build/test passes in the devcontainer.
- [Developers want extra tools immediately] -> Add only what is needed now and expand later when concrete needs appear.

## Migration Plan

- Add `.devcontainer/` configuration.
- Document how to enter the container and run build/test commands.
- Reopen the generator change’s verification step so it depends on successful in-container execution.
- Use the devcontainer as the baseline environment for subsequent Go feature work.

## Open Questions

- Which exact Go version should be pinned for the first container image?
- Should the initial container include only Go tooling, or also common YAML/debug helpers?
- Do we want a `postCreateCommand` immediately, or keep setup commands explicit at first?
