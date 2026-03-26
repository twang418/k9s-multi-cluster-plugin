## ADDED Requirements

### Requirement: Executable integration tests are enabled only against a real CLI
The project SHALL enable fixture-backed end-to-end integration tests only when a
real CLI entrypoint or built binary exists and a real test runner can execute
the suite.

#### Scenario: Success cases are enabled feature by feature
- **WHEN** a real CLI implementation supports a documented render behavior
- **THEN** the corresponding fixture-backed success case MAY be enabled without
  waiting for unrelated behaviors to exist

#### Scenario: Success output verification uses a file diff
- **WHEN** a fixture-backed success case verifies rendered output
- **THEN** the harness MUST render the template to a real output file and diff
  that file against the expected golden file to confirm an exact match

#### Scenario: Rendered test workspaces can be preserved for inspection
- **WHEN** the end-to-end suite runs with its auto-clean flag enabled
- **THEN** each test workspace MUST be removed after the test completes by
  default
- **AND WHEN** the suite runs with the auto-clean flag disabled
- **THEN** the rendered output workspace MUST remain on disk so the generated
  template can be inspected after the test run

#### Scenario: Manual rendered output can be written to a gitignored directory
- **WHEN** a developer runs the real CLI manually to inspect generated output
- **THEN** the repository MAY use a local `output/` directory for rendered YAML
- **AND** that directory MUST be gitignored so inspection artifacts do not
  become tracked changes

#### Scenario: Failure cases are enabled against the real CLI
- **WHEN** the CLI implements the relevant validation and error paths
- **THEN** the corresponding fixture-backed failure cases MAY be enabled using
  the documented non-zero exit, output, and stderr assertions

#### Scenario: No fake suite is added before the real stack exists
- **WHEN** the repository still lacks a real CLI entrypoint or test runner
- **THEN** this change MUST remain non-runnable and MUST NOT introduce fake
  failing tests or placeholder executables
