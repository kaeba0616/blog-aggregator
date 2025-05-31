package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	DBURL           string `json:"db_url"`
	CurrentUserName string `json:"Current_User_Name"`
}

func (c *Config) SetUser(user_name string) error {
	c.CurrentUserName = user_name
	return write(*c)
}

func Read() (Config, error) {
	fullPath, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	file, err := os.Open(fullPath)
	if err != nil {
		return Config{}, err
	}

	defer func() { _ = file.Close() }()

	decoder := json.NewDecoder(file)
	cfg := Config{}

	err = decoder.Decode(&cfg)
	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}
