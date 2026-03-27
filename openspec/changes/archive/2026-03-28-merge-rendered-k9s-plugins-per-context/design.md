## Context

The current generator resolves the active cluster from kubeconfig, renders one plugin document, and writes it to a single caller-provided output path. K9s also supports context-specific plugin files under `$XDG_DATA_HOME/k9s/clusters/<cluster>/<context>/plugins.yaml`, and one cluster can appear under multiple kubeconfig contexts. This change adds a K9s-native installation path while preserving the existing standalone render workflow.

## Goals / Non-Goals

**Goals:**
- Write generated plugins into the K9s context-specific data-home layout for every kubeconfig context that references the active cluster.
- Merge generated plugin definitions with any existing `plugins.yaml` content so unrelated plugins stay intact.
- Preserve the current explicit output-path workflow for users and tests that want a standalone rendered file.
- Keep behavior deterministic and easy to test with fixture-backed command and e2e coverage.

**Non-Goals:**
- Render different plugin content per context within the same cluster.
- Add support for multiple template files or multiple override files in one run.
- Manage plugins in `$XDG_CONFIG_HOME` or other non-context-specific K9s plugin directories.

## Decisions

- Add a K9s installation mode alongside the existing explicit output mode. The command should continue to support writing a standalone file, but it also needs a path that targets the K9s data-home layout so users do not need a separate copy step.
  - Alternative considered: replacing `--output` with an implicit K9s install destination. Rejected because it would break the current workflow and make tests less explicit.
- Resolve install targets by enumerating all kubeconfig contexts whose `context.cluster` matches the active cluster selected by `current-context`. The rendered plugin content is cluster-specific today, so reusing the same rendered output for each matching context is the simplest behavior and matches the stated assumption that those contexts share plugin content for now.
  - Alternative considered: writing only the active context path. Rejected because K9s loads plugins per context path, and users would still need manual duplication for other contexts on the same cluster.
- Merge at the plugin-key level. Generated plugin definitions should overwrite existing plugins with the same name, while unrelated existing plugins remain in the destination file.
  - Alternative considered: overwriting the entire destination file. Rejected because it would remove unrelated manually managed plugins.
- Normalize installed files to a complete `plugins:` mapping. The generator already validates rendered YAML; emitting a single normalized structure makes merge behavior deterministic and simplifies tests.
  - Alternative considered: preserving the original destination-file shape. Rejected because K9s accepts multiple shapes, but preserving shape adds complexity without user value.
- Resolve the K9s data-home root from `XDG_DATA_HOME`, with a testable override in command inputs if needed, and fall back to the platform-appropriate default when the environment variable is unset.
  - Alternative considered: hard-coding one absolute base path. Rejected because it would be brittle across developer machines and CI.

## Risks / Trade-offs

- [Existing destination file uses an unexpected YAML shape] -> Normalize only after successfully parsing a supported plugin document shape, and fail with an actionable error when the file cannot be merged safely.
- [Plugin name collisions overwrite user-managed entries] -> Limit overwrite behavior to generated plugin keys and document that same-name plugins are intentionally refreshed by the generator.
- [Multiple context writes partially succeed] -> Stop on the first failure and surface the target path in the error so the user can correct permissions or malformed files.
- [Added CLI mode increases command complexity] -> Keep modes explicit and cover them with command-level tests so behavior stays discoverable.

## Migration Plan

- No data migration is required for existing users who keep using explicit output paths.
- Users who adopt the new K9s install mode will generate or update files under their K9s data-home directories on the next run.
- Rollback is limited to removing the new mode and deleting generated context plugin files if needed.

## Open Questions

- No open questions remain for proposal-level planning; per-context customization stays explicitly out of scope for this change.
