## ADDED Requirements

### Requirement: CLI generates K9s plugin YAML from template and cluster overrides
The system SHALL generate a K9s plugin YAML file by loading the active cluster from kubeconfig, selecting matching override values, and rendering a gomplate-style template.

#### Scenario: Active cluster matches regex override rule
- **WHEN** the generator runs with a kubeconfig whose active cluster matches a regex-based override rule and a template containing `{{ .image | default "busybox" }}`
- **THEN** the generated K9s plugin YAML MUST render the template using the matched `image` value from the override data

#### Scenario: Active cluster matches explicit list override rule
- **WHEN** the generator runs with a kubeconfig whose active cluster matches a list-based override rule
- **THEN** the generated K9s plugin YAML MUST render the template using the matched replacement values for that cluster

#### Scenario: No override rule matches active cluster
- **WHEN** the generator runs with a kubeconfig whose active cluster matches no override rule
- **THEN** the generated K9s plugin YAML MUST render using template-defined defaults and MUST NOT inject a hidden code-defined fallback for missing values

### Requirement: Generator resolves the active cluster from kubeconfig
The system SHALL determine the active cluster using `current-context` and the referenced kubeconfig context and cluster entries.

#### Scenario: Current context resolves to a cluster
- **WHEN** the kubeconfig contains multiple contexts and clusters and `current-context` references a valid context
- **THEN** the generator MUST use the cluster referenced by that context as the active cluster for override matching

#### Scenario: Current context cannot be resolved
- **WHEN** `current-context` is missing or points to a context or cluster that does not exist in kubeconfig
- **THEN** the generator MUST fail with an actionable error that identifies the unresolved kubeconfig reference

### Requirement: Generated output is written to a configurable file location
The system SHALL write the generated K9s plugin YAML to an output file location that can be overridden by command input.

#### Scenario: Output path is provided
- **WHEN** the generator command is run with an explicit output path
- **THEN** the generated plugin YAML MUST be written to that path instead of a hard-coded location

#### Scenario: Output path cannot be written
- **WHEN** the requested output path is invalid or not writable
- **THEN** the generator MUST fail with an actionable write error and MUST NOT report success
