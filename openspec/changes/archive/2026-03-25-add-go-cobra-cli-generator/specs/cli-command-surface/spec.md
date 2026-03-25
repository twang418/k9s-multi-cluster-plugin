## ADDED Requirements

### Requirement: Cobra CLI exposes a generation command with explicit path inputs
The system SHALL provide a Cobra-based CLI command that accepts the kubeconfig path, template folder path, override folder path, and output file path as explicit inputs.

#### Scenario: User provides all required paths
- **WHEN** the command is invoked with the required kubeconfig, template folder, override folder, and output path inputs
- **THEN** the command MUST run the generation workflow using those supplied locations

#### Scenario: Required path input is missing
- **WHEN** a required input path is omitted
- **THEN** the command MUST fail with usage guidance that identifies the missing required argument or flag

### Requirement: CLI input paths are test-friendly and not hard-coded to one environment
The system SHALL allow callers to override filesystem locations so the generator can be exercised in temporary directories and local test environments.

#### Scenario: Template and override folders are outside default locations
- **WHEN** the command is run with template and override folders in arbitrary filesystem locations
- **THEN** the generator MUST read from those provided locations without requiring a repository-fixed directory layout

#### Scenario: Output path is in a temporary directory
- **WHEN** the command is run with an output file path inside a temporary or test-specific directory
- **THEN** the generator MUST write the generated plugin YAML there if the location is writable
