## ADDED Requirements

### Requirement: Feature verification can be completed inside the devcontainer
The repository SHALL support feature-completion verification by running the documented build and test workflow inside the devcontainer.

#### Scenario: Generator feature is verified in the devcontainer
- **WHEN** the Go/Cobra generator change reaches its final verification stage
- **THEN** its completion criteria MUST include successful build and test execution in the devcontainer

#### Scenario: Repository documents devcontainer verification commands
- **WHEN** the devcontainer is added to the repository
- **THEN** the repository MUST document how to run the standard build and test commands inside that environment
