## 1. Project Setup

- [x] 1.1 Initialize the repository as a Go module and add Cobra as a dependency
- [x] 1.2 Create the initial CLI entrypoint and root command structure for the generator
- [x] 1.3 Update `AGENTS.md` with the real Go build, test, and single-test commands once the stack exists

## 2. Command Surface

- [x] 2.1 Add a generation command that accepts kubeconfig path, template folder path, override folder path, and output file path inputs
- [x] 2.2 Validate required command inputs and return actionable usage errors for missing or invalid paths
- [x] 2.3 Define and document the first-pass file discovery convention for template and override folders

## 3. Input Loading And Resolution

- [x] 3.1 Implement kubeconfig loading and active-cluster resolution via `current-context`
- [x] 3.2 Implement override folder loading and cluster-rule matching for regex and list match types
- [x] 3.3 Implement actionable failures for unresolved kubeconfig references and semantically invalid override configuration

## 4. Template Rendering And Output

- [x] 4.1 Implement gomplate-style / Go-template-style rendering for the initial supported expression subset needed by the documented fixtures
- [x] 4.2 Render template-defined defaults when no override value is supplied for a referenced field
- [x] 4.3 Write the generated K9s plugin YAML to the configured output file path and surface write failures clearly

## 5. Verification

- [x] 5.1 Add focused automated tests for kubeconfig resolution, override matching, and template rendering
- [x] 5.2 Verify the command can generate plugin output from fixture-backed inputs in a temporary output location inside the devcontainer
- [x] 5.3 Update `README.md` with the actual CLI usage once the command behavior is implemented

## 6. Devcontainer-Backed Completion

- [x] 6.1 Complete the `add-go-devcontainer` prerequisite change
- [x] 6.2 Run `go build ./...` in the devcontainer successfully
- [x] 6.3 Run `go test ./...` in the devcontainer successfully
