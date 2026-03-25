# k9s-multi-cluster-plugin

This project is for a CLI that generates K9s plugin configuration for multiple clusters from a shared template plus cluster-specific overrides.

## Approach

Keep 4 things separate:

- a kubeconfig file 
- a K9s plugin template
- a CLI override file for cluster-specific values
- a generated K9s plugin file

The template stays close to normal K9s plugin YAML, but uses gomplate-style
template expressions for substitution and defaults. The override file is custom
input for this CLI and is not native K9s plugin syntax. The CLI reads cluster
names from kubeconfig, matches them against the override rules, and then renders
the template with the matched values.

## Example K9s Plugin Template

```yaml
plugins:
  debug:
    shortCut: Shift-D
    description: Add debug container
    dangerous: true
    scopes:
      - containers
    command: bash
    background: false
    confirm: true
    inputs:
      - name: profile
        label: Debug profile
        type: dropdown
        required: true
        default: sysadmin
        options:
          - general
          - baseline
          - restricted
          - netadmin
          - sysadmin
          - legacy
      - name: share_processes
        label: Share processes
        type: bool
        required: true
        default: true
    args:
      - -c
      - >-
        kubectl debug -it --context "$CONTEXT" -n "$NAMESPACE" "$POD"
        --target "$NAME"
        --image "{{ .image | default "busybox" }}"
        --profile "$INPUT_PROFILE"
        $([ "$INPUT_SHARE_PROCESSES" = "true" ] && echo "--share-processes")
        -- sh
```

In this template, `{{ .image | default "busybox" }}` means the CLI should use
the matched cluster override for `image` when one exists, and otherwise fall
back to the default value `busybox` defined directly in the template.

## Template Syntax Direction

The intended template syntax is gomplate-style, based on Go templates.

- replacement values are exposed as template data fields such as `.image`
- defaults can be defined in the template itself with functions such as
  `default`
- future template use may include more than scalar substitution, so the syntax
  should stay compatible with richer gomplate-style expressions

The goal is to keep the template close to normal K9s YAML while still allowing
controlled template power where needed.

## Initial Folder Convention

The first implementation treats the template folder and override folder as
directory inputs with a simple convention:

- the template folder must contain exactly one `.yaml` or `.yml` file
- the override folder must contain exactly one `.yaml` or `.yml` file
- the template file must define exactly one plugin under `plugins`

This keeps the first CLI slice deterministic and testable while leaving room for
future expansion to multiple files.

## Example CLI Override File

```yaml
pluginOverrides:
  debug:
    clusters:
      - match:
          type: regex
          value: ".*org1.*"
        replace:
          image: "1111.dkr.ecr.ap-southeast-2.amazonaws.com/busybox:unstable-uclibc:1.37.0"
      - match:
          type: regex
          value: ".*org2.*"
        replace:
          image: "2222.dkr.ecr.ap-southeast-2.amazonaws.com/busybox:unstable-uclibc:1.37.0"
      - match:
          type: list
          values:
            - "org3-bu1"
            - "org4-bu1"
        replace:
          image: "3333.dkr.ecr.ap-southeast-2.amazonaws.com/busybox:unstable-uclibc:1.37.0"
```

## Expected Result

For the active cluster name found in kubeconfig, the CLI resolves the matching
rule and renders the gomplate-style template. If no rule matches, the CLI uses
the template-defined default value and writes the final K9s plugin YAML.

## CLI Usage

The initial command shape is:

```bash
go run . generate \
  --kubeconfig ./testdata/kubeconfig/active-org1.yaml \
  --template-dir ./testdata/template-single \
  --override-dir ./testdata/overrides-single \
  --output ./out/plugin.yaml
```

This writes the rendered K9s plugin YAML to the requested output path.

## Devcontainer

The repository includes a devcontainer so Go build and test commands can run in
a reproducible environment even when the host machine does not already have Go
installed.

Open the repository in the devcontainer, then run:

```bash
go mod download
go build ./...
go test ./...
```

The devcontainer is the intended environment for final feature verification.
