package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

const envFileName = ".env"

// LoadDotEnv loads environment variables from a .env file.
// It first checks ENV_FILE, then searches current and parent directories.
func LoadDotEnv() error {
	if envFilePath := os.Getenv("ENV_FILE"); envFilePath != "" {
		if err := godotenv.Load(envFilePath); err != nil {
			return fmt.Errorf("load env file %q: %w", envFilePath, err)
		}

		return nil
	}

	envFilePath, err := findDotEnvPath()
	if err != nil {
		return err
	}
	if envFilePath == "" {
		return nil
	}

	if err := godotenv.Load(envFilePath); err != nil {
		return fmt.Errorf("load env file %q: %w", envFilePath, err)
	}

	return nil
}

func findDotEnvPath() (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("get working directory: %w", err)
	}

	for {
		candidate := filepath.Join(currentDir, envFileName)
		_, statErr := os.Stat(candidate)
		if statErr == nil {
			return candidate, nil
		}

		if !errors.Is(statErr, os.ErrNotExist) {
			return "", fmt.Errorf("check env file %q: %w", candidate, statErr)
		}

		parentDir := filepath.Dir(currentDir)
		if parentDir == currentDir {
			return "", nil
		}
		currentDir = parentDir
	}
}
