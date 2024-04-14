package hive

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Hive struct {
	dir        string
	configName string
	config     *Config
}

func New(dir string, configName string) (*Hive, error) {

	raw, err := os.ReadFile(filepath.Join(dir, ".load-hive", fmt.Sprintf("%v.yaml", configName)))
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(raw, &cfg); err != nil {
		return nil, err
	}

	project := Hive{
		dir:        dir,
		configName: configName,
		config:     &cfg,
	}

	if project.config.Hive.UUID == "" {
		return nil, fmt.Errorf("config hive uuid not set")
	}

	if project.config.Hive.Port == "" {
		return nil, fmt.Errorf("config hive port not set")
	}

	return &project, nil
}

func (receiver *Hive) MainDir() string {
	return fmt.Sprintf("/var/lib/load-hive-%v-main", receiver.config.Hive.UUID)
}

func (receiver *Hive) NodeDir() string {
	return fmt.Sprintf("/var/lib/load-hive-%v-node", receiver.config.Hive.UUID)
}
