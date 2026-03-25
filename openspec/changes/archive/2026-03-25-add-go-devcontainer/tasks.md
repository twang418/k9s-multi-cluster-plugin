## 1. Devcontainer Setup

- [x] 1.1 Add a `.devcontainer/` configuration with a Go-capable container image
- [x] 1.2 Configure the container with the minimum tooling needed to run repository Go commands
- [x] 1.3 Document how to open the repository and run commands inside the devcontainer

## 2. Verification Workflow

- [x] 2.1 Verify the devcontainer can run `go mod download`
- [x] 2.2 Verify the devcontainer can run `go build ./...`
- [x] 2.3 Verify the devcontainer can run `go test ./...`

## 3. Generator Change Follow-Up

- [x] 3.1 Update `add-go-cobra-cli-generator` to depend on devcontainer-backed verification for final completion
- [x] 3.2 Re-run generator build and test verification in the devcontainer before considering that feature complete
