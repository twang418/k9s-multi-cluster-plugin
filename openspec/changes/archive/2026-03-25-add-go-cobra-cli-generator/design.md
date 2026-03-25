## Context

The repository currently contains product intent, OpenSpec contract artifacts, and canonical fixtures, but no implementation stack. This change introduces the first real application code: a Go CLI built with Cobra that reads a kubeconfig, loads template and override inputs from folders on disk, and generates K9s plugin YAML for the active cluster.

The generator implementation also now depends on a separate devcontainer prerequisite change so final build and test verification can happen in a reproducible Go environment instead of relying on host-machine setup.

The existing docs already establish important constraints:
- template rendering should follow gomplate-style / Go-template-style expressions
- unmatched clusters should rely on template-defined defaults rather than hidden code fallbacks
- output should remain close to normal K9s plugin YAML
- output location must be configurable so future tests can write into temporary directories

## Goals / Non-Goals

**Goals:**
- Add a Go module and Cobra-based CLI entrypoint.
- Implement a command that accepts template folder, override folder, kubeconfig file, and output path overrides.
- Resolve the active cluster from kubeconfig using `current-context`.
- Load override data, select the matching cluster rule, and render gomplate-style template expressions into final YAML.
- Produce output shaped for K9s plugin configuration and suitable for future fixture-backed tests.

**Non-Goals:**
- Implement the full runnable end-to-end suite in this change.
- Support every possible gomplate function; only the subset needed for documented template behavior is required initially.
- Solve overlapping-rule precedence beyond the already deferred requirement.
- Add multiple subcommands or a large CLI surface beyond the minimum generation workflow.

## Decisions

- Use Go with Cobra for the initial CLI framework.
  - Rationale: the user explicitly requested Go and Cobra, and Cobra provides conventional command parsing, help text, and future extensibility.
  - Alternative considered: implementing a minimal `flag`-based CLI. Rejected because Cobra better fits an intended multi-command CLI without much extra complexity.

- Start with a single generation command and explicit input flags.
  - Rationale: this keeps the first feature slice small while making input/output paths testable and deterministic.
  - Alternative considered: implicit default paths only. Rejected because configurable paths are important for temp-dir based testing and local experimentation.

- Treat the template folder and override folder as directory inputs, even if the initial implementation only reads one file from each.
  - Rationale: this matches the requested interface and leaves room for future expansion to multiple templates or layered overrides.
  - Alternative considered: only accepting one template file and one override file. Rejected because it would diverge from the requested folder-based workflow.

- Use Go-template-compatible rendering with a constrained data model keyed by values like `.image`.
  - Rationale: this aligns with the documented gomplate-style direction while allowing a minimal first implementation.
  - Alternative considered: inventing a custom placeholder parser. Rejected because the repo has already chosen the gomplate-style direction.

- Resolve unmatched clusters by rendering the template with no override values beyond the template's own defaults.
  - Rationale: this preserves template-owned defaults and avoids hidden fallback logic in code.
  - Alternative considered: injecting code-defined default values. Rejected because it would compete with the template as the source of truth.

- Write generated output to an explicit file path that can be overridden.
  - Rationale: K9s expects plugin YAML, but tests and local workflows need a configurable destination.
  - Alternative considered: writing only to stdout. Rejected because the requested behavior is to generate a plugin file and because file output fits the existing fixture model better.

## Risks / Trade-offs

- [Template rendering grows beyond the initial subset] -> Start with only the gomplate-style features needed by the documented fixtures and expand later with explicit spec updates.
- [Folder semantics are underspecified] -> Define a simple first-pass convention in implementation, document it, and keep the command flags explicit.
- [K9s plugin YAML validity is assumed rather than deeply validated] -> Preserve the template structure and focus initial validation on readable input errors plus deterministic output generation.
- [Kubeconfig parsing edge cases become broader than the first feature slice] -> Implement the documented active-cluster path first and return actionable errors for unsupported or unresolved references.
- [Future tests depend on exact output formatting] -> Keep generated YAML stable and document any formatting assumptions when the real test harness is added.

## Migration Plan

- Add the Go module and Cobra dependency.
- Implement the CLI entrypoint and generation command.
- Add input loading and rendering packages.
- Update `AGENTS.md` with real build/test commands once the stack exists.
- Add the devcontainer prerequisite and use it for final generator verification.
- Keep the separate executable-suite OpenSpec change as the follow-up for runnable integration tests.

## Open Questions

- What exact file discovery convention should the template folder and override folder use for the first implementation?
- Should the first command generate one combined plugin file from one selected template, or support multiple template files in a folder immediately?
- How much validation against K9s plugin schema should the first feature slice enforce beyond well-formed YAML generation?
