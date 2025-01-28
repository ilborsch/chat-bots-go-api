package config

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
	"time"
)

var configPath, _ = filepath.Abs("../../config/test.yaml")

func TestMustLoad(t *testing.T) {
	err := os.Setenv("CONFIG_PATH", configPath)
	require.NoError(t, err)
	defer os.Unsetenv("CONFIG_PATH")

	config := MustLoad()
	require.NotNil(t, config)

	assert.Equal(t, "test", config.Env)
	assert.Equal(t, 8080, config.Port)
	assert.Equal(t, "testuser", config.MySQLConfig.Username)
	assert.Equal(t, "testpass", config.MySQLConfig.Password)
	assert.Equal(t, 3306, config.MySQLConfig.Port)
	assert.Equal(t, "testsecret", config.SSOConfig.Secret)
	assert.Equal(t, 12345, config.SSOConfig.AppID)
	assert.Equal(t, 9090, config.SSOConfig.Port)
	assert.Equal(t, 30*time.Second, config.SSOConfig.Timeout)
}

func TestGetConfigPath(t *testing.T) {
	err := os.Setenv("CONFIG_PATH", configPath)
	require.NoError(t, err)

	gotPath := getConfigPath()
	assert.Equal(t, configPath, gotPath)
}

func TestMustLoadByPath(t *testing.T) {
	t.Run("empty path", func(t *testing.T) {
		assert.Panics(t, func() {
			MustLoadByPath("")
		})
	})

	t.Run("non-existent file", func(t *testing.T) {
		assert.Panics(t, func() {
			MustLoadByPath("nonexistent.yaml")
		})
	})

	t.Run("valid config path", func(t *testing.T) {
		config := MustLoadByPath(configPath)
		require.NotNil(t, config)

		assert.Equal(t, "test", config.Env)
		assert.Equal(t, 8080, config.Port)
		assert.Equal(t, "testuser", config.MySQLConfig.Username)
		assert.Equal(t, "testpass", config.MySQLConfig.Password)
		assert.Equal(t, 3306, config.MySQLConfig.Port)
		assert.Equal(t, "testsecret", config.SSOConfig.Secret)
		assert.Equal(t, 12345, config.SSOConfig.AppID)
		assert.Equal(t, 9090, config.SSOConfig.Port)
		assert.Equal(t, 30*time.Second, config.SSOConfig.Timeout)
	})
}
