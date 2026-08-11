package config

import (
	"encoding/json"
	"fmt"
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

func Read() Config {
	configPath, err := getConfigFilePath()
	if err != nil {
		panic(err)
	}

	fmt.Println(configPath)
	configData, err := os.ReadFile(configPath)
	if err != nil {
		panic(err)
	}

	var config Config
	if err := json.Unmarshal(configData, &config); err != nil {
		panic(err)
	}

	return config
}
