package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func (cfg *Config) SetUser(name string) error {
	cfg.CurrentUserName = name
	return write(*cfg)
}

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to detect home directory %w", err)
	}

	return filepath.Join(homeDir, configFileName), nil
}

func Read() (Config, error) {

	fileConfig, _ := getConfigFilePath()

	data, err := os.ReadFile(fileConfig)
	if err != nil {
		return Config{}, fmt.Errorf("failed to read config file %q: %w", fileConfig, err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return Config{}, err
	}

	return config, nil
}

func write(cfg Config) error {
	fileConfig, err := getConfigFilePath()
	if err != nil {
		return err
	}

	jsonData, err := json.MarshalIndent(cfg, "", " ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err = os.WriteFile(fileConfig, jsonData, 0600); err != nil {
		return fmt.Errorf("failed to write config file %q: %w", fileConfig, err)
	}

	return nil
}
