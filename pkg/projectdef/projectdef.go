package projectdef

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	SpecConfigFile = "specscore-spec-repo.yaml"
	CodeConfigFile = "specscore-code-repo.yaml"
)

const worktreeScheme = "worktree://"

type PlanningConfig struct {
	WhatsNext string `yaml:"whats_next"`
}

// HubConfig controls optional Synchestra Hub integration. Presence of Host
// (non-empty) is the opt-in signal; when omitted, no Hub-related behavior
// (e.g. the hub-view-link lint rule) activates.
type HubConfig struct {
	Host string `yaml:"host"`
}

type SpecConfig struct {
	Title     string          `yaml:"title"`
	StateRepo string          `yaml:"state_repo"`
	Repos     []string        `yaml:"repos"`
	Planning  *PlanningConfig `yaml:"planning,omitempty"`
	Hub       *HubConfig      `yaml:"hub,omitempty"`
	Extras    map[string]any  `yaml:",inline"`
}

// HubHost returns the configured Synchestra Hub base URL, or "" if Hub
// integration is not enabled for this project.
func (c SpecConfig) HubHost() string {
	if c.Hub == nil {
		return ""
	}
	return strings.TrimRight(c.Hub.Host, "/")
}

// WhatsNextMode returns the effective whats_next mode, defaulting to "disabled".
func (c SpecConfig) WhatsNextMode() string {
	if c.Planning != nil && c.Planning.WhatsNext != "" {
		return c.Planning.WhatsNext
	}
	return "disabled"
}

// ParseStateRepo parses the state_repo field.
// Returns (mode, branch):
//   - ("worktree", branchName) for "worktree://branchName"
//   - ("repo", "") for any other non-empty value
//   - ("", "") if state_repo is empty
func (c SpecConfig) ParseStateRepo() (mode, branch string) {
	if c.StateRepo == "" {
		return "", ""
	}
	if strings.HasPrefix(c.StateRepo, worktreeScheme) {
		b := c.StateRepo[len(worktreeScheme):]
		if b == "" {
			return "", ""
		}
		return "worktree", b
	}
	return "repo", ""
}

type CodeConfig struct {
	SpecRepos []string `yaml:"spec_repos"`
}

func WriteSpecConfig(dir string, cfg SpecConfig) error {
	return WriteYAML(filepath.Join(dir, SpecConfigFile), cfg)
}

func WriteCodeConfig(dir string, cfg CodeConfig) error {
	return WriteYAML(filepath.Join(dir, CodeConfigFile), cfg)
}

func ReadSpecConfig(dir string) (SpecConfig, error) {
	var cfg SpecConfig
	data, err := os.ReadFile(filepath.Join(dir, SpecConfigFile))
	if err != nil {
		return cfg, fmt.Errorf("reading spec config: %w", err)
	}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return cfg, fmt.Errorf("parsing spec config: %w", err)
	}
	return cfg, nil
}

func ReadCodeConfig(dir string) (CodeConfig, error) {
	var cfg CodeConfig
	data, err := os.ReadFile(filepath.Join(dir, CodeConfigFile))
	if err != nil {
		return cfg, fmt.Errorf("reading code config: %w", err)
	}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return cfg, fmt.Errorf("parsing code config: %w", err)
	}
	return cfg, nil
}

// WriteYAML marshals v to YAML and writes it to path.
func WriteYAML(path string, v any) error {
	data, err := yaml.Marshal(v)
	if err != nil {
		return fmt.Errorf("marshalling YAML: %w", err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("writing %s: %w", path, err)
	}
	return nil
}
