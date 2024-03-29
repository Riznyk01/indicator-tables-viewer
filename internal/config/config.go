package config

import (
	"github.com/pelletier/go-toml"
	"log"
)

type Config struct {
	Host     string `toml:"host"`
	Port     string `toml:"port"`
	Path     string `toml:"path"`
	Username string `toml:"username"`
	DBName   string `toml:"db_name"`
	Password string `toml:"password"`
}

func MustLoad() *Config {

	configPath := "config/config_1.toml"
	// check if file exists
	cfg, err := toml.LoadFile(configPath)
	if err != nil {
		log.Fatalf("error loading config file: %s", err)
	}

	var config Config

	if err := cfg.Unmarshal(&config); err != nil {
		log.Fatalf("error decoding config: %s", err)
	}

	return &config
}
