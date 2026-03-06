package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

const DefaultConfigFileName = ".genie.yaml"

type Config struct {
	ProtectedBranch string   `yaml:"protected_branch"`
	RequireScope    bool     `yaml:"require_scope"`
	EnforceAll      bool     `yaml:"enforce_all"`
	AllowedModules  []string `yaml:"allowed_modules,omitempty"`
	Types           []string `yaml:"types,omitempty"`
}

func NewDefaultConfig() *Config {
	return &Config{
		ProtectedBranch: "master",
		RequireScope:    false,
		EnforceAll:      false,
		Types:           []string{"feat", "fix", "docs", "style", "refactor", "perf", "test", "chore", "revert"},
	}
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (c *Config) Save(path string) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}
