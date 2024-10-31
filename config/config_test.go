package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDefaultConfig(t *testing.T) {
	expectedConfig := Config{
		Aliases: make(map[string][]string),
	}

	defaultConfig := getDefaultConfig()

	assert.NotNil(t, defaultConfig)
	assert.ElementsMatch(t, expectedConfig.Aliases, defaultConfig.Aliases)
}

func TestGetConfigFilePath(t *testing.T) {
	expectedDir := "/tmp/test-config"
	t.Setenv("XDG_CONFIG_HOME", expectedDir)
	expectedPath := filepath.Join(expectedDir, "gh-pairing-with", "config.yml")

	configFilePath, err := getConfigFilePath()
	assert.Nil(t, err)
	assert.Equal(t, expectedPath, configFilePath)

	// Test when XDG_CONFIG_HOME is not set
	t.Setenv("XDG_CONFIG_HOME", "")
	homeDir, err := os.UserHomeDir()
	assert.Nil(t, err)

	expectedPath = filepath.Join(homeDir, ".config", "gh-pairing-with", "config.yml")
	configFilePath, err = getConfigFilePath()
	assert.Nil(t, err)
	assert.Equal(t, expectedPath, configFilePath)
}
