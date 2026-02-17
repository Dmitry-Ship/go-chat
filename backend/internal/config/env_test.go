package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadDotEnvFindsEnvInParentDirectory(t *testing.T) {
	const key = "CODEX_TEST_PARENT_ENV"
	const value = "from-parent"

	restoreEnv(t, key)
	restoreEnv(t, "ENV_FILE")
	require.NoError(t, os.Unsetenv(key))
	require.NoError(t, os.Unsetenv("ENV_FILE"))

	projectDir := t.TempDir()
	backendDir := filepath.Join(projectDir, "backend")
	require.NoError(t, os.MkdirAll(backendDir, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(projectDir, ".env"), []byte(key+"="+value+"\n"), 0o644))

	restoreWorkingDir(t)
	require.NoError(t, os.Chdir(backendDir))

	require.NoError(t, LoadDotEnv())
	assert.Equal(t, value, os.Getenv(key))
}

func TestLoadDotEnvUsesExplicitEnvFilePath(t *testing.T) {
	const key = "CODEX_TEST_EXPLICIT_ENV"
	const value = "from-explicit"

	restoreEnv(t, key)
	restoreEnv(t, "ENV_FILE")
	require.NoError(t, os.Unsetenv(key))

	envFilePath := filepath.Join(t.TempDir(), ".custom.env")
	require.NoError(t, os.WriteFile(envFilePath, []byte(key+"="+value+"\n"), 0o644))
	require.NoError(t, os.Setenv("ENV_FILE", envFilePath))

	require.NoError(t, LoadDotEnv())
	assert.Equal(t, value, os.Getenv(key))
}

func TestLoadDotEnvReturnsNilWhenNoEnvFileExists(t *testing.T) {
	restoreEnv(t, "ENV_FILE")
	require.NoError(t, os.Unsetenv("ENV_FILE"))

	emptyDir := t.TempDir()
	restoreWorkingDir(t)
	require.NoError(t, os.Chdir(emptyDir))

	require.NoError(t, LoadDotEnv())
}

func restoreWorkingDir(t *testing.T) {
	t.Helper()

	currentDir, err := os.Getwd()
	require.NoError(t, err)

	t.Cleanup(func() {
		require.NoError(t, os.Chdir(currentDir))
	})
}

func restoreEnv(t *testing.T, key string) {
	t.Helper()

	previousValue, hadValue := os.LookupEnv(key)
	t.Cleanup(func() {
		var err error
		if hadValue {
			err = os.Setenv(key, previousValue)
		} else {
			err = os.Unsetenv(key)
		}
		require.NoError(t, err)
	})
}
