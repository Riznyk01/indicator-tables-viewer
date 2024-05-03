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
	"indicator-tables-viewer/internal/logg"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
)

const (
	configPath          = "path to the configuration file"
	updateExists        = "program update exists"
	title               = "updates checker"
	errWhileDownloading = "error occurred while downloading update"
	updatedSuccessfully = "program update downloaded successfully.\nstarting the program."
	errOccur            = "error occurred while"
	loadingVer          = "loading ver info file"
	savingVer           = "saving ver info file"
	readingVer          = "reading local ver info file"
	newInfoFile         = "creating new version info file"
	readingNewVer       = "reading downloaded ver info file"
	remoteByte          = "remote []byte"
	localByte           = "local []byte"
	remoteConvErr       = "remoteVer converting error"
	localConvErr        = "localVer converting error"
	failedExePath       = "failed to get executable path"
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

	cfgPath = os.Getenv("CFG_PATH")

	cfgPathFlag := flag.String("CFG_PATH", "", "path to the config")
	flag.Parse()
	if *cfgPathFlag != "" {
		cfgPath = *cfgPathFlag
	}
	log.Printf("the path of the config is: %s", cfgPath)
	log.Printf("%s: %s", configPath, cfgPath)
	cfg := config.MustLoad(cfgPath)

	if cfg.Env == "dev" {
		exeDir = cfg.CodePath
	} else if cfg.Env == "prod" {
		exePath, err := os.Executable()
		if err != nil {
			log.Printf("%s: %v", failedExePath, err)
			return
		}
		exeDir = filepath.Dir(exePath)
		log.Printf("exeDir variable: %s", exeDir)
	}

	logger := logg.SetupLogger(cfg)
	logger.V(1).Info("launcher started")
	l := NewLauncher(cfg, logger)
	a := app.New()
	update := a.NewWindow(title)
	info := widget.NewLabel("start")

	if cfg.AutoUpdate {
		go func() {

			ex, err := l.checkUpdate()
			if err != nil {
				info.SetText(err.Error())
			}
			if ex && err == nil {

				info.SetText(updateExists)
				logger.V(1).Info(updateExists)

				err = downloader.Download(cfg.UpdatePath, cfg.RemoteExeFilename, cfg.CodePath+cfg.LocalExeFilename)
				if err != nil {
					logger.V(1).Error(err, errWhileDownloading+cfg.RemoteExeFilename)
					info.SetText(fmt.Sprintf("%s %s: %v", errWhileDownloading, cfg.RemoteExeFilename, err))
				} else {
					err = downloader.Download(cfg.UpdatePath, cfg.VerRemoteFilePath, cfg.CodePath+cfg.DownloadedVerFile)
					if err != nil {
						logger.V(1).Error(err, errWhileDownloading+cfg.VerRemoteFilePath)
						info.SetText(fmt.Sprintf("%s %s: %v", errWhileDownloading, cfg.VerRemoteFilePath, err))
					} else {
						err = os.Rename("ver_remote", "ver")
						if err != nil {
							log.Fatal(err)
						}
						logger.V(1).Info(updatedSuccessfully)
						info.SetText(updatedSuccessfully)
					}
				}
				l.run(exeDir)
				os.Exit(0)
			} else {
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
	resp, err := http.Get(l.cfg.UpdatePath + "/" + l.cfg.VerRemoteFilePath)
	if err != nil {
		log.Printf("%s %s: %v", errOccur, loadingVer, err)
		return false, err
	}
	defer resp.Body.Close()

	file, err := os.Create(l.cfg.CodePath + l.cfg.DownloadedVerFile)
	if err != nil {
		l.logger.V(1).Error(err, errOccur, newInfoFile)
		return false, err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		l.logger.V(1).Error(err, errOccur, savingVer)
		return false, err
	}

	localVer, err := os.ReadFile(l.cfg.CodePath + l.cfg.VerLocalFilePath)
	if err != nil {
		l.logger.V(1).Error(err, errOccur, readingVer)
		return false, err
	}

	remoteVer, err := os.ReadFile(l.cfg.CodePath + l.cfg.DownloadedVerFile)
	if err != nil {
		l.logger.V(1).Error(err, errOccur, readingNewVer)
		return false, err
	}

	l.logger.V(1).Info(remoteByte + string(remoteVer))
	l.logger.V(1).Info(localByte + string(remoteVer))
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
	if l.cfg.Env == "dev" {
		exeDir = l.cfg.CodePath
	}

	log.Printf("path for start viewer: %s\\%s\n", exeDir, l.cfg.LocalExeFilename)
	var cmd *exec.Cmd
	cmd = exec.Command(exeDir + "\\" + l.cfg.LocalExeFilename)
	cmd.Env = append(os.Environ(), "CONFIG_PATH="+cfgPath)
	if err := cmd.Start(); err != nil {
		l.logger.Error(err, errOccur+l.cfg.LocalExeFilename)
		return
	}
}
