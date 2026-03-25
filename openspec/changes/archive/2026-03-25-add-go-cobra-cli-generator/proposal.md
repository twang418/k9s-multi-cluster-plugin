## Why

The repository already defines the desired fixture contract and rendering behavior, but it still lacks a real implementation for generating K9s plugin YAML from a template, kubeconfig, and override data. Adding a Go CLI with Cobra now creates the first concrete product slice and enables future feature-by-feature testing against the documented fixtures.

## What Changes

- Add a Go-based CLI using Cobra as the command framework.
- Implement a command that accepts a template folder, kubeconfig file, and override folder, with output location configurable for testing and local use.
- Load the active cluster from kubeconfig, resolve matching override values, and render gomplate-style template expressions into a generated K9s plugin file.
- Generate output that matches K9s plugin YAML expectations and preserves template-defined defaults when no override matches.
- Establish the initial project structure, dependency manifest, and documentation needed to build and run the CLI locally.

## Capabilities

### New Capabilities
- `k9s-plugin-generator`: Generate a K9s plugin YAML file from template and override inputs using the active kubeconfig cluster.
- `cli-command-surface`: Provide a Cobra-based CLI interface for supplying input folders/files and output location overrides.

### Modified Capabilities

None.

## Impact

- Adds the first implementation stack to the repository: Go plus Cobra.
- Introduces a real buildable CLI entrypoint, dependency manifest, and package structure.
- Creates the implementation foundation needed for later runnable end-to-end tests.
- Affects future docs and agent guidance because the repo will no longer be documentation-only.
- Depends on the new `add-go-devcontainer` change for reproducible final build/test verification.
