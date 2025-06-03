package config

import (
	"encoding/json"
	"fmt"
	"os"
)

const config_file = ".gatorconfig.json"

type Config struct {
	Url           string `json:"db_url"`
	Curr_username string `json:"current_user_name"`
}

func Read() (*Config, error) {
	filepath, err := getConfigFilePath()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("couldnt read from config file: %w", err)
	}
	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("couldnt unmarshal json into Config struct: %w", err)
	}
	return &config, nil
}
func (cfg Config) SetUser(username string) error {
	cfg.Curr_username = username
	err := write(cfg)
	if err != nil {
		return err
	}
	return nil
}

func write(cfg Config) error {
	marshalled, err := json.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("error marshaling Config struct into json: %w", err)
	}
	filepath, err := getConfigFilePath()
	if err != nil {
		return fmt.Errorf("couldnt file .gatorconfig.json on home directory")
	}
	err = os.WriteFile(filepath, marshalled, 666)
	if err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}
	return nil
}

func getConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("couldnt find home directory: %w", err)
	}

	filepath := home + "/" + config_file
	_, err2 := os.Stat(filepath)
	if err2 != nil {
		return "", fmt.Errorf("couldnt find the .gatorconfig.json file in the home directory")
	}
	return filepath, nil
}
