// This file is safe to edit. Once it exists it will not be overwritten
package config

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// Load - loads config from yamlfile
func Load(prefix string) (*Config, error) {
	// config with default values
	c := Config{}

	// get env name
	if len(prefix) > 0 {
		prefix = strings.ToUpper(prefix) + "_"
	}

	configFilePathEnvName := prefix + "CONFIG_PATH"

	// get config file path from env
	configFilePath := os.Getenv(configFilePathEnvName)

	// if not set, use default config file is config.yml
	if len(configFilePath) == 0 {
		configFilePath = "config.yml"
		slog.Info("env value is empty, set config file path as default ", "env name", configFilePathEnvName, "value", configFilePath)
	} else {
		slog.Info("get config file path from env", "env name", configFilePathEnvName, "value", configFilePath)
	}

	_, err := os.Stat(configFilePath)
	if os.IsNotExist(err) {
		slog.Info("config file not found, using config with default values")

		// save config with default values
		errSave := Save(&c, configFilePath)
		if errSave != nil {
			return nil, fmt.Errorf("config.Save: %w", errSave)
		}

		return &c, nil
	}

	data, err := os.ReadFile(configFilePath)
	if err != nil {
		return nil, fmt.Errorf("os.ReadFile: %w, file: %s", err, configFilePath)
	}

	err = yaml.Unmarshal(data, &c)
	if err != nil {
		return nil, fmt.Errorf("yaml.Unmarshal: %w, file: %s", err, configFilePath)
	}

	slog.Info("config loaded from file", "file path", configFilePath)

	return &c, nil
}

// Save - saves config to yaml file
func Save(config *Config, filename string) error {
	if len(filename) == 0 {
		return errors.New("filename cannot be empty")
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		return err
	}

	slog.Info("config saved to file", "file path", filename)

	return nil
}
