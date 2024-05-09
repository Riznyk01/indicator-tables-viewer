package logg

import (
	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
	"indicator-tables-viewer/internal/config"
	"indicator-tables-viewer/internal/filemanager"
	"log"
)

func SetupLogger(cfg *config.Config) *logr.Logger {
	var logFilePath string

	if cfg.Env == "prod" {
		logFilePath = cfg.LocalPath + "\\"
	}

	logFilePath += cfg.LogDirName + cfg.LogFileName + "_" + cfg.LocalExeFilename[:len(cfg.LocalExeFilename)-4] + cfg.LogFileExt

	err := filemanager.CheckLogFileSize(logFilePath, cfg.LogFileSize)
	if err != nil {
		log.Print(err)
	}

	cfgLog := zap.NewDevelopmentConfig()
	cfgLog.OutputPaths = []string{logFilePath}

	log, err := cfgLog.Build()
	if err != nil {
		panic(err)
	}
	logger := zapr.NewLogger(log)

	logger.V(1).Info("", "log file path", logFilePath)

	return &logger
}
