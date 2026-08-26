package config

import (
	"encoding/json"
	"errors"
	"os"
)

const configFileName string = ".rssconfig.json"

var home, _ = os.UserHomeDir()

type Config struct {
	Db_url           string `json:"db_url"`
	Current_username string `json:"current_username"`
}

func Read() (Config, error) {
	configBytes, err := os.ReadFile(home + "/" + configFileName)
	if err != nil {
		return Config{}, errors.New("Can't read config file")
	}
	var data Config
	json.Unmarshal(configBytes, &data)
	return data, nil
}

func (c Config) SetUser() error {
	configBytes, err := json.Marshal(c)
	if err != nil {
		return errors.New("Config file corrupted")
	}
	err = os.WriteFile(home+"/"+configFileName, configBytes, 0600)
	return err
}
