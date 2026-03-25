# Testdata

This directory contains canonical fixture inputs and expected outputs for the
end-to-end integration suite.

Current status:

- the repository contains a real CLI entrypoint and Go test coverage
- the fixture-backed end-to-end suite now executes the real CLI from `go test`
- overlapping override precedence remains a documented deferred case

These files define the contract and golden outputs the end-to-end harness uses.

Repository rule:

- keep the suite black-box by invoking only the real CLI entrypoint or built
  binary
- use exact golden-file comparison for successful renders
- enable overlapping-rule precedence coverage only after the precedence rule is
  documented
- use `go test ./e2e -args -e2e-auto-clean=false` when you want the rendered
  output workspace to remain on disk for inspection

Fixture groups:

- `kubeconfig/`: multi-cluster kubeconfig inputs that select an active cluster by
  `current-context`
- `template/`: K9s plugin templates using gomplate-style expressions
- `template-single/`: single-file template folder fixture for the first CLI slice
- `overrides/`: cluster matching and replacement data
- `overrides-single/`: single-file override folder fixture for the first CLI slice
- `expected/`: golden outputs for successful render cases
