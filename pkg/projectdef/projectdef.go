// Package projectdef provides the specscore-project.yaml schema and read/write operations.
package projectdef

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	SpecConfigFile    = "synchestra-spec-repo.yaml"
	StateConfigFile   = "synchestra-state-repo.yaml"
	CodeConfigFile    = "synchestra-code-repo.yaml"
	EmbeddedStateFile = "synchestra-state.yaml"
)

const worktreeScheme = "worktree://"

// PlanningConfig holds planning-related settings from synchestra-spec-repo.yaml.
type PlanningConfig struct {
	WhatsNext string `yaml:"whats_next"`
}

// SpecConfig represents the contents of synchestra-spec-repo.yaml.
type SpecConfig struct {
	Title     string          `yaml:"title"`
	StateRepo string          `yaml:"state_repo"`
	Repos     []string        `yaml:"repos"`
	Planning  *PlanningConfig `yaml:"planning,omitempty"`
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

// StateConfig represents the contents of synchestra-state-repo.yaml.
type StateConfig struct {
	Title     string   `yaml:"title"`
	MainRepo  string   `yaml:"main_repo"`
	SpecRepos []string `yaml:"spec_repos"`
	CodeRepos []string `yaml:"code_repos,omitempty"`
}

// CodeConfig represents the contents of synchestra-code-repo.yaml.
type CodeConfig struct {
	SpecRepos []string `yaml:"spec_repos"`
}

// EmbeddedStateConfig lives on the orphan branch.
type EmbeddedStateConfig struct {
	Title        string           `yaml:"title"`
	Mode         string           `yaml:"mode"`
	SourceBranch string           `yaml:"source_branch"`
	Sync         *EmbeddedSyncCfg `yaml:"sync,omitempty"`
}

// EmbeddedSyncCfg controls sync policy for embedded state.
type EmbeddedSyncCfg struct {
	Pull string `yaml:"pull"`
	Push string `yaml:"push"`
}

// WriteSpecConfig writes a SpecConfig to the given directory.
func WriteSpecConfig(dir string, cfg SpecConfig) error {
	return writeYAML(filepath.Join(dir, SpecConfigFile), cfg)
}

// WriteStateConfig writes a StateConfig to the given directory.
func WriteStateConfig(dir string, cfg StateConfig) error {
	return writeYAML(filepath.Join(dir, StateConfigFile), cfg)
}

// WriteCodeConfig writes a CodeConfig to the given directory.
func WriteCodeConfig(dir string, cfg CodeConfig) error {
	return writeYAML(filepath.Join(dir, CodeConfigFile), cfg)
}

// ReadSpecConfig reads a SpecConfig from the given directory.
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

// ReadStateConfig reads a StateConfig from the given directory.
func ReadStateConfig(dir string) (StateConfig, error) {
	var cfg StateConfig
	data, err := os.ReadFile(filepath.Join(dir, StateConfigFile))
	if err != nil {
		return cfg, fmt.Errorf("reading state config: %w", err)
	}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return cfg, fmt.Errorf("parsing state config: %w", err)
	}
	return cfg, nil
}

// ReadCodeConfig reads a CodeConfig from the given directory.
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

// WriteEmbeddedStateConfig writes an EmbeddedStateConfig to the given directory.
func WriteEmbeddedStateConfig(dir string, cfg EmbeddedStateConfig) error {
	return writeYAML(filepath.Join(dir, EmbeddedStateFile), cfg)
}

// ReadEmbeddedStateConfig reads an EmbeddedStateConfig from the given directory.
func ReadEmbeddedStateConfig(dir string) (EmbeddedStateConfig, error) {
	var cfg EmbeddedStateConfig
	data, err := os.ReadFile(filepath.Join(dir, EmbeddedStateFile))
	if err != nil {
		return cfg, fmt.Errorf("reading embedded state config: %w", err)
	}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return cfg, fmt.Errorf("parsing embedded state config: %w", err)
	}
	return cfg, nil
}

func writeYAML(path string, v any) error {
	data, err := yaml.Marshal(v)
	if err != nil {
		return fmt.Errorf("marshalling YAML: %w", err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("writing %s: %w", path, err)
	}
	return nil
}
