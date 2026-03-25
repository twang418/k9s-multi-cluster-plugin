# Testdata

This directory contains canonical fixture inputs and expected outputs for the
planned end-to-end integration suite.

Current status:

- the fixtures are intentionally non-runnable today
- the repository now contains a real CLI entrypoint and Go test coverage
- the full fixture-backed end-to-end suite is still separate follow-up work

These files are still useful because they define the contract that a future
implementation must satisfy.

Repository rule:

- do not add fake failing tests, placeholder binaries, or invented harness code
  just to make the suite appear executable
- wait until a real project stack exists, then wire these fixtures into the real
  integration tests
- enable each fixture-backed test only when the corresponding feature behavior
  exists in the real implementation

Fixture groups:

- `kubeconfig/`: multi-cluster kubeconfig inputs that select an active cluster by
  `current-context`
- `template/`: K9s plugin templates using gomplate-style expressions
- `template-single/`: single-file template folder fixture for the first CLI slice
- `overrides/`: cluster matching and replacement data
- `overrides-single/`: single-file override folder fixture for the first CLI slice
- `expected/`: golden outputs for successful render cases
