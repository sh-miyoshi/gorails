package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type GlobalConfig struct {
	GoModulePath string `yaml:"go-module-path"`
}

var inst GlobalConfig

func Init(filePath string) error {
	fp, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open config file: %w", err)
	}
	defer fp.Close()

	if err := yaml.NewDecoder(fp).Decode(&inst); err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}

	// TODO validation

	return nil
}

func Get() *GlobalConfig {
	return &inst
}
