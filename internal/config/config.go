package config

import (
	"github.com/pelletier/go-toml"
	"log"
	"os"
	"time"
)

type Config struct {
	Username          string        `toml:"username"`
	Host              string        `toml:"host"`
	Port              string        `toml:"port"`
	Path              string        `toml:"path"`
	Password          string        `toml:"password"`
	LocalHost         string        `toml:"host_local"`
	LocalPort         string        `toml:"port_local"`
	LocalPath         string        `toml:"path_local"`
	LocalPassword     string        `toml:"password_local"`
	DBName            string        `toml:"db_name"`
	WindowHeight      float32       `toml:"window_height"`
	WindowWidth       float32       `toml:"window_width"`
	InfoTimeout       time.Duration `toml:"info_timeout"`
	GoToUpdateTimeout time.Duration `toml:"go_to_update_timeout"`
	LocalMode         bool          `toml:"local_mode"`
	UpdatePath        string        `toml:"update_path"`
	Ver               string        `toml:"update_version"`
	AutoUpdate        bool          `toml:"auto_update"`
	XlsExportPath     string        `toml:"excel_export_path"`
	IconPath          string        `toml:"icon_path"`
	VerFilePath       string        `toml:"path_to_ver_file"`
	RemoteExeFilename string        `toml:"remote_exe_filename"`
	LocalExeFilename  string        `toml:"local_exe_filename"`
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
