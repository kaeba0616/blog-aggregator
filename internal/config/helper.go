package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func getConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", nil
	}
	fullpath := filepath.Join(home, configFileName)
	return fullpath, nil
}

func write(cfg Config) error {
	fullpath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	file, err := os.Create(fullpath)
	if err != nil {
		return err
	}

	defer func() { _ = file.Close() }()

	encoder := json.NewEncoder(file)

	err = encoder.Encode(&cfg)
	if err != nil {
		return err
	}
	return nil
}
