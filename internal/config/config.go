package config

import (
  "encoding/json"
  "os"
  "path"
)

const configFileName = ".gatorconfig.json"

type Config struct {
  DBUrl string `json:"db_url"`
  CurrentUserName string `json:"current_user_name"`
}

func getConfigFilePath() (string, error) {
  home, err := os.UserHomeDir()
  if err != nil {
    return "", err
  }
  return path.Join(home, configFileName), nil
}

func Read() (Config, error) {
  config := Config{}

  path, err := getConfigFilePath()
  if err != nil {
    return Config{}, err
  }

  dat, err := os.ReadFile(path)
  if err != nil {
    return Config{}, err
  }
  if err := json.Unmarshal(dat, &config); err != nil {
    return Config{}, err
  }
    
  return config, nil
}

func (c Config) SetUser(name string) error {
  c.CurrentUserName = name
  
  err := write(c)
  if err != nil {
    return err
  }
  return nil
}

func write(c Config) error {
  dat, err := json.Marshal(c)

  path, err := getConfigFilePath()
  if err != nil {
    return err
  }
  
  if err := os.WriteFile(path, dat, 0666); err != nil {
    return err
  }
  return nil
}
