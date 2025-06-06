package config

import (
	"encoding/json"
	"fmt"
	"os"
)

const config_file = ".gatorconfig.json"

type Config struct {
	DbUrl         string `json:"db_url"`
	Curr_username string `json:"current_user_name"`
}

func Read() (Config, error) {
	filepath, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}
	file, err := os.Open(filepath)
	if err != nil {
		return Config{}, fmt.Errorf("couldnt open config file: %w", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	cfg := Config{}
	err = decoder.Decode(&cfg)
	if err != nil {
		return Config{}, fmt.Errorf("couldnt decode json into Config struct: %w", err)
	}
	return cfg, nil
}
func (cfg *Config) SetUser(username string) error {
	cfg.Curr_username = username
	return write(*cfg)
}

func write(cfg Config) error {
	filepath, err := getConfigFilePath()
	if err != nil {
		return err
	}
	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	err = encoder.Encode(cfg)
	if err != nil {
		return fmt.Errorf("error encoding Config struct into json")
	}
	return nil
}

func getConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("couldnt find home directory: %w", err)
	}
	filepath := home + "/" + config_file
	return filepath, nil
}

func GetLoggedUserName() (string, error) {
	cfg, err := Read()
	if err != nil {
		return "", err
	}
	return cfg.Curr_username, nil

}
