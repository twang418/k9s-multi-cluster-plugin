## ADDED Requirements

### Requirement: CLI can install rendered plugins into K9s context plugin paths
The system SHALL support writing rendered plugin definitions into `$XDG_DATA_HOME/k9s/clusters/<cluster>/<context>/plugins.yaml` for every kubeconfig context that references the active cluster selected by `current-context`.

#### Scenario: Active cluster is referenced by multiple contexts
- **WHEN** the generate command runs in K9s install mode and the active cluster is referenced by multiple kubeconfig contexts
- **THEN** the system MUST write the rendered plugin definitions to `plugins.yaml` under each matching `<cluster>/<context>` directory

#### Scenario: XDG data home is resolved for install output
- **WHEN** the generate command runs in K9s install mode
- **THEN** the system MUST resolve the base install directory from `XDG_DATA_HOME` or its documented default before creating the cluster and context plugin paths

### Requirement: Installed context plugin files are merged by plugin key
The system SHALL merge generated plugin definitions into each target `plugins.yaml` file without removing unrelated existing plugins.

#### Scenario: Destination file already contains unrelated plugins
- **WHEN** a target context `plugins.yaml` already defines plugins with names not present in the generated output
- **THEN** the system MUST preserve those unrelated plugins and add the generated plugins into the same resulting `plugins` mapping

#### Scenario: Destination file already contains the generated plugin name
- **WHEN** a target context `plugins.yaml` already defines a plugin with the same name as a generated plugin
- **THEN** the system MUST replace that existing plugin definition with the newly generated plugin definition

### Requirement: Standalone output path workflow remains available
The system SHALL continue to support writing the rendered plugin YAML to a caller-selected standalone output path.

#### Scenario: Explicit output path is requested
- **WHEN** the generate command is run with an explicit output-path workflow instead of K9s install mode
- **THEN** the system MUST write the rendered plugin YAML to that requested path and MUST NOT require the K9s context-directory layout
