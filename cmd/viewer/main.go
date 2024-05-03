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
	"indicator-tables-viewer/internal/repository"
	"indicator-tables-viewer/internal/text"
	"indicator-tables-viewer/internal/ui"
	"log"
	"os"
	"path/filepath"
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

	logger.V(1).Info("viewer started")

	logger.V(1).Info("the dir to the exe file", "path", cfg.LocalPath)
	logger.V(1).Info("", "the path of the config is", cfgPath)

	a := app.New()
	logger.V(1).Info("resources", "path", cfg.LocalPath+cfg.IconPath)
	r, _ := loadRecourseFromPath(cfg.LocalPath + cfg.IconPath)
	a.SetIcon(r)

	sizer := newTermTheme()
	a.Settings().SetTheme(sizer)

	login := a.NewWindow("Login Form")

	usernameEntry := widget.NewEntry()
	usernameEntry.SetText(cfg.Username)
	usernameEntry.SetPlaceHolder("username")

	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder("password")

	settingsButton := widget.NewButton("settings", func() {
		newSettingsWindow(a, cfg, cfgPath, usernameEntry)
	})

	logger.V(1).Info("path to the ver file", "path", cfg.CodePath+cfg.VerLocalFilePath)
	localVer, err := os.ReadFile(cfg.CodePath + cfg.VerLocalFilePath)
	if err != nil {
		logger.Error(err, text.ErrOccur, readingVer)
	}

	versionLabel := widget.NewLabel("")
	verInfo := formatter.VersionFormatter(localVer, cfg.Env)
	logger.V(1).Info("", "version", verInfo)
	versionLabel.SetText("version: " + verInfo)

	username := container.NewGridWithColumns(4, widget.NewLabel(""), usernameEntry, passwordEntry, widget.NewLabel(""))

	dbPath := widget.NewEntry()
	dbPath.SetText(cfg.RemotePathToDb + "/" + cfg.DBName)

	checkboxLocalMode := widget.NewCheck("local db", func(checked bool) {
		if checked {
			cfg.LocalMode = true
			dbPath.SetText("in program folder")
			logger.V(1).Info("db", "path", dbPath.Text)
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
			err = config.UpdateConfig(cfg, cfgPath)
			if err != nil {
				errDialog := dialog.NewInformation("error", err.Error(), login)
				errDialog.Show()
			} else {
				logger.V(1).Info("LocalModeCheckbox state in the config file has been updated", "checkbox state", cfg.LocalMode)
			}
		}
	})

	checkboxYearDB := widget.NewCheck("year db (qtr. if not)", func(checked bool) {
		if checked {
			cfg.YearDB = true
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
			err = config.UpdateConfig(cfg, cfgPath)
			if err != nil {
				errDialog := dialog.NewInformation("error", err.Error(), login)
				errDialog.Show()
			} else {
				logger.V(1).Info("checkboxYearDB state in the config file has been updated", "checkbox state", cfg.YearDB)
			}
		}
	})

	connStr := widget.NewLabel("")

	checkboxLocalMode.SetChecked(cfg.LocalMode)
	logger.V(1).Info("checkboxLocalMode is set according to the configuration")
	checkboxYearDB.SetChecked(cfg.YearDB)
	logger.V(1).Info("checkboxYearDB is set according to the configuration")
	loginButton := widget.NewButton("login", func() {
		db, connectionString, err := repository.NewFirebirdDB(cfg, usernameEntry.Text, passwordEntry.Text)
		if err != nil {
			logger.Error(err, "")
			errDialog := dialog.NewInformation("error", err.Error(), login)
			errDialog.Show()
		} else {
			login.Hide()
			logger.V(1).Info("connected")
			connStr.SetText("connection: " + connectionString)
			repo := repository.NewRepository(db)
			newViewerWindow(a, logger, repo, cfg, connStr)
		}
	})

	textRow := container.NewGridWithColumns(4, widget.NewLabel(""), widget.NewLabel("username"), widget.NewLabel("password"), widget.NewLabel(""))
	checkRow := container.NewGridWithColumns(3, widget.NewLabel(""), loginButton, widget.NewLabel(""))
	settingsRow := container.NewGridWithColumns(5, versionLabel, widget.NewLabel(""), checkboxLocalMode, checkboxYearDB, settingsButton)

	loginW := container.NewGridWithRows(4, textRow, username, checkRow, settingsRow)
	login.SetContent(loginW)
	login.Resize(fyne.NewSize(1000, 100))
	login.CenterOnScreen()
	login.Show()
	a.Run()
}

func newSettingsWindow(app fyne.App, cfg *config.Config, cfgPath string, usernameEntry *widget.Entry) {
	settings := app.NewWindow("settings")
	settings.Resize(fyne.NewSize(500, 100))
	settings.CenterOnScreen()

	dbName := widget.NewEntry()
	dbName.SetText(cfg.DBName)
	dbNameCols := container.NewGridWithColumns(2, widget.NewLabel("db file name: "), dbName)

	usernameSettings := widget.NewEntry()
	usernameSettings.SetText(cfg.Username)
	usernameSettingsCols := container.NewGridWithColumns(2, widget.NewLabel("username: "), usernameSettings)

	remotePort := widget.NewEntry()
	remotePort.SetText(cfg.Port)
	remoteHost := widget.NewEntry()
	remoteHost.SetText(cfg.Host)
	remoteHostCols := container.NewGridWithColumns(4, widget.NewLabel("remote db settings"), widget.NewLabel("[host:port]: "), remoteHost, remotePort)

	remotePath := widget.NewEntry()
	remotePath.SetText(cfg.RemotePathToDb)
	remotePathCols := container.NewGridWithColumns(3, widget.NewLabel(""), widget.NewLabel("path to remote DB: "), remotePath)

	remoteYearDBDir := widget.NewEntry()
	remoteYearDBDir.SetText(cfg.RemoteYearDbDir)
	remoteYearDBDirCols := container.NewGridWithColumns(3, widget.NewLabel(""), widget.NewLabel("dir name for year DB: "), remoteYearDBDir)

	remoteQuarterDBDir := widget.NewEntry()
	remoteQuarterDBDir.SetText(cfg.RemoteQuarterDbDir)
	remoteQuarterDBDirCols := container.NewGridWithColumns(3, widget.NewLabel(""), widget.NewLabel("dir name for quarter DB: "), remoteQuarterDBDir)

	localYearDBDir := widget.NewEntry()
	localYearDBDir.SetText(cfg.LocalYearDbDir)
	localYearDBDirCols := container.NewGridWithColumns(3, widget.NewLabel(""), widget.NewLabel("dir name for year DB: "), localYearDBDir)

	localQuarterDBDir := widget.NewEntry()
	localQuarterDBDir.SetText(cfg.LocalQuarterDbDir)
	localQuarterDBDirCols := container.NewGridWithColumns(3, widget.NewLabel(""), widget.NewLabel("dir name for quarter DB: "), localQuarterDBDir)

	localPort := widget.NewEntry()
	localPort.SetText(cfg.LocalPort)
	localHost := widget.NewEntry()
	localHost.SetText(cfg.LocalHost)
	localHostCols := container.NewGridWithColumns(4, widget.NewLabel("local db settings"), widget.NewLabel("[host:port]: "), localHost, localPort)

	infoTimeout := widget.NewEntry()
	infoTimeout.SetText(fmt.Sprintf("%v", cfg.InfoTimeout))
	infoTimeoutCols := container.NewGridWithColumns(3, widget.NewLabel("other settings "), widget.NewLabel("info messages timeout: "), infoTimeout)

	xlsExport := widget.NewEntry()

	selectDirButton := widget.NewButton("Choose Folder", func() {
		dirDialog := dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {
			if err == nil && uri != nil {
				fmt.Printf("Selected directory: %s", uri.Path())
				cfg.XlsExportPath = uri.Path()
				xlsExport.SetText(fmt.Sprintf("%v", cfg.XlsExportPath))
			}
		}, settings)
		dirDialog.Show()
	})

	xlsExport.SetText(fmt.Sprintf("%v", cfg.XlsExportPath))
	xlsExportCols := container.NewGridWithColumns(4, widget.NewLabel(""), widget.NewLabel("XLS export path (program dir if empty):"), xlsExport, selectDirButton)

	saveSettingsButton := widget.NewButton("save settings", func() {

		newInfoTimeout, err := time.ParseDuration(infoTimeout.Text)
		if err != nil {
			log.Printf("error occurred while parsing info messages timeout: %v\n", err)
			return
		}

		cfg.Username = usernameEntry.Text
		cfg.Host = remoteHost.Text
		cfg.Port = remotePort.Text
		cfg.RemotePathToDb = remotePath.Text
		cfg.LocalHost = localHost.Text
		cfg.LocalPort = localPort.Text
		cfg.DBName = dbName.Text
		cfg.InfoTimeout = newInfoTimeout
		cfg.XlsExportPath = xlsExport.Text

		err = config.UpdateConfig(cfg, cfgPath)

		if err != nil {
			log.Println(err)
			errDialog := dialog.NewInformation("settings", err.Error(), settings)
			errDialog.Show()
		} else {
			successDialog := dialog.NewInformation("settings", "config has been changed", settings)
			successDialog.Show()
			usernameEntry.SetText(cfg.Username)
		}
		settings.Close()
	})

	buttonsRowCols := container.NewGridWithColumns(3, widget.NewLabel(""), widget.NewLabel(""), saveSettingsButton)

	settingsRows := container.NewGridWithRows(12,
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
		buttonsRowCols)

	settings.SetContent(settingsRows)
	settings.Show()
}
func newViewerWindow(app fyne.App, logger *logr.Logger, repo *repository.Repository, cfg *config.Config, connStr *widget.Label) {
	log.Printf("Main window is started")
	var tableName string
	statData := newData()

	w, h := ui.SetResolution()
	windowSize := fyne.NewSize(w, h)
	tableSize := fyne.NewSize(w, h)

	window := app.NewWindow("Indicator tables viewer")
	window.FullScreen()
	window.Resize(windowSize)
	window.SetMaster()
	tablesList, _ := repo.GetTable()

	t := newTable(statData)

	dropdown := widget.NewSelect(tablesList, func(selected string) {
		tableName = selected[:7]
		formName := "F" + selected[1:3]
		logger.V(1).Info("selected", "table name", tableName)
		logger.V(1).Info("selected", "table name", formName)
		// get where columns names is located
		colNameLocation, _ := repo.GetColNameLocation(tableName)
		// get the tables' header
		statData[0], _ = repo.GetHeader(colNameLocation)

		formatter.LineSplit(statData)
		// set the headers' height
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
			rowHeight := rowHeightCount(statData[rows])
			t.SetRowHeight(rows, rowHeight)
		}

		setColumnWidth(t, statData)
		t.Refresh()
	})

	info := widget.NewLabel("")

	exportFileButton := widget.NewButton("export to excel", func() {
		err := filemanager.ExportToExcel(statData, tableName, cfg.XlsExportPath)
		if err != nil {
			info.SetText(err.Error())
			<-time.After(cfg.InfoTimeout)
			info.SetText("")
		} else {
			info.SetText("file saved successfully")
			<-time.After(cfg.InfoTimeout)
			info.SetText("")
		}
	})

	updateDateButton := widget.NewButton("update DB correction date", func() {
		err := repo.UpdateDBCorrectionDate(time.Now())
		if err != nil {
			info.SetText(err.Error())
			<-time.After(cfg.InfoTimeout)
			info.SetText("")
		} else {
			info.SetText("date updated successfully")
			<-time.After(cfg.InfoTimeout)
			info.SetText("")
		}
	})

	horizontalContent := container.NewHBox(
		widget.NewLabel("select an indicators table for view:"),
		dropdown,
		exportFileButton,
		updateDateButton,
	)

	scr := container.NewVScroll(t)
	scr.SetMinSize(tableSize)

	infoLine := container.NewGridWithColumns(2, info, connStr)

	window.SetContent(container.NewVBox(horizontalContent, scr, infoLine))
	window.Show()
}

func rowHeightCount(rowToCount []string) float32 {
	count := 0
	for _, field := range rowToCount {
		q := strings.Count(field, "\n")
		if q > count {
			count = q
		}
		//log.Printf("the filed to count is: %v, the field len is: %v\n", count, len(field))
	}
	//log.Printf("quanity of the \\n is: %v\n", count)
	if count == 0 {
		return 24
	}
	return float32(count) * 24
}

// setColumnWidth set the columns width after fetching new data
// depending on data len and splitting by \n
func setColumnWidth(t *widget.Table, statData [][]string) {

	for c := 0; c < len(statData[0]); c++ {
		var maxLen int
		for r := 0; r < len(statData); r++ {
			if len(statData[r][c]) > maxLen {
				maxLen = formatter.MaxLengthAfterSplit(statData[r][c])
			}
		}
		log.Printf("max text len for the %v column is: %v", c, maxLen)
		t.SetColumnWidth(c, float32(maxLen)*7)
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
