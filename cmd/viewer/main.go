package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/go-logr/logr"
	_ "github.com/nakagami/firebirdsql"
	"indicator-tables-viewer/internal/config"
	"indicator-tables-viewer/internal/filemanager"
	"indicator-tables-viewer/internal/formatter"
	"indicator-tables-viewer/internal/logg"
	"indicator-tables-viewer/internal/models"
	"indicator-tables-viewer/internal/repository"
	"indicator-tables-viewer/internal/translator"
	"indicator-tables-viewer/internal/ui"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	readingVer = "reading local ver info file"
)

func main() {
	cfgPath := os.Getenv("CONFIG_PATH")
	cfg := config.MustLoad(cfgPath)

	exePath, err := os.Executable()
	if err != nil {
		log.Print(err)
	}
	cfg.LocalPath = filepath.Dir(exePath)
	logger := logg.SetupLogger(cfg)

	lang, err := translator.LoadTranslations(cfg.Lang)
	if err != nil {
		logger.V(1).Error(err, "error occurred while loading localization")
	}
	logger.V(1).Info("language file", cfg.Lang, lang)

	logger.V(1).Info("viewer started")
	logger.V(1).Info("the dir to the exe file", "path", cfg.LocalPath)
	logger.V(1).Info("", "the path of the config is", cfgPath)

	err = filemanager.MakeDirIfNotExist(cfg.XlsExportPath)
	logger.Error(err, "")

	a := app.New()
	logger.V(1).Info("resources", "path", cfg.LocalPath+cfg.IconPath)
	r, _ := loadRecourseFromPath(cfg.LocalPath + cfg.IconPath)
	a.SetIcon(r)

	sizer := newTermTheme(cfg.FontSize)
	a.Settings().SetTheme(sizer)

	logger.V(1).Info("path to the ver file", "path", cfg.VerLocalFilePath)
	localVer, err := os.ReadFile(cfg.VerLocalFilePath)
	if err != nil {
		logger.Error(err, lang["ErrOccur"], readingVer)
	}

	verInfo := formatter.VersionFormatter(localVer)
	logger.V(1).Info("", lang["Version"], verInfo)

	login := a.NewWindow(fmt.Sprintf("%s %s %s", lang["LoginFormTitle"], lang["Version"], verInfo))

	usernameEntry := widget.NewEntry()
	usernameEntry.SetText(cfg.Username)
	usernameEntry.SetPlaceHolder(lang["Username"])

	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder(lang["Password"])

	settingsButton := widget.NewButton(lang["SettingsButtonText"], func() {
		newSettingsWindow(a, cfg, cfgPath, usernameEntry, lang)
	})

	username := container.NewGridWithColumns(4, widget.NewLabel(""), usernameEntry, passwordEntry, widget.NewLabel(""))

	dbPath := widget.NewEntry()
	dbPath.SetText(cfg.RemotePathToDb + "/" + cfg.DBName)

	checkboxLocalMode := widget.NewCheck(lang["LocalDB"], func(checked bool) {
		if checked {
			cfg.LocalMode = true
			dbPath.SetText("in program folder")
			logger.V(1).Info("db", "path", dbPath.Text)
			logger.V(1).Info("config to update", "cfg", cfg)
			err = config.UpdateConfig(cfg, cfgPath)
			if err != nil {
				errDialog := dialog.NewInformation("error", err.Error(), login)
				errDialog.Show()
				logger.Error(err, "")
			} else {
				logger.V(1).Info("LocalModeCheckbox state in the config file has been updated", "checkbox state", cfg.LocalMode)
			}
		} else {
			cfg.LocalMode = false
			dbPath.SetText(cfg.RemotePathToDb)
			logger.V(1).Info("db", "path", dbPath.Text)
			logger.V(1).Info("config to update", "cfg", cfg)
			err = config.UpdateConfig(cfg, cfgPath)
			if err != nil {
				errDialog := dialog.NewInformation("error", err.Error(), login)
				errDialog.Show()
			} else {
				logger.V(1).Info("LocalModeCheckbox state in the config file has been updated", "checkbox state", cfg.LocalMode)
			}
		}
	})

	checkboxYearDB := widget.NewCheck(lang["YearDB"], func(checked bool) {
		if checked {
			cfg.YearDB = true
			logger.V(1).Info("config to update", "cfg", cfg)
			err = config.UpdateConfig(cfg, cfgPath)
			if err != nil {
				errDialog := dialog.NewInformation("error", err.Error(), login)
				errDialog.Show()
				logger.Error(err, "")
			} else {
				logger.V(1).Info("checkboxYearDB state in the config file has been updated", "checkbox state", cfg.YearDB)
			}
		} else {
			cfg.YearDB = false
			logger.V(1).Info("config to update", "cfg", cfg)
			err = config.UpdateConfig(cfg, cfgPath)
			if err != nil {
				errDialog := dialog.NewInformation("error", err.Error(), login)
				errDialog.Show()
			} else {
				logger.V(1).Info("checkboxYearDB state in the config file has been updated", "checkbox state", cfg.YearDB)
			}
		}
	})

	checkboxLocalMode.SetChecked(cfg.LocalMode)
	logger.V(1).Info("checkboxLocalMode is set according to the configuration")
	checkboxYearDB.SetChecked(cfg.YearDB)
	logger.V(1).Info("checkboxYearDB is set according to the configuration")
	loginButton := widget.NewButton(lang["Login"], func() {
		db, connectionString, err := repository.NewFirebirdDB(cfg, usernameEntry.Text, passwordEntry.Text)
		if err != nil {
			logger.Error(err, "")
			errDialog := dialog.NewInformation("error", err.Error(), login)
			errDialog.Show()
		} else {
			login.Hide()
			logger.V(1).Info("connected")
			connStr := widget.NewLabel(fmt.Sprintf("%s  %s", lang["AppName"], connectionString))

			repo := repository.NewRepository(db)
			newViewerWindow(a, logger, repo, cfg, connStr, lang)
		}
	})

	checkRow := container.NewGridWithColumns(3, widget.NewLabel(""), loginButton, widget.NewLabel(""))
	settingsRow := container.NewGridWithColumns(3, checkboxLocalMode, checkboxYearDB, settingsButton)

	loginW := container.NewGridWithRows(3, username, checkRow, settingsRow)

	login.SetContent(loginW)
	login.Resize(fyne.NewSize(800, 180))
	login.CenterOnScreen()
	login.Show()
	a.Run()
}

func newSettingsWindow(app fyne.App, cfg *config.Config, cfgPath string, usernameEntry *widget.Entry, lang models.Translations) {
	settings := app.NewWindow(lang["SettingsTitle"])
	settings.Resize(fyne.NewSize(500, 100))
	settings.CenterOnScreen()

	dbName := widget.NewEntry()
	dbName.SetText(cfg.DBName)
	dbNameCols := container.NewGridWithColumns(2, widget.NewLabel(lang["DBFileName"]), dbName)

	username := widget.NewEntry()
	username.SetText(cfg.Username)
	usernameSettingsCols := container.NewGridWithColumns(2, widget.NewLabel(lang["Username"]+":"), username)

	remoteHost := widget.NewEntry()
	remoteHost.SetText(cfg.Host)
	remotePort := widget.NewEntry()
	remotePort.SetText(cfg.Port)
	remoteHostCols := container.NewGridWithColumns(4, widget.NewLabel(lang["RemoteDBSettings"]), widget.NewLabel("[host:port]: "), remoteHost, remotePort)

	remotePath := widget.NewEntry()
	remotePath.SetText(cfg.RemotePathToDb)
	remotePathCols := container.NewGridWithColumns(3, widget.NewLabel(""), widget.NewLabel(lang["PathToRemoteDB"]), remotePath)

	remoteYearDBDir := widget.NewEntry()
	remoteYearDBDir.SetText(cfg.RemoteYearDbDir)
	remoteYearDBDirCols := container.NewGridWithColumns(3, widget.NewLabel(""), widget.NewLabel(lang["DirForYearDB"]), remoteYearDBDir)

	remoteQuarterDBDir := widget.NewEntry()
	remoteQuarterDBDir.SetText(cfg.RemoteQuarterDbDir)
	remoteQuarterDBDirCols := container.NewGridWithColumns(3, widget.NewLabel(""), widget.NewLabel(lang["DirForQuarterDB"]), remoteQuarterDBDir)

	localYearDBDir := widget.NewEntry()
	localYearDBDir.SetText(cfg.LocalYearDbDir)
	localYearDBDirCols := container.NewGridWithColumns(3, widget.NewLabel(""), widget.NewLabel(lang["DirForYearDB"]), localYearDBDir)

	localQuarterDBDir := widget.NewEntry()
	localQuarterDBDir.SetText(cfg.LocalQuarterDbDir)
	localQuarterDBDirCols := container.NewGridWithColumns(3, widget.NewLabel(""), widget.NewLabel(lang["DirForQuarterDB"]), localQuarterDBDir)

	localHost := widget.NewEntry()
	localHost.SetText(cfg.LocalHost)
	localPort := widget.NewEntry()
	localPort.SetText(cfg.LocalPort)
	localHostCols := container.NewGridWithColumns(4, widget.NewLabel(lang["LocalDBSettings"]), widget.NewLabel("[host:port]: "), localHost, localPort)

	infoTimeout := widget.NewEntry()
	infoTimeout.SetText(fmt.Sprintf("%v", cfg.InfoTimeout))
	infoTimeoutCols := container.NewGridWithColumns(3, widget.NewLabel(lang["OtherSettings"]), widget.NewLabel(lang["InfoTimeout"]), infoTimeout)

	xlsExport := widget.NewEntry()

	selectDirButton := widget.NewButton(lang["ChooseFolderButtonText"], func() {
		dirDialog := dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {
			if err == nil && uri != nil {
				fmt.Printf("%ы %s", lang["SelectedDirectory"], uri.Path())
				cfg.XlsExportPath = uri.Path()
				xlsExport.SetText(fmt.Sprintf("%v", cfg.XlsExportPath))
			}
		}, settings)
		dirDialog.Show()
	})

	xlsExport.SetText(fmt.Sprintf("%v", cfg.XlsExportPath))
	xlsExportCols := container.NewGridWithColumns(4, widget.NewLabel(""), widget.NewLabel(lang["XLSExportPath"]), xlsExport, selectDirButton)

	resolutionWidthMultiplier := widget.NewEntry()
	resolutionWidthMultiplier.SetText(fmt.Sprintf("%v", cfg.WidthMultiplier))

	resolutionHeightMultiplier := widget.NewEntry()
	resolutionHeightMultiplier.SetText(fmt.Sprintf("%v", cfg.HeightMultiplier))

	resolutionMultiplierCols := container.NewGridWithColumns(4, widget.NewLabel(""), widget.NewLabel(lang["MultiplierSettings"]), resolutionWidthMultiplier, resolutionHeightMultiplier)

	setFontSize := widget.NewEntry()
	setFontSize.SetText(fmt.Sprintf("%v", cfg.FontSize))
	setFontSizeCols := container.NewGridWithColumns(3, widget.NewLabel(lang["SetFontSize"]), widget.NewLabel(""), setFontSize)

	saveSettingsButton := widget.NewButton(lang["SaveSettingsButtonText"], func() {

		newInfoTimeout, err := time.ParseDuration(infoTimeout.Text)
		if err != nil {
			log.Printf("%s parsing info messages timeout: %v\n", lang["ErrOccur"], err)
			return
		}

		cfg.DBName = dbName.Text
		cfg.Username = username.Text
		cfg.Host = remoteHost.Text
		cfg.Port = remotePort.Text
		cfg.RemotePathToDb = remotePath.Text
		cfg.RemoteYearDbDir = remoteYearDBDir.Text
		cfg.RemoteQuarterDbDir = remoteQuarterDBDir.Text
		cfg.LocalYearDbDir = localYearDBDir.Text
		cfg.LocalQuarterDbDir = localQuarterDBDir.Text
		cfg.LocalHost = localHost.Text
		cfg.LocalPort = localPort.Text
		cfg.InfoTimeout = newInfoTimeout
		cfg.XlsExportPath = xlsExport.Text

		widthMultiplFloat, _ := stringToFloat(resolutionWidthMultiplier.Text)
		heightMultiplFloat, _ := stringToFloat(resolutionHeightMultiplier.Text)
		fontSizeFloat, _ := stringToFloat(setFontSize.Text)

		cfg.WidthMultiplier = widthMultiplFloat
		cfg.HeightMultiplier = heightMultiplFloat
		cfg.FontSize = fontSizeFloat

		err = config.UpdateConfig(cfg, cfgPath)

		if err != nil {
			log.Println(err)
			errDialog := dialog.NewInformation(lang["SettingsTitle"], err.Error(), settings)
			errDialog.Show()
		} else {
			successDialog := dialog.NewInformation(lang["SettingsTitle"], "config has been changed", settings)
			successDialog.Show()
			usernameEntry.SetText(cfg.Username)
		}
		settings.Close()
	})

	buttonsRowCols := container.NewGridWithColumns(3, widget.NewLabel(""), widget.NewLabel(""), saveSettingsButton)

	settingsRows := container.NewGridWithRows(14,
		dbNameCols,
		usernameSettingsCols,
		remoteHostCols,
		remotePathCols,
		remoteQuarterDBDirCols,
		remoteYearDBDirCols,
		localHostCols,
		localQuarterDBDirCols,
		localYearDBDirCols,
		infoTimeoutCols,
		xlsExportCols,
		resolutionMultiplierCols,
		setFontSizeCols,
		buttonsRowCols)

	settings.SetContent(settingsRows)
	settings.Show()
}
func newViewerWindow(app fyne.App, logger *logr.Logger, repo *repository.Repository, cfg *config.Config, connStr *widget.Label, lang models.Translations) {
	log.Printf("Main window is started")
	var tableName string
	statData := newData()

	w, h := ui.SetResolution(cfg)

	window := app.NewWindow(lang["AppName"])
	window.Resize(fyne.NewSize(w, h))
	window.SetMaster()
	tablesList, _ := repo.GetTable()

	t := newTable(statData)

	dropdown := widget.NewSelect(tablesList, func(selected string) {
		tableName = selected[:7]

		var formName string

		if selected[3] == 75 {
			// if P20K form is chose
			formName = "F" + selected[1:3] + "K"
			tableName = selected[:8]
		} else {
			formName = "F" + selected[1:3]
		}

		logger.V(1).Info("selected", "table name", tableName)
		logger.V(1).Info("selected", "form name", formName)
		// get where columns names is located
		colNameLocation, _ := repo.GetColNameLocation(tableName)
		// get the tables' header
		statData[0], _ = repo.GetHeader(colNameLocation)

		formatter.LineSplit(statData)

		headerHeight := rowHeightCount(statData[0])
		t.SetRowHeight(0, headerHeight)

		indicators, _ := repo.GetIndicatorNumbers(tableName)
		for m := 1; m < len(statData); m++ {
			statData[m] = make([]string, 14)
		}
		if len(indicators) != 0 {
			for ind, _ := range indicators {
				statData[ind+1] = []string{
					repo.GetIndicator(formName, indicators[ind].P1, indicators[ind].STRTAB, indicators[ind].SHSTR, indicators[ind].NZAP),
					repo.GetIndicator(formName, indicators[ind].P2, indicators[ind].STRTAB, indicators[ind].SHSTR, indicators[ind].NZAP),
					repo.GetIndicator(formName, indicators[ind].P3, indicators[ind].STRTAB, indicators[ind].SHSTR, indicators[ind].NZAP),
					repo.GetIndicator(formName, indicators[ind].P4, indicators[ind].STRTAB, indicators[ind].SHSTR, indicators[ind].NZAP),
					repo.GetIndicator(formName, indicators[ind].P5, indicators[ind].STRTAB, indicators[ind].SHSTR, indicators[ind].NZAP),
					repo.GetIndicator(formName, indicators[ind].P6, indicators[ind].STRTAB, indicators[ind].SHSTR, indicators[ind].NZAP),
					repo.GetIndicator(formName, indicators[ind].P7, indicators[ind].STRTAB, indicators[ind].SHSTR, indicators[ind].NZAP),
					repo.GetIndicator(formName, indicators[ind].P8, indicators[ind].STRTAB, indicators[ind].SHSTR, indicators[ind].NZAP),
					repo.GetIndicator(formName, indicators[ind].P9, indicators[ind].STRTAB, indicators[ind].SHSTR, indicators[ind].NZAP),
					repo.GetIndicator(formName, indicators[ind].P10, indicators[ind].STRTAB, indicators[ind].SHSTR, indicators[ind].NZAP),
					repo.GetIndicator(formName, indicators[ind].P11, indicators[ind].STRTAB, indicators[ind].SHSTR, indicators[ind].NZAP),
					repo.GetIndicator(formName, indicators[ind].P12, indicators[ind].STRTAB, indicators[ind].SHSTR, indicators[ind].NZAP),
					repo.GetIndicator(formName, indicators[ind].P13, indicators[ind].STRTAB, indicators[ind].SHSTR, indicators[ind].NZAP),
					repo.GetIndicator(formName, indicators[ind].P14, indicators[ind].STRTAB, indicators[ind].SHSTR, indicators[ind].NZAP),
				}
			}
		}

		// set the data rows height
		for rows := 1; rows < len(statData); rows++ {
			// set the row height
			rowHeight := rowHeightCount(statData[rows]) * (cfg.FontSize / 13)
			t.SetRowHeight(rows, rowHeight)
		}

		setColumnWidth(t, statData, cfg.FontSize)
		t.Refresh()
	})

	exportFileButton := widget.NewButton(lang["ExportButtonText"], func() {
		err := filemanager.ExportToExcel(statData, tableName, cfg.XlsExportPath)
		if err != nil {
			window.SetTitle(err.Error())
			<-time.After(cfg.InfoTimeout)
			window.SetTitle(connStr.Text)
		} else {
			window.SetTitle(lang["FileSaved"])
			<-time.After(cfg.InfoTimeout)
			window.SetTitle(connStr.Text)
		}
	})

	horizontalContent := container.NewHBox(
		widget.NewLabel(lang["SelectTable"]),
		dropdown,
		exportFileButton,
	)

	scr := container.NewVScroll(t)
	scr.SetMinSize(fyne.NewSize(w, h))
	window.SetTitle(connStr.Text)
	window.SetContent(container.NewVBox(horizontalContent, scr))
	window.Show()
}

// rowHeightCount set the headers' height
func rowHeightCount(rowToCount []string) float32 {
	count := 0
	for _, field := range rowToCount {
		q := strings.Count(field, "\n")
		if q > count {
			count = q
		}
	}
	if count == 0 {
		return 24
	}
	return float32(count) * 24
}

// setColumnWidth set the columns width after fetching new data
// depending on data len and splitting by \n
func setColumnWidth(t *widget.Table, statData [][]string, fontSize float32) {

	for c := 0; c < len(statData[0]); c++ {
		var maxLen int
		for r := 0; r < len(statData); r++ {
			if len(statData[r][c]) > maxLen {
				maxLen = formatter.MaxLengthAfterSplit(statData[r][c])
			}
		}
		log.Printf("max text len for the %v column is: %v", c, maxLen)
		t.SetColumnWidth(c, float32(maxLen)*7*(fontSize/13))
	}
}

func newData() [][]string {
	rows := 150
	data := make([][]string, rows)
	for i := 0; i < rows; i++ {
		row := make([]string, 14)
		for cell := 0; cell < len(row)-1; cell++ {
			row[cell] = "-"
		}
		data[i] = row
	}
	return data
}

func newTable(statData [][]string) *widget.Table {
	tbl := widget.NewTable(
		func() (int, int) {
			return len(statData), len(statData[0])
		},
		func() fyne.CanvasObject {
			label := widget.NewLabel("-")
			return label
		},
		func(tci widget.TableCellID, c fyne.CanvasObject) {
			c.(*widget.Label).SetText(statData[tci.Row][tci.Col])
		},
	)
	return tbl
}

func stringToFloat(s string) (float32, error) {
	f, err := strconv.ParseFloat(s, 32)
	if err != nil {
		return 0.8, err
	}
	return float32(f), nil
}
