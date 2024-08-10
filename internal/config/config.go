package config

import (
	"github.com/pelletier/go-toml"
	"log"
	"os"
	"time"
)

type Config struct {
	Username            string        `toml:"username"`
	Host                string        `toml:"host"`
	Port                string        `toml:"port"`
	RemotePathToDb      string        `toml:"remote_path_to_db"`
	LocalHost           string        `toml:"host_local"`
	LocalPort           string        `toml:"port_local"`
	LocalPath           string        `toml:"local_path"`
	DBName              string        `toml:"db_name"`
	WindowHeight        float32       `toml:"window_height"`
	WindowWidth         float32       `toml:"window_width"`
	InfoTimeout         time.Duration `toml:"info_timeout"`
	LocalMode           bool          `toml:"local_mode"`
	UpdateURL           string        `toml:"update_url"`
	AutoUpdate          bool          `toml:"auto_update"`
	XlsExportPath       string        `toml:"excel_export_path"`
	IconPath            string        `toml:"icon_path"`
	VerRemoteFilePath   string        `toml:"path_to_remote_ver_file"`
	VerLocalFilePath    string        `toml:"path_to_ver_file"`
	UpdateArch          string        `toml:"update_arch"`
	LocalExeFilename    string        `toml:"local_exe_filename"`
	LauncherExeFilename string        `toml:"launcher_exe_filename"`
	LogFileExt          string        `toml:"log_file_extension"`
	DownloadedVerFile   string        `toml:"downloaded_ver_file"`
	LocalYearDbDir      string        `toml:"local_year_db_dir"`
	LocalQuarterDbDir   string        `toml:"local_quarter_db_dir"`
	RemoteYearDbDir     string        `toml:"remote_year_db_dir"`
	RemoteQuarterDbDir  string        `toml:"remote_quarter_db_dir"`
	YearDB              bool          `toml:"year_db"`
	LogFileName         string        `toml:"log_file_name"`
	LogDirName          string        `toml:"log_dir_name"`
	LogFileSize         int64         `toml:"log_file_size"`
	W1Size              float32       `toml:"w1_size"`
	H1Size              float32       `toml:"h1_size"`
	W2Size              float32       `toml:"w2_size"`
	H2Size              float32       `toml:"h2_size"`
	Lang                string        `toml:"lang"`
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
func UpdateConfig(config *Config, configPath string) error {

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
