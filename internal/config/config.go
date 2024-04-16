package config

import (
	"github.com/pelletier/go-toml"
	"log"
	"os"
	"time"
)

type Config struct {
	Username      string        `toml:"username"`
	Host          string        `toml:"host"`
	Port          string        `toml:"port"`
	Path          string        `toml:"path"`
	Password      string        `toml:"password"`
	LocalHost     string        `toml:"host_local"`
	LocalPort     string        `toml:"port_local"`
	LocalPath     string        `toml:"path_local"`
	LocalPassword string        `toml:"password_local"`
	DBName        string        `toml:"db_name"`
	HeaderHeight  float32       `toml:"header_row_height"`
	RowHeight     float32       `toml:"row_height"`
	WindowHeight  float32       `toml:"window_height"`
	WindowWidth   float32       `toml:"window_width"`
	InfoTimeout   time.Duration `toml:"info_timeout"`
	LocalMode     bool          `toml:"local_mode"`
}

func MustLoad(configPath string) *Config {
	cfg, err := toml.LoadFile(configPath)
	if err != nil {
		log.Fatalf("error loading config file: %s", err)
	}

	var config Config

	if err := cfg.Unmarshal(&config); err != nil {
		log.Fatalf("error decoding config: %s", err)
	}

	log.Printf("the config values: %v", config)

	return &config
}

// UpdateConfig updates the config file on the disk
func UpdateConfig(config Config, configPath string) error {

	file, err := os.OpenFile(configPath, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := toml.NewEncoder(file)
	if err = encoder.Encode(config); err != nil {
		return err
	}

	return nil
}

func SaveLocalModeCheckboxState(config Config, configPath string) error {
	file, err := os.OpenFile(configPath, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := toml.NewEncoder(file)
	if err = encoder.Encode(config); err != nil {
		return err
	}

	return nil
}
