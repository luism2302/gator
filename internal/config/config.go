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

func Read() (Config, error) {
	configPath, err := getConfigFilePath()
	data, err := os.ReadFile(configPath)
	if err != nil {
		return Config{}, fmt.Errorf("error: couldnt read config file: %w", err)
	}
	var cfg Config
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return Config{}, fmt.Errorf("error: couldnt unmarshal into Config struct: %w", err)
	}
	return cfg, nil
}

func (cfg *Config) SetUser(username string) error {
	cfg.CurrentUserName = username
	err := write(*cfg)
	if err != nil {
		return err
	}
	return nil
}

func getConfigFilePath() (string, error) {
	usrHome, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("error: couldnt get users home directory: %w", err)
	}
	configPath := usrHome + "/" + configFileName
	return configPath, nil
}

func write(cfg Config) error {
	marshaled, err := json.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("error: couldnt marshal into json: %w", err)
	}
	configPath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	err = os.WriteFile(configPath, marshaled, 0644)
	if err != nil {
		return fmt.Errorf("error: couldnt write to config file: %w", err)
	}
	return nil
}
