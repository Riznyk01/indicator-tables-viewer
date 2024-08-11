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
	"indicator-tables-viewer/internal/models"
	"indicator-tables-viewer/internal/translator"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
)

const (
	updateDoesntExists     = "update doesn't exist"
	updateCheckingErr      = "update checking error"
	updDeletedSuccessfully = "archive deleted successfully"
	failedExePath          = "failed to get executable path"
	remoteConvErr          = "remoteVer converting error"
	localConvErr           = "localVer converting error"
	localByte              = "local []byte"
	configPath             = "path to the configuration file"
	readingVer             = "reading local ver info file"
	readingNewVer          = "reading downloaded ver info file"
	remoteByte             = "remote []byte"
	errWhileDownloading    = "downloading update"
)

var cfgPath string

type Launcher struct {
	cfg    *config.Config
	logger *logr.Logger
	lang   models.Translations
}

func NewLauncher(cfg *config.Config, logger *logr.Logger, lang models.Translations) *Launcher {
	return &Launcher{
		cfg:    cfg,
		logger: logger,
		lang:   lang,
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

	lang, err := translator.LoadTranslations(cfg.Lang)
	if err != nil {
		logger.V(1).Error(err, "error occurred while loading localization")
	}
	logger.V(1).Info("language file", cfg.Lang, lang)

	l := NewLauncher(cfg, logger, lang)
	a := app.New()
	update := a.NewWindow(l.lang["Title"])
	info := widget.NewLabel(l.lang["Start"])

	if cfg.AutoUpdate {
		go func() {
			l.logger.V(1).Info(fmt.Sprintf("update URL %s/%s:", l.cfg.UpdateURL, l.cfg.VerRemoteFilePath))

			err := downloader.Download(l.cfg.UpdateURL, l.cfg.VerRemoteFilePath, l.cfg.DownloadedVerFile, l.lang)
			if err != nil {
				l.logger.V(1).Error(err, l.lang["ErrOccur"]+errWhileDownloading+l.cfg.VerRemoteFilePath)
			}

			ex, err := l.checkUpdate()
			if err != nil {
				info.SetText(err.Error())
				logger.V(1).Error(err, fmt.Sprintf("%s", updateCheckingErr))
			}
			logger.V(1).Info("update", "exist", ex)
			if ex && err == nil {
				info.SetText(l.lang["UpdateExists"])
				err = downloader.Download(cfg.UpdateURL, cfg.UpdateArch, cfg.UpdateArch, l.lang)
				if err != nil {
					info.SetText(fmt.Sprintf("%s %s: %v", errWhileDownloading, cfg.UpdateArch, err))
					logger.V(1).Error(err, l.lang["ErrOccur"]+errWhileDownloading+cfg.UpdateArch)
				} else {

					err = filemanager.Unzip(cfg.UpdateArch, exeDir)
					if err != nil {
						info.SetText(fmt.Sprintf("%s %s: %v", l.lang["ErrOccur"], l.lang["ErrWhileExtracting"], err))
						logger.V(1).Error(err, l.lang["ErrOccur"]+l.lang["ErrWhileExtracting"])

					} else {
						err = os.Rename("ver_remote", "ver")
						if err != nil {
							log.Fatal(err)
						}
						updArchPath := exeDir + "\\" + cfg.UpdateArch
						err = os.Remove(updArchPath)
						if err != nil {
							logger.V(1).Info(fmt.Sprintf("%s %s %s", l.lang["ErrOccur"], l.lang["WhileDeleting"], updArchPath))
						}
						info.SetText(fmt.Sprintf("%s", l.lang["UpdatedSuccessfully"]))
						logger.V(1).Info(l.lang["UpdatedSuccessfully"])
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
		l.logger.V(1).Error(err, l.lang["ErrOccur"], readingVer)
		return false, err
	}

	remoteVer, err := os.ReadFile(l.cfg.DownloadedVerFile)
	if err != nil {
		l.logger.V(1).Error(err, l.lang["ErrOccur"], readingNewVer)
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
		l.logger.Error(err, l.lang["ErrOccur"]+l.cfg.LocalExeFilename)
		return
	}
}
