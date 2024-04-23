package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
	"indicator-tables-viewer/internal/config"
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
	cfgPath           = "build/config/config_prod.toml"
	remoteVerInfoPath = "ver_remote"
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
	downUpd             = "downloading updating file"
	failedExePath       = "failed to get executable path"
	fileCreating        = "creating file"
	doesntExist         = "doesn't exist"
)

func main() {
	log.Printf("the path of the config is: %s", cfgPath)
	logFile, err := os.OpenFile("logfile.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("error occurred while opening logfile:", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	log.Printf("%s: %s", configPath, cfgPath)

	cfg := config.MustLoad(cfgPath)

	a := app.New()
	update := a.NewWindow(title)
	info := widget.NewLabel("start")

	if cfg.AutoUpdate {
		go func() {
			ex, err := updateAvailable(cfg)
			if err != nil {
				info.SetText(err.Error())
			}
			if ex && err == nil {

				info.SetText(updateExists)
				log.Printf(updateExists)

				err = fileDownloader(cfg.UpdatePath, cfg.RemoteExeFilename, cfg.LocalExeFilename)
				if err != nil {
					log.Printf("%s %s: %v", errWhileDownloading, cfg.RemoteExeFilename, err)
					info.SetText(fmt.Sprintf("%s %s: %v", errWhileDownloading, cfg.RemoteExeFilename, err))
				} else {
					err = fileDownloader(cfg.UpdatePath, cfg.VerFilePath, remoteVerInfoPath)
					if err != nil {
						log.Printf("%s %s: %v", errWhileDownloading, cfg.VerFilePath, err)
						info.SetText(fmt.Sprintf("%s %s: %v", errWhileDownloading, cfg.VerFilePath, err))
					} else {
						err = os.Rename("ver_remote", "ver")
						if err != nil {
							log.Fatal(err)
						}

						log.Printf(updatedSuccessfully)
						info.SetText(updatedSuccessfully)
					}
				}
				runViewer(cfg)
				os.Exit(0)
			} else {
				runViewer(cfg)
				fmt.Printf(err.Error())
				os.Exit(0)
			}
		}()
	} else {
		runViewer(cfg)
		os.Exit(0)
	}
	update.SetContent(info)
	update.CenterOnScreen()
	update.Resize(fyne.NewSize(300, 100))
	update.Show()
	a.Run()
}
func updateAvailable(cfg *config.Config) (bool, error) {
	resp, err := http.Get(cfg.UpdatePath + "/" + cfg.VerFilePath)
	if err != nil {
		log.Printf("%s %s: %v", errOccur, loadingVer, err)
		return false, err
	}
	defer resp.Body.Close()

	file, err := os.Create(remoteVerInfoPath)
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

	localVer, err := os.ReadFile(cfg.VerFilePath)
	if err != nil {
		log.Printf("%s %s: %v", errOccur, readingVer, err)
		return false, err
	}

	remoteVer, err := os.ReadFile(remoteVerInfoPath)
	if err != nil {
		log.Printf("%s %s: %v", errOccur, readingNewVer, err)
		return false, err
	}

	log.Printf("%s %v\n%s %v", remoteByte, string(remoteVer), localByte, string(localVer))
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

func fileDownloader(updatePath, fileName, savePath string) error {
	resp, err := http.Get(updatePath + "/" + fileName)
	if err != nil {
		log.Printf("%s %s: %v", errOccur, downUpd, err)
		return err
	}
	defer resp.Body.Close()

	file, err := os.Create(savePath)
	if err != nil {
		log.Printf("%s %s: %v", errOccur, fileCreating, err)
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		log.Printf("%s copying file to the dir: %v", errOccur, err)
		return err
	}
	return nil
}

func runViewer(cfg *config.Config) {
	exePath, err := os.Executable()
	if err != nil {
		log.Printf("%s: %v", failedExePath, err)
		return
	}
	exeDir := filepath.Dir(exePath)
	viewerExe := filepath.Join(exeDir, cfg.LocalExeFilename)
	if _, err = os.Stat(viewerExe); os.IsNotExist(err) {
		log.Printf("%s %s", cfg.LocalExeFilename, doesntExist)
		return
	}
	cmd := exec.Command(viewerExe)
	if err = cmd.Start(); err != nil {
		log.Printf("%s starting %s: %v", errOccur, cfg.LocalExeFilename, err)
		return
	}
}
