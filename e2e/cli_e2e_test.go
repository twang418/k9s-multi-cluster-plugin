package e2e_test

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"
)

var (
	buildBinaryOnce sync.Once
	builtBinaryPath string
	buildBinaryErr  error
	autoCleanE2E    = flag.Bool("e2e-auto-clean", true, "remove rendered e2e test workspaces after each test")
)

func TestCLIEndToEndSuccessCases(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		kubeconfigFile  string
		templateFile    string
		overridesFile   string
		expectedOutFile string
	}{
		{
			name:            "renders debug plugin for active org1 cluster",
			kubeconfigFile:  fixturePath("kubeconfig", "active-org1.yaml"),
			templateFile:    fixturePath("template", "debug-template.yaml"),
			overridesFile:   fixturePath("overrides", "standard-overrides.yaml"),
			expectedOutFile: fixturePath("expected", "rendered-org1.yaml"),
		},
		{
			name:            "renders debug plugin for active org3-bu1 cluster",
			kubeconfigFile:  fixturePath("kubeconfig", "active-org3-bu1.yaml"),
			templateFile:    fixturePath("template", "debug-template.yaml"),
			overridesFile:   fixturePath("overrides", "standard-overrides.yaml"),
			expectedOutFile: fixturePath("expected", "rendered-org3-bu1.yaml"),
		},
		{
			name:            "renders template default image when active cluster has no matching override",
			kubeconfigFile:  fixturePath("kubeconfig", "active-unmatched.yaml"),
			templateFile:    fixturePath("template", "debug-template.yaml"),
			overridesFile:   fixturePath("overrides", "standard-overrides.yaml"),
			expectedOutFile: fixturePath("expected", "unchanged-unmatched.yaml"),
		},
		{
			name:            "writes unchanged template when no matching expression exists in the template",
			kubeconfigFile:  fixturePath("kubeconfig", "active-org1.yaml"),
			templateFile:    fixturePath("template", "missing-placeholder-template.yaml"),
			overridesFile:   fixturePath("overrides", "standard-overrides.yaml"),
			expectedOutFile: fixturePath("expected", "unchanged-no-placeholder.yaml"),
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			workspace := newTestWorkspace(t)
			outputPath := filepath.Join(workspace, "generated", "plugin.yaml")
			templateDir := copyFixtureToSingleFileDir(t, workspace, "template", test.templateFile)
			overrideDir := copyFixtureToSingleFileDir(t, workspace, "overrides", test.overridesFile)

			result := runCLI(t,
				"generate",
				"--kubeconfig", test.kubeconfigFile,
				"--template-dir", templateDir,
				"--override-dir", overrideDir,
				"--output", outputPath,
			)

			if result.exitCode != 0 {
				t.Fatalf("expected exit code 0, got %d; stderr=%s", result.exitCode, result.stderr)
			}

			assertFilesMatchDiff(t, test.expectedOutFile, outputPath)
		})
	}
}

func TestCLIEndToEndFailureCases(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                string
		kubeconfigFile      string
		templateFile        string
		overridesFile       string
		outputPathFactory   func(*testing.T, string) string
		expectedStderrParts []string
	}{
		{
			name:          "fails for missing kubeconfig path",
			templateFile:  fixturePath("template", "debug-template.yaml"),
			overridesFile: fixturePath("overrides", "standard-overrides.yaml"),
			outputPathFactory: func(t *testing.T, workspace string) string {
				return filepath.Join(workspace, "plugin.yaml")
			},
			expectedStderrParts: []string{"read kubeconfig", "no such file"},
		},
		{
			name:           "fails for malformed kubeconfig",
			kubeconfigFile: fixturePath("kubeconfig", "malformed.yaml"),
			templateFile:   fixturePath("template", "debug-template.yaml"),
			overridesFile:  fixturePath("overrides", "standard-overrides.yaml"),
			outputPathFactory: func(t *testing.T, workspace string) string {
				return filepath.Join(workspace, "plugin.yaml")
			},
			expectedStderrParts: []string{"parse kubeconfig"},
		},
		{
			name:           "fails for malformed template yaml",
			kubeconfigFile: fixturePath("kubeconfig", "active-org1.yaml"),
			templateFile:   fixturePath("template", "malformed-template.yaml"),
			overridesFile:  fixturePath("overrides", "standard-overrides.yaml"),
			outputPathFactory: func(t *testing.T, workspace string) string {
				return filepath.Join(workspace, "plugin.yaml")
			},
			expectedStderrParts: []string{"template", "plugin"},
		},
		{
			name:           "fails for malformed overrides yaml",
			kubeconfigFile: fixturePath("kubeconfig", "active-org1.yaml"),
			templateFile:   fixturePath("template", "debug-template.yaml"),
			overridesFile:  fixturePath("overrides", "malformed-overrides.yaml"),
			outputPathFactory: func(t *testing.T, workspace string) string {
				return filepath.Join(workspace, "plugin.yaml")
			},
			expectedStderrParts: []string{"parse overrides"},
		},
		{
			name:           "fails for missing current-context resolution",
			kubeconfigFile: fixturePath("kubeconfig", "missing-current-context.yaml"),
			templateFile:   fixturePath("template", "debug-template.yaml"),
			overridesFile:  fixturePath("overrides", "standard-overrides.yaml"),
			outputPathFactory: func(t *testing.T, workspace string) string {
				return filepath.Join(workspace, "plugin.yaml")
			},
			expectedStderrParts: []string{"current-context", "does not match any context"},
		},
		{
			name:           "fails for missing cluster reference resolution",
			kubeconfigFile: fixturePath("kubeconfig", "missing-cluster-for-context.yaml"),
			templateFile:   fixturePath("template", "debug-template.yaml"),
			overridesFile:  fixturePath("overrides", "standard-overrides.yaml"),
			outputPathFactory: func(t *testing.T, workspace string) string {
				return filepath.Join(workspace, "plugin.yaml")
			},
			expectedStderrParts: []string{"references missing cluster"},
		},
		{
			name:           "fails for unsupported match type",
			kubeconfigFile: fixturePath("kubeconfig", "active-org1.yaml"),
			templateFile:   fixturePath("template", "debug-template.yaml"),
			overridesFile:  fixturePath("overrides", "invalid-match-type-overrides.yaml"),
			outputPathFactory: func(t *testing.T, workspace string) string {
				return filepath.Join(workspace, "plugin.yaml")
			},
			expectedStderrParts: []string{"unsupported match type"},
		},
		{
			name:           "fails for missing replacement values",
			kubeconfigFile: fixturePath("kubeconfig", "active-org1.yaml"),
			templateFile:   fixturePath("template", "debug-template.yaml"),
			overridesFile:  fixturePath("overrides", "missing-replace-value-overrides.yaml"),
			outputPathFactory: func(t *testing.T, workspace string) string {
				return filepath.Join(workspace, "plugin.yaml")
			},
			expectedStderrParts: []string{"missing replacement values"},
		},
		{
			name:           "fails for unwritable output path",
			kubeconfigFile: fixturePath("kubeconfig", "active-org1.yaml"),
			templateFile:   fixturePath("template", "debug-template.yaml"),
			overridesFile:  fixturePath("overrides", "standard-overrides.yaml"),
			outputPathFactory: func(t *testing.T, workspace string) string {
				path := filepath.Join(workspace, "existing-directory")
				if err := os.Mkdir(path, 0o755); err != nil {
					t.Fatalf("create output directory: %v", err)
				}
				return path
			},
			expectedStderrParts: []string{"write output", "directory"},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			workspace := newTestWorkspace(t)
			kubeconfigPath := test.kubeconfigFile
			if kubeconfigPath == "" {
				kubeconfigPath = filepath.Join(workspace, "missing-kubeconfig.yaml")
			}
			outputPath := test.outputPathFactory(t, workspace)
			templateDir := copyFixtureToSingleFileDir(t, workspace, "template", test.templateFile)
			overrideDir := copyFixtureToSingleFileDir(t, workspace, "overrides", test.overridesFile)

			result := runCLI(t,
				"generate",
				"--kubeconfig", kubeconfigPath,
				"--template-dir", templateDir,
				"--override-dir", overrideDir,
				"--output", outputPath,
			)

			if result.exitCode == 0 {
				t.Fatalf("expected non-zero exit code; stdout=%s", result.stdout)
			}
			assertNoMeaningfulOutput(t, outputPath)
			for _, part := range test.expectedStderrParts {
				if !strings.Contains(result.stderr, part) {
					t.Fatalf("expected stderr to contain %q, got %s", part, result.stderr)
				}
			}
		})
	}
}

type cliResult struct {
	exitCode int
	stdout   string
	stderr   string
}

func runCLI(t *testing.T, args ...string) cliResult {
	t.Helper()

	cmd := exec.Command(cliBinaryPath(t), args...)
	cmd.Dir = repoRoot(t)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err == nil {
		return cliResult{exitCode: 0, stdout: stdout.String(), stderr: stderr.String()}
	}

	exitErr, ok := err.(*exec.ExitError)
	if !ok {
		t.Fatalf("run cli: %v", err)
	}

	return cliResult{exitCode: exitErr.ExitCode(), stdout: stdout.String(), stderr: stderr.String()}
}

func cliBinaryPath(t *testing.T) string {
	t.Helper()

	buildBinaryOnce.Do(func() {
		buildDir, err := os.MkdirTemp("", "k9s-multi-cluster-plugin-bin-")
		if err != nil {
			buildBinaryErr = fmt.Errorf("create build temp dir: %w", err)
			return
		}

		builtBinaryPath = filepath.Join(buildDir, "k9s-multi-cluster-plugin")
		cmd := exec.Command("go", "build", "-o", builtBinaryPath, ".")
		cmd.Dir = repoRoot(t)
		output, err := cmd.CombinedOutput()
		if err != nil {
			buildBinaryErr = fmt.Errorf("build cli binary: %w: %s", err, strings.TrimSpace(string(output)))
		}
	})

	if buildBinaryErr != nil {
		t.Fatal(buildBinaryErr)
	}

	return builtBinaryPath
}

func copyFixtureToSingleFileDir(t *testing.T, workspace, dirName, sourcePath string) string {
	t.Helper()

	targetDir := filepath.Join(workspace, dirName)
	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		t.Fatalf("create %s dir: %v", dirName, err)
	}

	targetPath := filepath.Join(targetDir, filepath.Base(sourcePath))
	data := mustReadFile(t, sourcePath)
	if err := os.WriteFile(targetPath, data, 0o644); err != nil {
		t.Fatalf("write %s fixture copy: %v", dirName, err)
	}

	return targetDir
}

func newTestWorkspace(t *testing.T) string {
	t.Helper()

	workspace, err := os.MkdirTemp("", "k9s-multi-cluster-plugin-e2e-")
	if err != nil {
		t.Fatalf("create test workspace: %v", err)
	}

	if *autoCleanE2E {
		t.Cleanup(func() {
			if err := os.RemoveAll(workspace); err != nil {
				t.Fatalf("remove test workspace %q: %v", workspace, err)
			}
		})
	} else {
		t.Logf("preserving e2e workspace: %s", workspace)
	}

	return workspace
}

func mustReadFile(t *testing.T, path string) []byte {
	t.Helper()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read file %q: %v", path, err)
	}

	return data
}

func assertNoMeaningfulOutput(t *testing.T, path string) {
	t.Helper()

	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		t.Fatalf("stat output path %q: %v", path, err)
	}

	if info.IsDir() {
		return
	}

	if info.Size() != 0 {
		t.Fatalf("expected output file %q to be absent or empty, size=%d", path, info.Size())
	}
}

func assertFilesMatchDiff(t *testing.T, expectedPath, actualPath string) {
	t.Helper()

	cmd := exec.Command("diff", "-u", expectedPath, actualPath)
	output, err := cmd.CombinedOutput()
	if err == nil {
		return
	}

	exitErr, ok := err.(*exec.ExitError)
	if !ok {
		t.Fatalf("run diff for %q and %q: %v", expectedPath, actualPath, err)
	}

	if exitErr.ExitCode() == 1 {
		t.Fatalf("generated output did not match golden file:\n%s", output)
	}

	t.Fatalf("diff failed for %q and %q: %v\n%s", expectedPath, actualPath, err, output)
}

func fixturePath(parts ...string) string {
	allParts := append([]string{repoRoot(nil), "testdata"}, parts...)
	return filepath.Join(allParts...)
}

func repoRoot(t *testing.T) string {
	if t != nil {
		t.Helper()
	}

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		if t != nil {
			t.Fatal("resolve repository root")
		}
		panic("resolve repository root")
	}

	return filepath.Clean(filepath.Join(filepath.Dir(filename), ".."))
}
