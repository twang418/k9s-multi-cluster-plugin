package generator

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"gopkg.in/yaml.v3"
)

type Request struct {
	KubeconfigPath string
	TemplateDir    string
	OverrideDir    string
	OutputPath     string
}

type Result struct {
	ActiveCluster string
	OutputPath    string
}

type kubeconfig struct {
	CurrentContext string              `yaml:"current-context"`
	Contexts       []namedContext      `yaml:"contexts"`
	Clusters       []namedCluster      `yaml:"clusters"`
	Users          []namedUser         `yaml:"users"`
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

func Generate(req Request) (Result, error) {
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

	activeCluster, err := loadActiveCluster(req.KubeconfigPath)
	if err != nil {
		return Result{}, err
	}

	templateBytes, pluginName, err := loadTemplate(templatePath)
	if err != nil {
		return Result{}, err
	}

	replacements, err := loadReplacements(overridePath, pluginName, activeCluster)
	if err != nil {
		return Result{}, err
	}

	rendered, err := renderTemplate(templateBytes, replacements)
	if err != nil {
		return Result{}, err
	}

	if err := writeOutput(req.OutputPath, rendered); err != nil {
		return Result{}, err
	}

	return Result{ActiveCluster: activeCluster, OutputPath: req.OutputPath}, nil
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
	if req.OutputPath == "" {
		return errors.New("output path is required")
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
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("read kubeconfig %q: %w", path, err)
	}

	var cfg kubeconfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return "", fmt.Errorf("parse kubeconfig %q: %w", path, err)
	}

	if cfg.CurrentContext == "" {
		return "", fmt.Errorf("kubeconfig %q is missing current-context", path)
	}

	for _, ctx := range cfg.Contexts {
		if ctx.Name == cfg.CurrentContext {
			if ctx.Context.Cluster == "" {
				return "", fmt.Errorf("kubeconfig context %q does not reference a cluster", cfg.CurrentContext)
			}
			for _, cluster := range cfg.Clusters {
				if cluster.Name == ctx.Context.Cluster {
					return cluster.Name, nil
				}
			}
			return "", fmt.Errorf("kubeconfig context %q references missing cluster %q", cfg.CurrentContext, ctx.Context.Cluster)
		}
	}

	return "", fmt.Errorf("kubeconfig current-context %q does not match any context", cfg.CurrentContext)
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

func writeOutput(path string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create output directory for %q: %w", path, err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("write output %q: %w", path, err)
	}
	return nil
}
