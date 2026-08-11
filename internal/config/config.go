package config

import (
	"encoding/json"
	"os"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DbUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	configFilePath := homeDir + "/" + configFileName

	return configFilePath, nil
}

func Read() (Config, error) {
	configPath, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	configData, err := os.ReadFile(configPath)
	if err != nil {
		return Config{}, err
	}

	var config Config
	if err := json.Unmarshal(configData, &config); err != nil {
		return Config{}, err
	}

	return config, nil
}

func write(cfg *Config) error {
	// Get config path
	configPath, err := getConfigFilePath()
	if err != nil {
		return err
	}
	// Marshal config to JSON
	cfgJson, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	// Write file
	if err := os.WriteFile(configPath, cfgJson, 0644); err != nil {
		return err
	}

	return nil
}

func (cfg *Config) SetUser() error {
	if err := write(cfg); err != nil {
		return err
	}

	return nil
}
