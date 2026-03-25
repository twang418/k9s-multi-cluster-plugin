## ADDED Requirements

### Requirement: Canonical integration fixtures are defined
The project SHALL define a canonical end-to-end integration fixture layout for the planned CLI under `testdata/`, including kubeconfig fixtures, template fixtures, override fixtures, and expected output fixtures, so the same files can be reused across implementations.

#### Scenario: Required fixture set is documented
- **WHEN** the integration test specification is prepared for implementation
- **THEN** it MUST name fixture files for active org1, active org3-bu1, active unmatched, malformed kubeconfig, missing current-context, missing cluster-for-context, debug template, malformed template, missing-placeholder template, standard overrides, malformed overrides, invalid match type overrides, missing replace value overrides, overlapping overrides, rendered org1 output, rendered org3-bu1 output, unchanged unmatched output, and unchanged no-placeholder output

#### Scenario: Kubeconfig fixtures represent multiple available clusters
- **WHEN** a kubeconfig fixture is used by the integration suite
- **THEN** it MUST contain multiple cluster and context entries, and `current-context` MUST determine which cluster is treated as active for that test case

### Requirement: Integration tests cover successful rendering flows
The integration specification SHALL define named success cases that execute the real CLI process against files on disk and verify exit status, output file creation, and exact output contents for matched-cluster rendering.

#### Scenario: Template expression defines a default value
- **WHEN** a template contains a gomplate-style expression like `{{ .image | default "busybox" }}`
- **THEN** the CLI MUST treat `busybox` as the fallback value for `image` when no matched override provides that replacement

#### Scenario: Org1 regex match renders expected output
- **WHEN** the CLI runs with `active-org1.yaml`, `debug-template.yaml`, `standard-overrides.yaml`, and a writable output path
- **THEN** the process MUST exit with code `0`, create the output file, render the image template expression with the org1 image value, and produce content that exactly matches `rendered-org1.yaml`

#### Scenario: Org3 list match renders expected output
- **WHEN** the CLI runs with `active-org3-bu1.yaml`, `debug-template.yaml`, `standard-overrides.yaml`, and a writable output path
- **THEN** the process MUST exit with code `0`, create the output file, render the image template expression with the org3-bu1 image value, and produce content that exactly matches `rendered-org3-bu1.yaml`

### Requirement: Defaulted output is treated as a successful result when replacement is unnecessary
The CLI SHALL succeed when the active cluster matches no override rule or when the template contains no matching gomplate expression to render for that value.

#### Scenario: No override rule matches active cluster
- **WHEN** the CLI runs with `active-unmatched.yaml`, `debug-template.yaml`, `standard-overrides.yaml`, and a writable output path
- **THEN** the process MUST exit with code `0`, create the output file, render the template-defined default image value, and produce content that exactly matches `unchanged-unmatched.yaml`

#### Scenario: Template contains no matching expression to render
- **WHEN** the CLI runs with `active-org1.yaml`, `missing-placeholder-template.yaml`, `standard-overrides.yaml`, and a writable output path
- **THEN** the process MUST exit with code `0`, create the output file, and produce content that exactly matches `unchanged-no-placeholder.yaml`, with no dependency on K9s plugin input variables that are not defined by that template

### Requirement: Integration tests cover deterministic input failures
The integration specification SHALL define named failure cases for invalid input files and invalid override configuration, and each case MUST assert non-zero exit status, absent or empty output, and actionable stderr content.

#### Scenario: Kubeconfig file is missing or malformed
- **WHEN** the CLI runs with a non-existent kubeconfig path or `malformed.yaml`
- **THEN** the process MUST exit non-zero, MUST NOT produce a meaningful output file, and stderr MUST identify kubeconfig loading or parsing failure

#### Scenario: Active kubeconfig references cannot be resolved
- **WHEN** the CLI runs with `missing-current-context.yaml` or `missing-cluster-for-context.yaml`
- **THEN** the process MUST exit non-zero, MUST NOT produce a meaningful output file, and stderr MUST identify the unresolved context or cluster reference

#### Scenario: Template or overrides yaml is malformed
- **WHEN** the CLI runs with `malformed-template.yaml` or `malformed-overrides.yaml`
- **THEN** the process MUST exit non-zero, MUST NOT produce a meaningful output file, and stderr MUST identify template or overrides parsing failure

#### Scenario: Override configuration is semantically invalid
- **WHEN** the CLI runs with `invalid-match-type-overrides.yaml` or `missing-replace-value-overrides.yaml`
- **THEN** the process MUST exit non-zero, MUST NOT produce a meaningful output file, and stderr MUST identify the unsupported match type or missing replacement values

#### Scenario: Output destination cannot be written
- **WHEN** the CLI runs with a non-writable or invalid output path
- **THEN** the process MUST exit non-zero, MUST NOT produce a meaningful output file, and stderr MUST identify the write failure

### Requirement: The suite remains black-box and implementation-agnostic
The integration suite SHALL exercise only the real CLI entrypoint or built binary, create isolated temporary output locations per test, and use exact golden-file comparison for success cases plus substring stderr assertions for failures.

#### Scenario: Test harness executes the CLI as an external process
- **WHEN** an integration test is implemented from this specification
- **THEN** it MUST pass explicit file paths for all inputs, invoke the CLI without calling internal parsing or rendering helpers directly, read output from disk after process completion, and compare success results to the expected fixture file

#### Scenario: Repository has fixtures but no implementation stack yet
- **WHEN** the repository has the documented fixtures and OpenSpec artifacts but no real CLI entrypoint or test runner
- **THEN** the suite MAY remain intentionally non-runnable, and the project MUST NOT add fake failing tests or placeholder executables solely to simulate progress

### Requirement: Overlapping override precedence is deferred until documented
The integration specification SHALL identify the overlapping-rule case as deferred and SHALL NOT require an executable precedence test until precedence behavior is explicitly documented.

#### Scenario: Multiple override rules match the same cluster before precedence is documented
- **WHEN** the fixture set includes `overlapping-overrides.yaml`
- **THEN** the specification MUST record the case for future coverage but MUST NOT require an active integration test that locks precedence behavior yet
