## ADDED Requirements

### Requirement: Executable integration tests are enabled only against a real CLI
The project SHALL enable fixture-backed end-to-end integration tests only when a
real CLI entrypoint or built binary exists and a real test runner can execute
the suite.

#### Scenario: Success cases are enabled feature by feature
- **WHEN** a real CLI implementation supports a documented render behavior
- **THEN** the corresponding fixture-backed success case MAY be enabled without
  waiting for unrelated behaviors to exist

#### Scenario: Failure cases are enabled against the real CLI
- **WHEN** the CLI implements the relevant validation and error paths
- **THEN** the corresponding fixture-backed failure cases MAY be enabled using
  the documented non-zero exit, output, and stderr assertions

#### Scenario: No fake suite is added before the real stack exists
- **WHEN** the repository still lacks a real CLI entrypoint or test runner
- **THEN** this change MUST remain non-runnable and MUST NOT introduce fake
  failing tests or placeholder executables
