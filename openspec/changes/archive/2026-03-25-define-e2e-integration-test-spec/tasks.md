## 1. Canonical fixtures

- [x] 1.1 Create the `testdata/kubeconfig/`, `testdata/template/`, `testdata/overrides/`, and `testdata/expected/` fixture layout defined by the spec
- [x] 1.2 Add canonical success fixtures for `active-org1.yaml`, `active-org3-bu1.yaml`, `active-unmatched.yaml`, `debug-template.yaml`, `standard-overrides.yaml`, `rendered-org1.yaml`, `rendered-org3-bu1.yaml`, `unchanged-unmatched.yaml`, and `unchanged-no-placeholder.yaml`, with kubeconfig fixtures containing multiple cluster and context entries
- [x] 1.3 Add canonical failure fixtures for malformed kubeconfig, missing current-context, missing cluster-for-context, malformed template, malformed overrides, invalid match type overrides, missing replace value overrides, and overlapping overrides

## 2. Contract And Documentation Alignment

- [x] 2.1 Align the spec, README, and fixture set around gomplate-style template rendering and template-defined defaults
- [x] 2.2 Document that fixture-backed tests remain intentionally non-runnable until a real CLI entrypoint and test runner exist
- [x] 2.3 Remove executable test implementation from this change and hand it off to a follow-up OpenSpec change

## 3. Deferred Follow-Up

- [x] 3.1 Record `applies documented precedence when multiple override rules match` as deferred until precedence behavior is explicitly documented
- [x] 3.2 Open a follow-up OpenSpec change for executable integration-test implementation work
