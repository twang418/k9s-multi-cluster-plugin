package cmd

import (
	"bytes"
	"path/filepath"
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
