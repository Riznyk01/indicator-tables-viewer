package config

import (
	"github.com/pelletier/go-toml"
	"log"
	"os"
)

var configPath = "config/config_2.toml"

type Config struct {
	Host         string  `toml:"host"`
	Port         string  `toml:"port"`
	Path         string  `toml:"path"`
	Username     string  `toml:"username"`
	DBName       string  `toml:"db_name"`
	Password     string  `toml:"password"`
	HeaderHeight float32 `toml:"header_row_height"`
	RowHeight    float32 `toml:"row_height"`
	WindowHeight float32 `toml:"window_height"`
	WindowWidth  float32 `toml:"window_width"`
}

func MustLoad() *Config {

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

// UpdatePath updates the config file on the disk if it has been changed
func UpdatePath(config *Config, newPath string) error {
	config.Path = newPath

	file, err := os.OpenFile(configPath, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := toml.NewEncoder(file)
	if err = encoder.Encode(*config); err != nil {
		return err
	}

	return nil
}
