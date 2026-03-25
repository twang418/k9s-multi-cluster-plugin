## Why

The repository now has a documented end-to-end integration contract and canonical
fixtures under `testdata/`, but it intentionally does not yet have a runnable
integration suite. Once a real CLI entrypoint and test runner exist, that work
should be tracked separately so the suite can be enabled feature by feature.

## What Changes

- Implement a real black-box integration test harness against the actual CLI
  entrypoint or built binary.
- Enable the documented success and failure cases from the existing fixture set.
- Keep tests feature-gated so each fixture-backed test is enabled only when the
  corresponding CLI behavior exists.
- Preserve exact golden-file comparison for successful renders and actionable
  stderr assertions for failures.

## Capabilities

### New Capabilities
- `cli-e2e-suite-execution`: Defines how the documented integration contract
  becomes a runnable suite once a real CLI stack exists.

### Modified Capabilities

None.

## Impact

- Depends on a future real CLI implementation and a real test runner.
- Builds directly on the contract and fixtures from
  `define-e2e-integration-test-spec`.
- Keeps runnable integration work explicitly separate from fixture-only spec
  definition.
