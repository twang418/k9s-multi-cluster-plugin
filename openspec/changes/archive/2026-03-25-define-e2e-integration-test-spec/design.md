## Context

The repository currently has product intent in `README.md` and a detailed draft in `integration-test-spec.md`, but no OpenSpec change artifacts that make the CLI's end-to-end integration contract apply-ready. The project has no implementation stack yet, so this change needs to stay implementation-agnostic while still being concrete enough to guide future fixture creation and test harness work.

The kubeconfig inputs for this CLI are expected to be realistic multi-cluster files,
not trivial single-cluster fixtures. The integration contract should therefore prove
that the CLI selects the active cluster via `current-context` from a kubeconfig that
contains multiple cluster and context definitions.

The template direction has also shifted from a custom substitution form to a
gomplate-style / Go-template-style expression model so defaults and future richer
template behavior can live directly in the template document.

## Goals / Non-Goals

**Goals:**
- Capture the CLI's end-to-end integration behavior in OpenSpec artifacts.
- Define stable fixture names, test case names, expected outputs, and failure cases.
- Ensure kubeconfig fixtures model multiple available clusters so active-cluster selection is part of the contract.
- Lock the decided behavior that unmatched clusters use template-defined default values when available and still exit successfully.
- Align the integration contract with gomplate-style template expressions and defaults.
- Leave runnable suite implementation to a separate follow-up change once a concrete test stack exists.

**Non-Goals:**
- Choosing a programming language, CLI framework, or test runner.
- Defining internal parser, renderer, or data structure design.
- Implementing executable end-to-end tests in this repository before a real CLI stack exists.
- Locking override precedence beyond the already noted deferred case.

## Decisions

- Keep the specification black-box and CLI-facing. This preserves value even if the eventual implementation language, binary name, or subcommand structure changes.
- Use a canonical `testdata/` fixture layout in the spec. This gives future implementation work deterministic paths and golden files without assuming any specific test framework.
- Make kubeconfig fixtures multi-cluster by default. This prevents the tests from passing with overly trivial inputs and verifies that `current-context` is the mechanism that selects the cluster to render for.
- Express substitutions and defaults with gomplate-style template syntax, using data fields such as `.image` and defaults such as `default "busybox"`. This keeps template behavior close to an established model and leaves room for future richer expressions without inventing a custom mini-language.
- Do not invent a fake CLI or failing test harness just to make the suite executable early. Until a real project stack exists, the fixtures and OpenSpec artifacts are the intentional non-runnable test state.
- Track executable harness work in a separate follow-up change so this change can complete as a contract-and-fixtures deliverable.
- Keep one stable invocation shape in the spec, but allow the harness adapter to change later if the real CLI shape differs. This keeps fixture and assertion contracts stable while avoiding premature lock-in to a command surface that does not exist yet.
- Treat unmatched clusters as successful render cases when the template defines defaults through gomplate-style expressions such as `{{ .image | default "busybox" }}`. This keeps environment-specific fallback behavior in the template itself and allows the final generated YAML to contain a concrete default value instead of an unresolved placeholder.
- Treat malformed files, unresolved kubeconfig references, unsupported match types, missing replacement values, and unwritable output paths as hard failures with non-zero exit status and no meaningful output file. These are deterministic operator errors and need actionable integration coverage.

## Risks / Trade-offs

- [Spec drift from future CLI flags] -> Keep the invocation shape isolated as a harness concern and preserve fixture names plus assertions even if flags change.
- [Fixtures become unrealistically simple] -> Keep kubeconfig fixtures multi-cluster so cluster selection logic is exercised instead of bypassed.
- [Template language becomes too open-ended] -> Document the supported gomplate-style subset clearly and prefer a constrained, validated data model over arbitrary template power.
- [Placeholder tests create false confidence] -> Prefer an explicitly non-runnable documented state over adding fake failing tests that are not tied to a real CLI entrypoint.
- [Spec becomes too implementation-specific] -> Express requirements in terms of process behavior, files, and outputs rather than language-specific helpers.
- [Undocumented precedence remains ambiguous] -> Keep overlapping-rule coverage explicitly deferred until precedence is documented, instead of guessing and freezing the wrong behavior.
- [Duplicate source-of-truth documents] -> Use this change as the OpenSpec source for implementation readiness; `integration-test-spec.md` can remain as supporting context until implementation begins.
