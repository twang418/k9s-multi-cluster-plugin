package generator

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
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

func TestLoadClusterSelectionIncludesMatchingContexts(t *testing.T) {
	t.Parallel()

	selection, err := loadClusterSelection(filepath.Join("..", "..", "testdata", "kubeconfig", "active-org1-multi-context.yaml"))
	if err != nil {
		t.Fatalf("loadClusterSelection returned error: %v", err)
	}
	if selection.ActiveCluster != "org1-dev" {
		t.Fatalf("expected active cluster org1-dev, got %q", selection.ActiveCluster)
	}
	want := []string{"org1-admin", "org1-context"}
	if !reflect.DeepEqual(selection.ContextNames, want) {
		t.Fatalf("expected context names %#v, got %#v", want, selection.ContextNames)
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

func TestGenerateInstallsMergedPluginsForAllMatchingContexts(t *testing.T) {
	t.Parallel()

	dataRoot := filepath.Join(t.TempDir(), "xdg-data")
	adminPath := filepath.Join(dataRoot, "k9s", "clusters", "org1-dev", "org1-admin", "plugins.yaml")
	contextPath := filepath.Join(dataRoot, "k9s", "clusters", "org1-dev", "org1-context", "plugins.yaml")

	writeTestFile(t, adminPath, "plugins:\n  existing-admin:\n    shortCut: Shift-A\n    description: Keep me\n")
	writeTestFile(t, contextPath, "plugins:\n  debug:\n    shortCut: Shift-X\n    description: Old debug\n  existing-context:\n    shortCut: Shift-C\n    description: Keep me too\n")

	result, err := Generate(Request{
		KubeconfigPath: filepath.Join("..", "..", "testdata", "kubeconfig", "active-org1-multi-context.yaml"),
		TemplateDir:    filepath.Join("..", "..", "testdata", "template-single"),
		OverrideDir:    filepath.Join("..", "..", "testdata", "overrides-single"),
		OutputMode:     OutputModeK9s,
		K9sDataDir:     dataRoot,
	})
	if err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}
	if result.InstallRoot != dataRoot {
		t.Fatalf("expected install root %q, got %q", dataRoot, result.InstallRoot)
	}
	if len(result.OutputPaths) != 2 {
		t.Fatalf("expected 2 output paths, got %#v", result.OutputPaths)
	}

	adminPlugins := readPluginsFile(t, adminPath)
	if _, ok := adminPlugins["existing-admin"]; !ok {
		t.Fatalf("expected existing-admin plugin to be preserved, got %#v", adminPlugins)
	}
	assertDebugImage(t, adminPlugins, "1111.dkr.ecr.ap-southeast-2.amazonaws.com/busybox:unstable-uclibc:1.37.0")

	contextPlugins := readPluginsFile(t, contextPath)
	if _, ok := contextPlugins["existing-context"]; !ok {
		t.Fatalf("expected existing-context plugin to be preserved, got %#v", contextPlugins)
	}
	assertDebugImage(t, contextPlugins, "1111.dkr.ecr.ap-southeast-2.amazonaws.com/busybox:unstable-uclibc:1.37.0")
	debug := contextPlugins["debug"].(map[string]any)
	if debug["description"] == "Old debug" {
		t.Fatalf("expected generated debug plugin to replace existing debug definition, got %#v", debug)
	}
}

func TestGenerateInstallFailsForMalformedExistingPlugins(t *testing.T) {
	t.Parallel()

	dataRoot := filepath.Join(t.TempDir(), "xdg-data")
	targetPath := filepath.Join(dataRoot, "k9s", "clusters", "org1-dev", "org1-admin", "plugins.yaml")
	writeTestFile(t, targetPath, "plugins: [broken\n")

	_, err := Generate(Request{
		KubeconfigPath: filepath.Join("..", "..", "testdata", "kubeconfig", "active-org1-multi-context.yaml"),
		TemplateDir:    filepath.Join("..", "..", "testdata", "template-single"),
		OverrideDir:    filepath.Join("..", "..", "testdata", "overrides-single"),
		OutputMode:     OutputModeK9s,
		K9sDataDir:     dataRoot,
	})
	if err == nil {
		t.Fatal("expected malformed existing plugins error")
	}
	if !strings.Contains(err.Error(), targetPath) {
		t.Fatalf("expected error to include target path, got %v", err)
	}
}

func writeTestFile(t *testing.T, path, content string) {
	t.Helper()

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("create parent directory for %q: %v", path, err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write test file %q: %v", path, err)
	}
}

func readPluginsFile(t *testing.T, path string) map[string]any {
	t.Helper()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read plugins file %q: %v", path, err)
	}

	var doc struct {
		Plugins map[string]any `yaml:"plugins"`
	}
	if err := yaml.Unmarshal(data, &doc); err != nil {
		t.Fatalf("unmarshal plugins file %q: %v", path, err)
	}
	return doc.Plugins
}

func assertDebugImage(t *testing.T, plugins map[string]any, image string) {
	t.Helper()

	debugValue, ok := plugins["debug"]
	if !ok {
		t.Fatalf("expected debug plugin in %#v", plugins)
	}
	debug, ok := debugValue.(map[string]any)
	if !ok {
		t.Fatalf("expected debug plugin map, got %#v", debugValue)
	}
	args, ok := debug["args"].([]any)
	if !ok {
		t.Fatalf("expected debug args slice, got %#v", debug["args"])
	}
	for _, arg := range args {
		argString, ok := arg.(string)
		if ok && strings.Contains(argString, image) {
			return
		}
	}
	t.Fatalf("expected debug args to contain image %q, got %#v", image, args)
}

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
