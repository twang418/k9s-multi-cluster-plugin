package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerateCommandRequiresFlags(t *testing.T) {
	t.Parallel()

	cmd := NewRootCommand()
	cmd.SetArgs([]string{"generate"})
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})

	if err := cmd.Execute(); err == nil {
		t.Fatal("expected missing flag error")
	}
}

func TestGenerateCommandWritesOutput(t *testing.T) {
	t.Parallel()

	outputPath := filepath.Join(t.TempDir(), "plugin.yaml")
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	cmd := NewRootCommand()
	cmd.SetOut(stdout)
	cmd.SetErr(stderr)
	cmd.SetArgs([]string{
		"generate",
		"--kubeconfig", filepath.Join("..", "testdata", "kubeconfig", "active-org1.yaml"),
		"--template-dir", filepath.Join("..", "testdata", "template-single"),
		"--override-dir", filepath.Join("..", "testdata", "overrides-single"),
		"--output", outputPath,
	})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute generate command: %v; stderr=%s", err, stderr.String())
	}
	if stdout.Len() == 0 {
		t.Fatal("expected command to report generated output")
	}
}

func TestGenerateCommandInstallsIntoK9sPaths(t *testing.T) {
	t.Parallel()

	dataRoot := filepath.Join(t.TempDir(), "xdg-data")
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	cmd := NewRootCommand()
	cmd.SetOut(stdout)
	cmd.SetErr(stderr)
	cmd.SetArgs([]string{
		"generate",
		"--kubeconfig", filepath.Join("..", "testdata", "kubeconfig", "active-org1-multi-context.yaml"),
		"--template-dir", filepath.Join("..", "testdata", "template-single"),
		"--override-dir", filepath.Join("..", "testdata", "overrides-single"),
		"--install-to-k9s",
		"--k9s-data-dir", dataRoot,
	})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute install generate command: %v; stderr=%s", err, stderr.String())
	}
	if !strings.Contains(stdout.String(), "generated 2 K9s plugin file(s)") {
		t.Fatalf("expected install output summary, got %q", stdout.String())
	}

	for _, contextName := range []string{"org1-admin", "org1-context"} {
		targetPath := filepath.Join(dataRoot, "k9s", "clusters", "org1-dev", contextName, "plugins.yaml")
		if _, err := os.Stat(targetPath); err != nil {
			t.Fatalf("expected generated plugins file %q: %v", targetPath, err)
		}
	}
}
