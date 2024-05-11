package main

import (
	"flag"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
	"github.com/go-logr/logr"
	"indicator-tables-viewer/internal/config"
	"indicator-tables-viewer/internal/downloader"
	"indicator-tables-viewer/internal/filemanager"
	"indicator-tables-viewer/internal/logg"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
)

const (
	errOccur               = "error occurred while"
	errWhileExtracting     = "extracting update files"
	configPath             = "path to the configuration file"
	updateExists           = "program update exists"
	title                  = "updates checker"
	errWhileDownloading    = "downloading update"
	updatedSuccessfully    = "program update downloaded successfully.\nstarting the program."
	loadingVer             = "loading ver info file"
	savingVer              = "saving ver info file"
	readingVer             = "reading local ver info file"
	newInfoFile            = "creating new version info file"
	readingNewVer          = "reading downloaded ver info file"
	remoteByte             = "remote []byte"
	localByte              = "local []byte"
	remoteConvErr          = "remoteVer converting error"
	localConvErr           = "localVer converting error"
	failedExePath          = "failed to get executable path"
	updateDoesntExists     = "update doesn't exist"
	updateCheckingErr      = "update checking error"
	updDeletedSuccessfully = "archive deleted successfully"
	whileDeleting          = "while deleting the"
)

var cfgPath string

type Launcher struct {
	cfg    *config.Config
	logger *logr.Logger
}

func NewLauncher(cfg *config.Config, logger *logr.Logger) *Launcher {
	return &Launcher{
		cfg:    cfg,
		logger: logger,
	}
}

func main() {
	var exeDir string

	exePath, err := os.Executable()
	if err != nil {
		log.Printf("%s: %v", failedExePath, err)
		return
	}
	exeDir = filepath.Dir(exePath)
	log.Printf("exeDir variable: %s", exeDir)

	cfgPath = os.Getenv("CFG_PATH")

	cfgPathFlag := flag.String("CFG_PATH", "", "path to the config")
	flag.Parse()
	if *cfgPathFlag != "" {
		cfgPath = exeDir + "\\" + *cfgPathFlag
	}
	log.Printf("the path of the config is: %s", cfgPath)
	log.Printf("%s: %s", configPath, cfgPath)
	cfg := config.MustLoad(cfgPath)

	logger := logg.SetupLogger(cfg)
	logger.V(1).Info("cfg", "cfg", cfg)
	logger.V(1).Info("launcher started")
	l := NewLauncher(cfg, logger)
	a := app.New()
	update := a.NewWindow(title)
	info := widget.NewLabel("start")

	if cfg.AutoUpdate {
		go func() {
			l.logger.V(1).Info(fmt.Sprintf("update URL %s/%s:", l.cfg.UpdateURL, l.cfg.VerRemoteFilePath))

			err := downloader.Download(l.cfg.UpdateURL, l.cfg.VerRemoteFilePath, l.cfg.DownloadedVerFile)
			if err != nil {
				l.logger.V(1).Error(err, errOccur+errWhileDownloading+l.cfg.VerRemoteFilePath)
			}

			ex, err := l.checkUpdate()
			if err != nil {
				info.SetText(err.Error())
				logger.V(1).Error(err, fmt.Sprintf("%s", updateCheckingErr))
			}
			logger.V(1).Info("update", "exist", ex)
			if ex && err == nil {
				info.SetText(updateExists)
				err = downloader.Download(cfg.UpdateURL, cfg.UpdateArch, cfg.UpdateArch)
				if err != nil {
					info.SetText(fmt.Sprintf("%s %s: %v", errWhileDownloading, cfg.UpdateArch, err))
					logger.V(1).Error(err, errOccur+errWhileDownloading+cfg.UpdateArch)
				} else {

					err = filemanager.Unzip(cfg.UpdateArch, exeDir)
					if err != nil {
						info.SetText(fmt.Sprintf("%s %s: %v", errOccur, errWhileExtracting, err))
						logger.V(1).Error(err, errOccur+errWhileExtracting)

					} else {
						err = os.Rename("ver_remote", "ver")
						if err != nil {
							log.Fatal(err)
						}
						updArchPath := exeDir + "\\" + cfg.UpdateArch
						err = os.Remove(updArchPath)
						if err != nil {
							logger.V(1).Info(fmt.Sprintf("%s %s %s", errOccur, whileDeleting, updArchPath))
						}
						info.SetText(fmt.Sprintf("%s", updatedSuccessfully))
						logger.V(1).Info(updatedSuccessfully)
						logger.V(1).Info(fmt.Sprintf("%s %s", updArchPath, updDeletedSuccessfully))
					}
				}
				l.run(exeDir)
				os.Exit(0)
			} else {
				logger.V(1).Info(updateDoesntExists)
				l.run(exeDir)
				if err != nil {
					fmt.Printf(err.Error())
				}
				os.Exit(0)
			}
		}()
	} else {
		l.run(exeDir)
		os.Exit(0)
	}
	update.SetContent(info)
	update.CenterOnScreen()
	update.Resize(fyne.NewSize(300, 100))
	update.Show()
	a.Run()
}
func (l *Launcher) checkUpdate() (bool, error) {

	localVer, err := os.ReadFile(l.cfg.VerLocalFilePath)
	if err != nil {
		l.logger.V(1).Error(err, errOccur, readingVer)
		return false, err
	}

	remoteVer, err := os.ReadFile(l.cfg.DownloadedVerFile)
	if err != nil {
		l.logger.V(1).Error(err, errOccur, readingNewVer)
		return false, err
	}

	l.logger.V(1).Info(localByte + string(localVer))
	l.logger.V(1).Info(remoteByte + string(remoteVer))
	// removing non-digit characters
	regex := regexp.MustCompile("[^0-9]+")
	cleanedRemoteStr, cleanedLocalStr := regex.ReplaceAllString(string(remoteVer), ""), regex.ReplaceAllString(string(localVer), "")

	remote, err := strconv.Atoi(cleanedRemoteStr)
	if err != nil {
		l.logger.V(1).Error(err, remoteConvErr)
	}
	log.Printf("remote int %v", remote)
	local, err := strconv.Atoi(cleanedLocalStr)
	if err != nil {
		l.logger.Error(err, localConvErr)
	}
	log.Printf("local int %v", local)
	return remote > local, nil
}

func (l *Launcher) run(exeDir string) {
	log.Printf("path for start viewer: %s\\%s\n", exeDir, l.cfg.LocalExeFilename)
	var cmd *exec.Cmd
	cmd = exec.Command(exeDir + "\\" + l.cfg.LocalExeFilename)
	cmd.Env = append(os.Environ(), "CONFIG_PATH="+cfgPath)
	if err := cmd.Start(); err != nil {
		l.logger.Error(err, errOccur+l.cfg.LocalExeFilename)
		return
	}
}
