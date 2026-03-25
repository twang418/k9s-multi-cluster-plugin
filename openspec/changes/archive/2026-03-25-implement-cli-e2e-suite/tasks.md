## 1. Prerequisites

- [x] 1.1 Confirm the repository has a real CLI entrypoint or built binary
- [x] 1.2 Confirm the repository has a real test runner and documented test command

## 2. Success Cases

- [x] 2.1 Implement `renders debug plugin for active org1 cluster`
- [x] 2.2 Implement `renders debug plugin for active org3-bu1 cluster`
- [x] 2.3 Implement `renders template default image when active cluster has no matching override`
- [x] 2.4 Implement `writes unchanged template when no matching expression exists in the template`

## 3. Harness Verification

- [x] 3.1 Verify the harness renders output and uses a file diff for exact golden-file verification
- [x] 3.2 Verify the harness uses isolated temporary output paths per test

## 4. Failure Cases

- [x] 4.1 Implement failure tests for missing kubeconfig path, malformed kubeconfig, malformed template YAML, and malformed overrides YAML
- [x] 4.2 Implement failure tests for missing current-context and missing cluster reference resolution
- [x] 4.3 Implement failure tests for unsupported match type, missing replacement values, and unwritable output paths
- [x] 4.4 Verify each failure case asserts non-zero exit status, absent or empty output, and stderr substring checks with actionable error context

## 5. Deferred Precedence Follow-Up

- [ ] 5.1 Add the overlapping-rule precedence integration test only after the documented precedence rule is accepted and reflected in the spec
