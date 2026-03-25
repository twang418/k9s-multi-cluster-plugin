## 1. Prerequisites

- [ ] 1.1 Confirm the repository has a real CLI entrypoint or built binary
- [ ] 1.2 Confirm the repository has a real test runner and documented test command

## 2. Success Cases

- [ ] 2.1 Implement `renders debug plugin for active org1 cluster`
- [ ] 2.2 Implement `renders debug plugin for active org3-bu1 cluster`
- [ ] 2.3 Implement `renders template default image when active cluster has no matching override`
- [ ] 2.4 Implement `writes unchanged template when no matching expression exists in the template`

## 3. Harness Verification

- [ ] 3.1 Verify the harness uses exact golden-file comparison for success outputs
- [ ] 3.2 Verify the harness uses isolated temporary output paths per test

## 4. Failure Cases

- [ ] 4.1 Implement failure tests for missing kubeconfig path, malformed kubeconfig, malformed template YAML, and malformed overrides YAML
- [ ] 4.2 Implement failure tests for missing current-context and missing cluster reference resolution
- [ ] 4.3 Implement failure tests for unsupported match type, missing replacement values, and unwritable output paths
- [ ] 4.4 Verify each failure case asserts non-zero exit status, absent or empty output, and stderr substring checks with actionable error context

## 5. Deferred Precedence Follow-Up

- [ ] 5.1 Add the overlapping-rule precedence integration test only after the documented precedence rule is accepted and reflected in the spec
