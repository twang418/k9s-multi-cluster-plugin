## Context

The repository already defines the integration-test contract and fixture set for
the planned CLI, but those artifacts are intentionally non-runnable because no
implementation stack exists yet. This follow-up change covers the point where a
real CLI entrypoint and test runner have been introduced and the documented
cases can become executable.

## Goals / Non-Goals

**Goals:**
- Implement a real black-box integration harness against the actual CLI.
- Enable fixture-backed tests only for behaviors that exist in the real
  implementation.
- Reuse the canonical `testdata/` fixtures and expected outputs.
- Preserve the documented success, no-op, and failure assertions from the
  contract change.

**Non-Goals:**
- Redefining fixture semantics already documented in the contract change.
- Introducing fake failing tests or placeholder executables.
- Locking overlapping-rule precedence before it is explicitly documented.

## Decisions

- Treat `define-e2e-integration-test-spec` as the source of truth for fixture
  names and behavioral expectations.
- Implement the suite as a real external-process harness, not by calling
  internal helpers directly.
- Enable tests incrementally as the corresponding CLI features land.
- Verify successful renders by diffing the generated output file against the
  expected golden file so mismatches surface as a real file diff.
- Create a dedicated test workspace per case and auto-clean it by default,
  with an explicit test flag to preserve rendered outputs when debugging.
- Keep a gitignored `output/` directory available for manual CLI renders when a
  developer wants to inspect generated YAML outside the test harness.
- Keep failure assertions focused on non-zero exit status, absent or empty
  output, and actionable stderr substrings.

## Risks / Trade-offs

- [Implementation arrives in slices] -> Enable tests feature by feature instead
  of requiring the whole matrix at once.
- [Harness drifts from contract] -> Reuse the existing fixture and golden-file
  names without renaming them in the follow-up change.
- [Fake progress pressure] -> Do not mark runnable-suite tasks complete until a
  real CLI entrypoint and test runner exist.
