package config

import (
	"encoding/json"
	"os"
)

const configFile = ".gatorconfig.json"

type Config struct {
	DbUrl       string `json:"db_url"`
	CurrentUser string `json:"current_user"`
}

func Read() (Config, error) {

	config := Config{}

	path, erro := getConfigFilePath()
	if erro != nil {
		return Config{}, erro
	}

	body, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	err = json.Unmarshal(body, &config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}

func (c Config) SetUser(user string) error {
	c.CurrentUser = user
	err := write(c)
	if err != nil {
		return err
	}

	return nil
}

func write(config Config) error {
	jsonData, err := json.Marshal(config)
	if err != nil {
		return err
	}

	path, erro := getConfigFilePath()
	if erro != nil {
		return erro
	}

	err = os.WriteFile(path, jsonData, 0666)
	if err != nil {
		return err
	}

	return nil

}

func getConfigFilePath() (string, error) {
	dir, erro := os.UserHomeDir()
	if erro != nil {
		return "", erro
	}

	return dir + "/" + configFile, nil
}
