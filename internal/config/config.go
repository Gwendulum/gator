package config

import (
	"encoding/json"
	"fmt"
	"os"
)

const filename = "/.gatorconfig.json"

type Config struct {
	DBUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func (cfg Config) SetUser(user string) error {
	cfg.CurrentUserName = user
	err := write(cfg)
	if err != nil {
		return fmt.Errorf("SetUser: %v", err)
	}
	return nil
}

func Read() (Config, error) {
	homepath, err := os.UserHomeDir()
	if err != nil {
		return Config{}, fmt.Errorf("Homedir: %v", err)
	}

	filepath := homepath + filename

	file, err := os.Open(filepath)

	if err != nil {
		return Config{}, fmt.Errorf("os.Open: %v", err)
	}
	defer file.Close()

	var cfg Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return Config{}, fmt.Errorf("Decoder: %v", err)
	}

	return cfg, nil
}

func write(cfg Config) error {
	homepath, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	filepath := homepath + filename

	jsonCfg, err := json.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("json.Marshal: %v", err)
	}
	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("os.Create: %v", err)
	}
	defer file.Close()

	file.Write(jsonCfg)

	return nil

}
