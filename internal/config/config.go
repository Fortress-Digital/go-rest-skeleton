package config

import (
	"flag"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type Application struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
	Env     string `yaml:"env"`
	Debug   bool   `yaml:"debug"`
}

type Server struct {
	Port         int `yaml:"port"`
	Timeout      int `yaml:"timeout"`
	ReadTimeout  int `yaml:"read_timeout"`
	WriteTimeout int `yaml:"write_timeout"`
}

type Database struct {
	Driver string `yaml:"driver"`
	Dsn    string `yaml:"dsn"`
}

type Supabase struct {
	Url string `yaml:"url"`
	Key string `yaml:"key"`
}

type Config struct {
	Application Application `yaml:"application"`
	Server      Server      `yaml:"server"`
	Database    Database    `yaml:"database"`
	Supabase    Supabase    `yaml:"supabase"`
}

type Configuration struct {
}

func (c Configuration) Register(app *App) error {
	cfg, err := newConfig()
	if err != nil {
		return err
	}

	app.Config = cfg

	return nil
}

func newConfig() (*Config, error) {
	configPath, err := parseConfig()
	if err != nil {
		return nil, err
	}

	config := &Config{}

	file, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	replaced := os.ExpandEnv(string(file))
	err = yaml.Unmarshal([]byte(replaced), config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func parseConfig() (string, error) {
	var configPath string

	flag.StringVar(&configPath, "config", "./config/config.yml", "path to config file")
	flag.Parse()

	if err := validateConfigPath(configPath); err != nil {
		return "", err
	}

	return configPath, nil
}

func validateConfigPath(path string) error {
	s, err := os.Stat(path)
	if err != nil {
		return err
	}
	if s.IsDir() {
		return fmt.Errorf("'%s' is a directory, not a normal file", path)
	}
	return nil
}
