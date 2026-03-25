## Why

The repository has product intent for a multi-cluster K9s plugin generator, but it does not yet have a committed end-to-end integration test contract for the CLI behavior described in `README.md`. Defining that contract now gives future implementation work a concrete target and locks the decided no-match behavior that unmatched clusters can fall back to template-defined defaults rendered through gomplate-style template expressions and still succeed.

## What Changes

- Define a new OpenSpec capability for end-to-end CLI integration test coverage.
- Specify the canonical fixture layout for kubeconfig, template, overrides, and expected outputs.
- Require kubeconfig fixtures to represent realistic multi-cluster configs where `current-context` selects the active cluster under test.
- Record that template expressions and defaults are expressed with gomplate-style / Go-template-style syntax.
- Specify named success and failure test cases, including exit code, filesystem, stdout/stderr, and golden-file expectations, as contract artifacts.
- Record the unmatched-cluster rule as a success case that uses template-defined default values when available.
- Keep executable test-harness implementation out of scope for this change and hand it off to a follow-up change.

## Capabilities

### New Capabilities
- `cli-e2e-integration-tests`: Defines the end-to-end integration test contract, fixtures, expected outputs, and failure handling for the planned CLI.

### Modified Capabilities

None.

## Impact

- Affects future CLI implementation by defining the observable behavior it must satisfy.
- Adds OpenSpec change artifacts under `openspec/changes/define-e2e-integration-test-spec/`.
- Establishes a reusable fixture and test naming scheme for the eventual integration harness.
- Leaves runnable suite implementation to a separate follow-up change once a real stack exists.
- Makes active-cluster selection meaningful by ensuring kubeconfig fixtures contain multiple cluster definitions.
