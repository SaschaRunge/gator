package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	Db_url            string `json:"db_url"`
	Current_user_name string `json:"current_user_name"`
}

func Read() (Config, error) {
	cfgPath, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	data, err := os.ReadFile(cfgPath)
	if err != nil {
		return Config{}, err
	}

	cfg := Config{}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func (c *Config) SetUser(user string) error {
	c.Current_user_name = user
	return write(*c)
}

func write(cfg Config) error {
	data, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	filePath, err := getConfigFilePath()
	if err != nil {
		return err
	}
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return err
	}
	return nil
}

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	cfgPath := filepath.Join(homeDir, configFileName)
	return cfgPath, err
}
