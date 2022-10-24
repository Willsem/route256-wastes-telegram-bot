package startup

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type MigrationConfig struct {
	Database DatabaseConfig `yaml:"database"`
}

func NewMigrationConfig(configFile string) (*MigrationConfig, error) {
	rawYAML, err := os.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("reading file error: %w", err)
	}

	cfg := &MigrationConfig{}
	if err = yaml.Unmarshal(rawYAML, cfg); err != nil {
		return nil, fmt.Errorf("yaml parsing error: %w", err)
	}

	return cfg, nil
}
