package generator

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadActiveCluster(t *testing.T) {
	t.Parallel()

	cluster, err := loadActiveCluster(filepath.Join("..", "..", "testdata", "kubeconfig", "active-org3-bu1.yaml"))
	if err != nil {
		t.Fatalf("loadActiveCluster returned error: %v", err)
	}
	if cluster != "org3-bu1" {
		t.Fatalf("expected org3-bu1, got %q", cluster)
	}
}

func TestLoadActiveClusterMalformedKubeconfig(t *testing.T) {
	t.Parallel()

	_, err := loadActiveCluster(filepath.Join("..", "..", "testdata", "kubeconfig", "malformed.yaml"))
	if err == nil {
		t.Fatal("expected malformed kubeconfig error")
	}
}

func TestLoadReplacements(t *testing.T) {
	t.Parallel()

	values, err := loadReplacements(filepath.Join("..", "..", "testdata", "overrides", "standard-overrides.yaml"), "debug", "org1-dev")
	if err != nil {
		t.Fatalf("loadReplacements returned error: %v", err)
	}
	if values["image"] != "1111.dkr.ecr.ap-southeast-2.amazonaws.com/busybox:unstable-uclibc:1.37.0" {
		t.Fatalf("unexpected image replacement: %#v", values["image"])
	}
	if got, ok := values["image"].(string); !ok || got == "" {
		t.Fatalf("expected string image replacement, got %#v", values["image"])
	}

	values, err = loadReplacements(filepath.Join("..", "..", "testdata", "overrides", "standard-overrides.yaml"), "debug", "sandbox-team-a")
	if err != nil {
		t.Fatalf("loadReplacements unmatched returned error: %v", err)
	}
	if len(values) != 0 {
		t.Fatalf("expected no values for unmatched cluster, got %#v", values)
	}
}

func TestLoadTemplateAllowsTemplateDirectivesOutsideQuotedStrings(t *testing.T) {
	t.Parallel()

	templatePath := filepath.Join(t.TempDir(), "template.yaml")
	content := []byte("plugins:\n  debug:\n    env: {{ .env | default (list) }}\n")
	if err := os.WriteFile(templatePath, content, 0o644); err != nil {
		t.Fatalf("write temp template: %v", err)
	}

	_, pluginName, err := loadTemplate(templatePath)
	if err != nil {
		t.Fatalf("loadTemplate returned error: %v", err)
	}
	if pluginName != "debug" {
		t.Fatalf("expected plugin name debug, got %q", pluginName)
	}
}

func TestLoadReplacementsInvalidConfiguration(t *testing.T) {
	t.Parallel()

	_, err := loadReplacements(filepath.Join("..", "..", "testdata", "overrides", "invalid-match-type-overrides.yaml"), "debug", "org1-dev")
	if err == nil {
		t.Fatal("expected invalid match type error")
	}

	_, err = loadReplacements(filepath.Join("..", "..", "testdata", "overrides", "missing-replace-value-overrides.yaml"), "debug", "org1-dev")
	if err == nil {
		t.Fatal("expected missing replacement values error")
	}
}

func TestRenderTemplateUsesDefault(t *testing.T) {
	t.Parallel()

	data, err := os.ReadFile(filepath.Join("..", "..", "testdata", "template", "debug-template.yaml"))
	if err != nil {
		t.Fatalf("read template: %v", err)
	}

	rendered, err := renderTemplate(data, map[string]any{})
	if err != nil {
		t.Fatalf("renderTemplate returned error: %v", err)
	}
	if string(rendered) == string(data) {
		t.Fatalf("expected rendered output to differ from template source")
	}
	if got := string(rendered); !contains(got, `--image "busybox"`) {
		t.Fatalf("expected default image in rendered output, got %s", got)
	}
}

func TestGenerateWritesOutput(t *testing.T) {
	t.Parallel()

	outputPath := filepath.Join(t.TempDir(), "plugin.yaml")
	result, err := Generate(Request{
		KubeconfigPath: filepath.Join("..", "..", "testdata", "kubeconfig", "active-unmatched.yaml"),
		TemplateDir:    filepath.Join("..", "..", "testdata", "template-single"),
		OverrideDir:    filepath.Join("..", "..", "testdata", "overrides-single"),
		OutputPath:     outputPath,
	})
	if err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}
	if result.OutputPath != outputPath {
		t.Fatalf("expected output path %q, got %q", outputPath, result.OutputPath)
	}

	data, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("read generated output: %v", err)
	}
	if got := string(data); !contains(got, `--image "busybox"`) {
		t.Fatalf("expected busybox fallback in output, got %s", got)
	}
}

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
