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

func TestAliasExists(t *testing.T) {
	config := Config{
		Aliases: map[string][]string{
			"alias1": {"user1", "user2"},
			"alias2": {"user3", "user4"},
			"alias3": {"user5", "user6"},
		},
	}

	assert.True(t, config.AliasExists("alias1"))
	assert.True(t, config.AliasExists("alias2"))
	assert.True(t, config.AliasExists("alias3"))
	assert.False(t, config.AliasExists("alias4"))
}

func TestGetAllAliases(t *testing.T) {
	config := Config{
		Aliases: map[string][]string{
			"alias1": {"user1", "user2"},
			"alias2": {"user3", "user4"},
			"alias3": {"user5", "user6"},
		},
	}

	assert.Equal(t,
		config.Aliases,
		config.GetAllAliases(),
	)
}

func TestExpandHandles(t *testing.T) {
	config := Config{
		Aliases: map[string][]string{
			"alias1": {"user1", "user2"},
			"alias2": {"user3", "user4"},
			"alias3": {"user5", "user6"},
		},
	}

	var handles []string

	handles = []string{"alias1", "alias2", "user5"}
	assert.ElementsMatch(t,
		[]string{"user1", "user2", "user3", "user4", "user5"},
		config.ExpandHandles(handles),
	)

	handles = []string{"user9"}
	assert.ElementsMatch(t,
		[]string{"user9"},
		config.ExpandHandles(handles),
	)

	handles = []string{"alias1", "user1"}
	assert.ElementsMatch(t,
		[]string{"user1", "user2"},
		config.ExpandHandles(handles),
	)
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
