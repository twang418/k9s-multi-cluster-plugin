## ADDED Requirements

### Requirement: Repository provides a reproducible Go devcontainer
The repository SHALL define a devcontainer configuration that provides a working Go toolchain suitable for building and testing the CLI.

#### Scenario: Contributor opens the repository in the devcontainer
- **WHEN** the repository is opened in the configured devcontainer
- **THEN** the environment MUST provide the Go toolchain needed to run the documented build and test commands for this repository

#### Scenario: Go commands run without relying on host machine setup
- **WHEN** a contributor or agent uses the devcontainer for repository work
- **THEN** `go build ./...` and `go test ./...` MUST be runnable inside the container without depending on a separately configured host Go installation
