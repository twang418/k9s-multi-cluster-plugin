# k9s-multi-cluster-plugin

This project is for a CLI that generates K9s plugin configuration for multiple clusters from a shared template plus cluster-specific overrides.

## Approach

Keep 4 things separate:

- a kubeconfig file 
- a K9s plugin template
- a CLI override file for cluster-specific values
- a generated K9s plugin file

The template stays close to normal K9s plugin YAML. The override file is custom input for this CLI and is not native K9s plugin syntax.
The CLI reads cluster names from kubeconfig, matches them against the override rules, and then applies replacements to the template.

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
        --image "{{ image }}"
        --profile "$INPUT_PROFILE"
        $([ "$INPUT_SHARE_PROCESSES" = "true" ] && echo "--share-processes")
        -- sh
```

In this template, `{{ image }}` is a placeholder that the CLI replaces based on the matched cluster name from kubeconfig.

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

For the active cluster name found in kubeconfig, the CLI resolves the matching rule, replaces `{{ image }}` in the template, and writes the final K9s plugin YAML.
