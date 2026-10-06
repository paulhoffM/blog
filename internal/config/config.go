package config

import (
	"os"
	"encoding/json"
	"fmt"
	"path/filepath"
)

const configFileName = "/.gatorconfig.json"

type Config struct {
	DbUrl string `json:"db_url"`
	CurrentUserName string `json:"current_user_name` 
}


func (cfg *Config) SetUser(userName string) error{
	cfg.CurrentUserName = userName
	return write(*cfg)
}


func Read() (Config, error) {
	fullPath, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}
	
	fileStream, err := os.Open(fullPath)
	if err != nil{
		return Config{}, fmt.Errorf("Unable to read /.gatorconfig.json: %w", err)
	}
	defer  fileStream.Close()
	
	decoder := json.NewDecoder(fileStream)
	cfg := Config{}
	
	err = decoder.Decode(&cfg)
	if err != nil{
		return Config{}, fmt.Errorf("Could not create config struct: %w", err)
	}
	return cfg, nil
}

func getConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err!= nil {
		return "",err
	}
	fullPath := filepath.Join(home + configFileName)
	return fullPath, nil
}

func write(cfg Config) error{
	fullPath, err := getConfigFilePath()
		if err != nil {
			return err
		}
		
	file, err := os.Create(fullPath)
		if err != nil {
			return err
		}
		defer file.Close()
		
		encoder := json.NewEncoder(file)
	err = encoder.Encode(cfg)
		if err != nil {
			return err
		}
	return  nil
}