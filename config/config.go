package config

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Aliases map[string][]string
}

func createConfigFileIfMissing(configFilePath string) error {
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		newConfigFile, err := os.OpenFile(
			configFilePath,
			os.O_RDWR|os.O_CREATE|os.O_EXCL,
			0666,
		)
		if err != nil {
			return err
		}

		var emptyConfig Config
		emptyConfig.Aliases = make(map[string][]string)

		blankConfigFile, err := yaml.Marshal(emptyConfig)
		if err != nil {
			return err
		}

		_, err = io.Writer.Write(newConfigFile, blankConfigFile)
		if err != nil {
			return err
		}

		defer newConfigFile.Close()
		return nil
	}

	return nil
}

func getConfigFilePath() (string, error) {
	const PairingWithDir = "gh-pairing-with"
	const ConfigYmlFileName = "config.yml"
	const DEFAULT_XDG_CONFIG_DIRNAME = ".config"

	configDir := os.Getenv("XDG_CONFIG_HOME")

	if configDir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		configDir = filepath.Join(homeDir, DEFAULT_XDG_CONFIG_DIRNAME)
	}

	pairingWithConfigDir := filepath.Join(configDir, PairingWithDir)
	return filepath.Join(pairingWithConfigDir, ConfigYmlFileName), nil
}

func LoadConfig() (*Config, error) {
	var config Config

	configFilePath, err := getConfigFilePath()

	if err != nil {
		return nil, err
	}

	configDir := filepath.Dir(configFilePath)
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		if err = os.MkdirAll(configDir, os.ModePerm); err != nil {
			return &config, err
		}
	}

	if err := createConfigFileIfMissing(configFilePath); err != nil {
		return &config, err
	}

	existingFile, err := os.ReadFile(configFilePath)
	if err != nil {
		return &config, err
	}

	err = yaml.Unmarshal(existingFile, &config)
	if err != nil {
		return &config, err
	}

	return &config, nil
}

func (c *Config) Save() error {
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	updatedFile, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	f, err := os.Create(configFilePath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Writer.Write(f, updatedFile)
	if err != nil {
		return fmt.Errorf("could not write to config file: %w", err)
	}

	return nil
}
