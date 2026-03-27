package generator

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"gopkg.in/yaml.v3"
)

type OutputMode string

const (
	OutputModeFile OutputMode = "file"
	OutputModeK9s  OutputMode = "k9s"
)

type Request struct {
	KubeconfigPath string
	TemplateDir    string
	OverrideDir    string
	OutputPath     string
	OutputMode     OutputMode
	K9sDataDir     string
}

type Result struct {
	ActiveCluster string
	OutputPath    string
	OutputPaths   []string
	InstallRoot   string
}

type kubeconfig struct {
	CurrentContext string         `yaml:"current-context"`
	Contexts       []namedContext `yaml:"contexts"`
	Clusters       []namedCluster `yaml:"clusters"`
	Users          []namedUser    `yaml:"users"`
}

type namedContext struct {
	Name    string        `yaml:"name"`
	Context contextDetail `yaml:"context"`
}

type contextDetail struct {
	Cluster string `yaml:"cluster"`
}

type namedCluster struct {
	Name string `yaml:"name"`
}

type namedUser struct {
	Name string `yaml:"name"`
}

type overrideDocument struct {
	PluginOverrides map[string]pluginOverride `yaml:"pluginOverrides"`
}

type pluginOverride struct {
	Clusters []clusterRule `yaml:"clusters"`
}

type clusterRule struct {
	Match   matchRule      `yaml:"match"`
	Replace map[string]any `yaml:"replace"`
}

type matchRule struct {
	Type   string   `yaml:"type"`
	Value  string   `yaml:"value"`
	Values []string `yaml:"values"`
}

type clusterSelection struct {
	ActiveCluster string
	ContextNames  []string
}

func Generate(req Request) (Result, error) {
	req = normalizeRequest(req)

	if err := validateRequest(req); err != nil {
		return Result{}, err
	}

	templatePath, err := discoverSingleYAML(req.TemplateDir, "template folder")
	if err != nil {
		return Result{}, err
	}

	overridePath, err := discoverSingleYAML(req.OverrideDir, "override folder")
	if err != nil {
		return Result{}, err
	}

	selection, err := loadClusterSelection(req.KubeconfigPath)
	if err != nil {
		return Result{}, err
	}

	templateBytes, pluginName, err := loadTemplate(templatePath)
	if err != nil {
		return Result{}, err
	}

	replacements, err := loadReplacements(overridePath, pluginName, selection.ActiveCluster)
	if err != nil {
		return Result{}, err
	}

	rendered, err := renderTemplate(templateBytes, replacements)
	if err != nil {
		return Result{}, err
	}

	if req.OutputMode == OutputModeFile {
		if err := writeOutput(req.OutputPath, rendered); err != nil {
			return Result{}, err
		}

		return Result{
			ActiveCluster: selection.ActiveCluster,
			OutputPath:    req.OutputPath,
			OutputPaths:   []string{req.OutputPath},
		}, nil
	}

	installRoot, err := resolveK9sDataRoot(req.K9sDataDir)
	if err != nil {
		return Result{}, err
	}

	targetPaths := resolveK9sPluginPaths(installRoot, selection.ActiveCluster, selection.ContextNames)
	if err := installMergedPlugins(targetPaths, rendered); err != nil {
		return Result{}, err
	}

	return Result{
		ActiveCluster: selection.ActiveCluster,
		OutputPaths:   targetPaths,
		InstallRoot:   installRoot,
	}, nil
}

func normalizeRequest(req Request) Request {
	if req.OutputMode == "" && req.OutputPath != "" {
		req.OutputMode = OutputModeFile
	}
	return req
}

func validateRequest(req Request) error {
	if req.KubeconfigPath == "" {
		return errors.New("kubeconfig path is required")
	}
	if req.TemplateDir == "" {
		return errors.New("template directory is required")
	}
	if req.OverrideDir == "" {
		return errors.New("override directory is required")
	}
	switch req.OutputMode {
	case OutputModeFile:
		if req.OutputPath == "" {
			return errors.New("output path is required for file output mode")
		}
	case OutputModeK9s:
		if req.OutputPath != "" {
			return errors.New("output path cannot be used with K9s install mode")
		}
	default:
		if req.OutputPath != "" {
			return fmt.Errorf("unsupported output mode %q", req.OutputMode)
		}
		return errors.New("exactly one output mode is required: use --output or --install-to-k9s")
	}
	if req.OutputMode == OutputModeK9s && req.K9sDataDir != "" && filepath.Clean(req.K9sDataDir) == "." {
		return errors.New("K9s data directory must be a valid path")
	}
	return nil
}

func discoverSingleYAML(dir, label string) (string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", fmt.Errorf("read %s %q: %w", label, dir, err)
	}

	var matches []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		ext := filepath.Ext(entry.Name())
		if ext == ".yaml" || ext == ".yml" {
			matches = append(matches, filepath.Join(dir, entry.Name()))
		}
	}

	sort.Strings(matches)
	if len(matches) == 0 {
		return "", fmt.Errorf("%s %q must contain exactly one .yaml or .yml file, found none", label, dir)
	}
	if len(matches) > 1 {
		return "", fmt.Errorf("%s %q must contain exactly one .yaml or .yml file, found %d", label, dir, len(matches))
	}
	return matches[0], nil
}

func loadActiveCluster(path string) (string, error) {
	selection, err := loadClusterSelection(path)
	if err != nil {
		return "", err
	}
	return selection.ActiveCluster, nil
}

func loadClusterSelection(path string) (clusterSelection, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return clusterSelection{}, fmt.Errorf("read kubeconfig %q: %w", path, err)
	}

	var cfg kubeconfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return clusterSelection{}, fmt.Errorf("parse kubeconfig %q: %w", path, err)
	}

	if cfg.CurrentContext == "" {
		return clusterSelection{}, fmt.Errorf("kubeconfig %q is missing current-context", path)
	}

	clusterNames := make(map[string]struct{}, len(cfg.Clusters))
	for _, cluster := range cfg.Clusters {
		clusterNames[cluster.Name] = struct{}{}
	}

	var activeCluster string
	for _, ctx := range cfg.Contexts {
		if ctx.Name == cfg.CurrentContext {
			if ctx.Context.Cluster == "" {
				return clusterSelection{}, fmt.Errorf("kubeconfig context %q does not reference a cluster", cfg.CurrentContext)
			}
			if _, ok := clusterNames[ctx.Context.Cluster]; !ok {
				return clusterSelection{}, fmt.Errorf("kubeconfig context %q references missing cluster %q", cfg.CurrentContext, ctx.Context.Cluster)
			}
			activeCluster = ctx.Context.Cluster
			break
		}
	}

	if activeCluster == "" {
		return clusterSelection{}, fmt.Errorf("kubeconfig current-context %q does not match any context", cfg.CurrentContext)
	}

	contextNames := make([]string, 0, len(cfg.Contexts))
	for _, ctx := range cfg.Contexts {
		if ctx.Context.Cluster == activeCluster {
			contextNames = append(contextNames, ctx.Name)
		}
	}
	sort.Strings(contextNames)

	return clusterSelection{ActiveCluster: activeCluster, ContextNames: contextNames}, nil
}

func loadTemplate(path string) ([]byte, string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, "", fmt.Errorf("read template %q: %w", path, err)
	}

	pluginNames, err := discoverPluginNames(data)
	if err != nil {
		return nil, "", fmt.Errorf("inspect template %q: %w", path, err)
	}

	if len(pluginNames) != 1 {
		return nil, "", fmt.Errorf("template %q must define exactly one plugin under plugins, found %d", path, len(pluginNames))
	}

	return data, pluginNames[0], nil
}

func loadReplacements(path, pluginName, activeCluster string) (map[string]any, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read overrides %q: %w", path, err)
	}

	var doc overrideDocument
	if err := yaml.Unmarshal(data, &doc); err != nil {
		return nil, fmt.Errorf("parse overrides %q: %w", path, err)
	}

	override, ok := doc.PluginOverrides[pluginName]
	if !ok {
		return map[string]any{}, nil
	}

	for _, rule := range override.Clusters {
		matched, err := rule.Match.matches(activeCluster)
		if err != nil {
			return nil, err
		}
		if !matched {
			continue
		}
		if len(rule.Replace) == 0 {
			return nil, fmt.Errorf("override for plugin %q and cluster %q is missing replacement values", pluginName, activeCluster)
		}
		result := make(map[string]any, len(rule.Replace))
		for key, value := range rule.Replace {
			result[key] = value
		}
		return result, nil
	}

	return map[string]any{}, nil
}

func discoverPluginNames(data []byte) ([]string, error) {
	lines := strings.Split(string(data), "\n")
	inPlugins := false
	var names []string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}

		if !inPlugins {
			if trimmed == "plugins:" {
				inPlugins = true
			}
			continue
		}

		if strings.HasPrefix(line, " ") {
			if strings.HasPrefix(line, "  ") && !strings.HasPrefix(line, "    ") {
				key := strings.TrimSpace(line)
				if strings.HasSuffix(key, ":") {
					names = append(names, strings.TrimSuffix(key, ":"))
				}
				continue
			}
			continue
		}

		break
	}

	if len(names) == 0 {
		return nil, errors.New("template does not define any plugin keys under plugins")
	}

	return names, nil
}

func (m matchRule) matches(clusterName string) (bool, error) {
	switch m.Type {
	case "regex":
		if m.Value == "" {
			return false, fmt.Errorf("regex match rule is missing value")
		}
		re, err := regexp.Compile(m.Value)
		if err != nil {
			return false, fmt.Errorf("invalid regex match rule %q: %w", m.Value, err)
		}
		return re.MatchString(clusterName), nil
	case "list":
		if len(m.Values) == 0 {
			return false, fmt.Errorf("list match rule is missing values")
		}
		for _, value := range m.Values {
			if value == clusterName {
				return true, nil
			}
		}
		return false, nil
	default:
		return false, fmt.Errorf("unsupported match type %q", m.Type)
	}
}

func renderTemplate(data []byte, values map[string]any) ([]byte, error) {
	tpl, err := template.New("plugin").Funcs(sprig.TxtFuncMap()).Option("missingkey=default").Parse(string(data))
	if err != nil {
		return nil, fmt.Errorf("parse template expressions: %w", err)
	}

	var rendered bytes.Buffer
	if err := tpl.Execute(&rendered, values); err != nil {
		return nil, fmt.Errorf("render template: %w", err)
	}

	var doc any
	if err := yaml.Unmarshal(rendered.Bytes(), &doc); err != nil {
		return nil, fmt.Errorf("rendered output is not valid YAML: %w", err)
	}

	return rendered.Bytes(), nil
}

func resolveK9sDataRoot(explicit string) (string, error) {
	if explicit != "" {
		return explicit, nil
	}

	if env := os.Getenv("XDG_DATA_HOME"); env != "" {
		return env, nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve home directory for K9s data root: %w", err)
	}

	switch runtime.GOOS {
	case "darwin":
		return filepath.Join(home, "Library", "Application Support"), nil
	case "windows":
		if localAppData := os.Getenv("LocalAppData"); localAppData != "" {
			return localAppData, nil
		}
		return filepath.Join(home, "AppData", "Local"), nil
	default:
		return filepath.Join(home, ".local", "share"), nil
	}
}

func resolveK9sPluginPaths(dataRoot, cluster string, contexts []string) []string {
	paths := make([]string, 0, len(contexts))
	for _, contextName := range contexts {
		paths = append(paths, filepath.Join(dataRoot, "k9s", "clusters", cluster, contextName, "plugins.yaml"))
	}
	return paths
}

func installMergedPlugins(targetPaths []string, rendered []byte) error {
	generatedPlugins, err := parsePluginMap(rendered, "rendered output")
	if err != nil {
		return err
	}

	for _, targetPath := range targetPaths {
		existingPlugins, err := loadExistingPlugins(targetPath)
		if err != nil {
			return err
		}

		merged := mergePluginMaps(existingPlugins, generatedPlugins)
		content, err := marshalPluginDocument(merged)
		if err != nil {
			return fmt.Errorf("marshal merged plugins for %q: %w", targetPath, err)
		}

		if err := writeOutput(targetPath, content); err != nil {
			return err
		}
	}

	return nil
}

func loadExistingPlugins(path string) (map[string]any, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return map[string]any{}, nil
		}
		return nil, fmt.Errorf("read existing plugins %q: %w", path, err)
	}

	plugins, err := parsePluginMap(data, path)
	if err != nil {
		return nil, err
	}
	return plugins, nil
}

func parsePluginMap(data []byte, source string) (map[string]any, error) {
	if len(bytes.TrimSpace(data)) == 0 {
		return map[string]any{}, nil
	}

	var wrapped struct {
		Plugins map[string]any `yaml:"plugins"`
	}
	if err := yaml.Unmarshal(data, &wrapped); err == nil && wrapped.Plugins != nil {
		return wrapped.Plugins, nil
	}

	var raw map[string]any
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("parse plugin document %q: %w", source, err)
	}
	if raw == nil {
		return map[string]any{}, nil
	}

	pluginsValue, hasPlugins := raw["plugins"]
	if hasPlugins {
		plugins, ok := pluginsValue.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("plugin document %q has unsupported plugins shape", source)
		}
		return plugins, nil
	}

	for name, value := range raw {
		if _, ok := value.(map[string]any); !ok {
			return nil, fmt.Errorf("plugin document %q has unsupported top-level entry %q", source, name)
		}
	}

	return raw, nil
}

func mergePluginMaps(existing, generated map[string]any) map[string]any {
	merged := make(map[string]any, len(existing)+len(generated))
	for key, value := range existing {
		merged[key] = value
	}
	for key, value := range generated {
		merged[key] = value
	}
	return merged
}

func marshalPluginDocument(plugins map[string]any) ([]byte, error) {
	return yaml.Marshal(map[string]any{"plugins": plugins})
}

func writeOutput(path string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create output directory for %q: %w", path, err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("write output %q: %w", path, err)
	}
	return nil
}
