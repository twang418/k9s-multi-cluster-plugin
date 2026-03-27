## 1. Command and request shape

- [ ] 1.1 Update `cmd/generate.go` to support a K9s install workflow alongside the existing explicit output-path workflow.
- [ ] 1.2 Extend generator request and result types to carry the selected output mode and any resolved install-root information needed for reporting.

## 2. Kubeconfig target resolution and merge behavior

- [ ] 2.1 Extend kubeconfig loading to collect every context name that references the active cluster selected by `current-context`.
- [ ] 2.2 Implement K9s context-path resolution under `$XDG_DATA_HOME/k9s/clusters/<cluster>/<context>/plugins.yaml`, including the documented default when `XDG_DATA_HOME` is unset.
- [ ] 2.3 Implement plugin-map merge logic that preserves unrelated plugins and replaces generated plugin keys in existing destination files.
- [ ] 2.4 Write generated plugin output to each matching context path and return actionable errors for malformed destination YAML or write failures.

## 3. Test coverage

- [ ] 3.1 Add generator tests for multi-context cluster resolution, K9s path resolution, and plugin-key merge behavior.
- [ ] 3.2 Add command tests for the explicit output-path workflow and the K9s install workflow.
- [ ] 3.3 Add e2e fixtures and cases that verify one active cluster writes merged `plugins.yaml` files for multiple contexts without removing unrelated plugins.

## 4. Documentation and verification

- [ ] 4.1 Update `README.md` with the K9s plugin discovery paths, the new install workflow, and merge semantics for existing `plugins.yaml` files.
- [ ] 4.2 Run the relevant Go test suites in the devcontainer and confirm the new workflow is covered before marking the change complete.
