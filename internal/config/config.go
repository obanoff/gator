package config

import (
	"database/sql"
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/obanoff/gator/internal/database"
)

type Database struct {
	Queries *database.Queries
	DB      *sql.DB
}

type State struct {
	Config *Config
	Logger *Logger
	DB     Database
}

type Config struct {
	DBUrl    string `json:"db_url"`
	Username string `json:"current_user_name"`
}

func Read() (Config, error) {
	path, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	file, err := os.Open(path)
	if err != nil {
		return Config{}, nil
	}
	defer file.Close()

	var config Config
	err = json.NewDecoder(file).Decode(&config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}

func (cfg *Config) SetUser(username string) error {
	path, err := getConfigFilePath()
	if err != nil {
		return err
	}

	cfg.Username = username

	data, err := json.MarshalIndent(&cfg, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(path, data, 0200)
	if err != nil {
		return err
	}

	return nil
}

func getConfigFilePath() (string, error) {
	const configFileName = ".gatorconfig.json"

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	path := filepath.Join(homeDir, configFileName)

	return path, nil
}
