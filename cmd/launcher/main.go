package main

import (
	"flag"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
	"indicator-tables-viewer/internal/config"
	"indicator-tables-viewer/internal/downloader"
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
	cfg *config.Config
}

func NewLauncher(cfg *config.Config) *Launcher {
	return &Launcher{
		cfg: cfg,
	}
}

func main() {
	var exeDir string
	logFile, err := os.OpenFile("logfile.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("error occurred while opening logfile:", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	if os.Getenv("CFG_PATH") == "" {
		cfgPath = "config/config_prod.toml"
	} else {
		cfgPath = os.Getenv("CFG_PATH")
	}

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

	l := NewLauncher(cfg)
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
				log.Printf(updateExists)

				err = downloader.Download(cfg.UpdatePath, cfg.RemoteExeFilename, cfg.CodePath+cfg.LocalExeFilename)
				if err != nil {
					log.Printf("%s %s: %v", errWhileDownloading, cfg.RemoteExeFilename, err)
					info.SetText(fmt.Sprintf("%s %s: %v", errWhileDownloading, cfg.RemoteExeFilename, err))
				} else {
					err = downloader.Download(cfg.UpdatePath, cfg.VerRemoteFilePath, cfg.CodePath+cfg.DownloadedVerFile)
					if err != nil {
						log.Printf("%s %s: %v", errWhileDownloading, cfg.VerRemoteFilePath, err)
						info.SetText(fmt.Sprintf("%s %s: %v", errWhileDownloading, cfg.VerRemoteFilePath, err))
					} else {
						err = os.Rename("ver_remote", "ver")
						if err != nil {
							log.Fatal(err)
						}

						log.Printf(updatedSuccessfully)
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
		log.Printf("%s %s: %v", errOccur, newInfoFile, err)
		return false, err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		log.Printf("%s %s: %v", errOccur, savingVer, err)
		return false, err
	}

	localVer, err := os.ReadFile(l.cfg.CodePath + l.cfg.VerLocalFilePath)
	if err != nil {
		log.Printf("%s %s: %v", errOccur, readingVer, err)
		return false, err
	}

	remoteVer, err := os.ReadFile(l.cfg.CodePath + l.cfg.DownloadedVerFile)
	if err != nil {
		log.Printf("%s %s: %v", errOccur, readingNewVer, err)
		return false, err
	}

	log.Printf("%s %v\n", remoteByte, string(remoteVer))
	log.Printf("%s %v\n", localByte, string(localVer))
	// removing non-digit characters
	regex := regexp.MustCompile("[^0-9]+")
	cleanedRemoteStr, cleanedLocalStr := regex.ReplaceAllString(string(remoteVer), ""), regex.ReplaceAllString(string(localVer), "")

	remote, err := strconv.Atoi(cleanedRemoteStr)
	if err != nil {
		log.Printf("%s: %v", remoteConvErr, err)
	}
	log.Printf("remote int %v", remote)
	local, err := strconv.Atoi(cleanedLocalStr)
	if err != nil {
		log.Printf("%s: %v", localConvErr, err)
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
		log.Printf("%s starting %s: %v", errOccur, l.cfg.LocalExeFilename, err)
		return
	}
}
