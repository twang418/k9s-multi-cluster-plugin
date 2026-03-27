## Why

The CLI currently writes one rendered plugin file to a caller-provided path, but K9s also supports context-specific plugin files under `$XDG_DATA_HOME/k9s/clusters/<cluster>/<context>/plugins.yaml`. Supporting that layout lets the generator place plugins where K9s already discovers them and avoids manual copy or merge steps when a cluster has multiple contexts.

## What Changes

- Add a generation mode that targets the K9s context-specific plugin directory layout under `$XDG_DATA_HOME/k9s/clusters/<cluster>/<context>/plugins.yaml`.
- Render plugin content once for the active cluster and write the same merged plugin set to every kubeconfig context that points at that cluster.
- Merge newly rendered plugin definitions with any existing `plugins.yaml` content at each context path instead of overwriting unrelated plugins.
- Keep the current explicit output-path workflow available for users who still want to write a standalone rendered file.
- Document the K9s plugin discovery layout and the CLI behavior for multi-context clusters.

## Capabilities

### New Capabilities
- `k9s-context-plugin-output`: Generate and merge rendered plugin definitions into K9s context-specific plugin files for every context that references the active cluster.

### Modified Capabilities

## Impact

- Affects CLI flags and generate-command behavior in `cmd/generate.go`.
- Affects kubeconfig parsing, output-path resolution, and YAML merge logic in `internal/generator/generator.go`.
- Adds or updates command, generator, and e2e coverage for context-specific output and merge behavior.
- Updates `README.md` to describe K9s output locations and the new workflow.
