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
		return Config{}, fmt.Errorf("error reading from config file: %w", err)
	}
	var cfg Config
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return Config{}, fmt.Errorf("error unmarshaling into Config struct: %w", err)
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
		return "", fmt.Errorf("error getting users home directory: %w", err)
	}
	configPath := usrHome + "/" + configFileName
	return configPath, nil
}

func write(cfg Config) error {
	marshaled, err := json.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("error marshaling struct to json: %w", err)
	}
	configPath, err := getConfigFilePath()
	if err != nil {
		return err
	}
	file, err := os.OpenFile(configPath, os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("error opening config file: %w", err)
	}
	defer file.Close()
	_, err = file.Write(marshaled)

	if err != nil {
		return fmt.Errorf("error writing to config file: %w", err)
	}
	return nil
}
