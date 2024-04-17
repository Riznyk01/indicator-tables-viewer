package main

import (
	"errors"
	"flag"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/go-vgo/robotgo"
	_ "github.com/nakagami/firebirdsql"
	"github.com/tealeg/xlsx"
	"indicator-tables-viewer/internal/config"
	"indicator-tables-viewer/internal/repository"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Resource interface {
	Name() string
	Content() []byte
}

type StaticResource struct {
	StaticName    string
	StaticContent []byte
}

func NewStaticResource(name string, content []byte) *StaticResource {
	return &StaticResource{
		StaticName:    name,
		StaticContent: content,
	}
}

func (r *StaticResource) Name() string {
	return r.StaticName
}

func (r *StaticResource) Content() []byte {
	return r.StaticContent
}

func main() {
	cfgPath := "config/config_prod.toml"

	pathPtr := flag.String("path", "", "the path to the file")
	flag.Parse()
	fmt.Printf("path flag: %v\n", *pathPtr)

	if *pathPtr != "" {
		cfgPath = *pathPtr
	}
	log.Printf("the path of the config is: %s", cfgPath)
	logFile, err := os.OpenFile("logfile.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("error occurred while opening logfile:", err)
	}
	defer logFile.Close()
	//log.SetOutput(logFile)
	log.SetOutput(os.Stdout)

	cfg := config.MustLoad(cfgPath)

	exePath, err := os.Executable()
	if err != nil {
		log.Println("error occurred while receiving exe file:", err)
		return
	}

	exeDir := filepath.Dir(exePath)
	log.Println("the dir to the exe file:", exeDir)

	cfg.LocalPath = exeDir

	a := app.New()

	r, _ := loadRecourseFromPath("cmd/viewer/data/Icon.png")
	a.SetIcon(r)

	sizer := newTermTheme()
	a.Settings().SetTheme(sizer)

	login := a.NewWindow("Login Form")

	usernameEntry := widget.NewEntry()
	usernameEntry.SetText(cfg.Username)
	usernameEntry.SetPlaceHolder("username")

	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetText(cfg.Password)
	passwordEntry.SetPlaceHolder("password")

	settingsButton := widget.NewButton("settings", func() {
		newSettingsWindow(a, cfg, cfgPath, usernameEntry, passwordEntry)
	})

	username := container.NewGridWithColumns(4, widget.NewLabel(""), usernameEntry, passwordEntry, widget.NewLabel(""))

	dbPath := widget.NewEntry()
	dbPath.SetText(cfg.Path + "/" + cfg.DBName)

	checkbox := widget.NewCheck("local db", func(checked bool) {
		if checked {
			cfg.LocalMode = true
			passwordEntry.SetText(cfg.LocalPassword)
			dbPath.SetText("in program folder")
			err = config.SaveLocalModeCheckboxState(*cfg, cfgPath)
			if err != nil {
				errDialog := dialog.NewInformation("error", err.Error(), login)
				errDialog.Show()
				log.Println(err)
			} else {
				log.Printf("LocalModeCheckbox state in the config has been updated (%v)", cfg.LocalMode)
			}
		} else {
			cfg.LocalMode = false
			dbPath.SetText(cfg.Path)
			passwordEntry.SetText(cfg.Password)
			err = config.SaveLocalModeCheckboxState(*cfg, cfgPath)
			if err != nil {
				errDialog := dialog.NewInformation("error", err.Error(), login)
				errDialog.Show()
			} else {
				log.Printf("LocalModeCheckbox state in the config has been updated (%v)", cfg.LocalMode)
			}
		}
	})
	checkbox.SetChecked(cfg.LocalMode)
	loginButton := widget.NewButton("login", func() {
		db, err := repository.NewFirebirdDB(cfg, usernameEntry.Text, passwordEntry.Text, cfg.LocalMode)
		if err != nil {
			log.Println(err)
			errDialog := dialog.NewInformation("error", err.Error(), login)
			errDialog.Show()
		} else {
			login.Hide()
			log.Printf("connected")
			repo := repository.NewRepository(db)
			newViewerWindow(a, repo, cfg)
		}
	})

	checkRow := container.NewGridWithColumns(3, widget.NewLabel(""), loginButton, widget.NewLabel(""))
	settingsRow := container.NewGridWithColumns(5, widget.NewLabel(""), widget.NewLabel(""), widget.NewLabel(""), checkbox, settingsButton)
	textRow := container.NewGridWithColumns(4, widget.NewLabel(""), widget.NewLabel("username"), widget.NewLabel("password"), widget.NewLabel(""))
	loginW := container.NewGridWithRows(6, textRow, username, checkRow, widget.NewLabel(""), widget.NewLabel(""), settingsRow)

	login.SetContent(loginW)
	login.Resize(fyne.NewSize(1000, 100))
	login.CenterOnScreen()
	login.Show()
	a.Run()
}

func loadRecourseFromPath(path string) (Resource, error) {

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open the file %v", err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		log.Println("error occurred while get information about the file:", err)
		return nil, err
	}
	size := fileInfo.Size()

	iconData := make([]byte, size)
	_, err = file.Read(iconData)
	if err != nil {
		log.Println("error occurred while reading the file:", err)
		return nil, err
	}
	name := filepath.Base(path)
	return NewStaticResource(name, iconData), nil
}
func newSettingsWindow(app fyne.App, cfg *config.Config, cfgPath string, usernameEntry *widget.Entry, passwordEntry *widget.Entry) {
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
	remotePath.SetText(cfg.Path)
	remotePathCols := container.NewGridWithColumns(3, widget.NewLabel(""), widget.NewLabel("path to DB: "), remotePath)

	remotePass := widget.NewEntry()
	remotePass.SetText(cfg.Password)
	remotePassCols := container.NewGridWithColumns(3, widget.NewLabel(""), widget.NewLabel("pass: "), remotePass)

	localPort := widget.NewEntry()
	localPort.SetText(cfg.LocalPort)
	localHost := widget.NewEntry()
	localHost.SetText(cfg.LocalHost)
	localHostCols := container.NewGridWithColumns(4, widget.NewLabel("local db settings"), widget.NewLabel("[host:port]: "), localHost, localPort)

	localPath := widget.NewEntry()
	localPath.SetText(cfg.LocalPath)
	localPathCols := container.NewGridWithColumns(3, widget.NewLabel(""), widget.NewLabel("path to DB (the program dir): "), localPath)

	localPass := widget.NewEntry()
	localPass.SetText(cfg.LocalPassword)
	localPassCols := container.NewGridWithColumns(3, widget.NewLabel(""), widget.NewLabel("pass: "), localPass)

	infoTimeout := widget.NewEntry()
	infoTimeout.SetText(fmt.Sprintf("%v", cfg.InfoTimeout))
	infoTimeoutCols := container.NewGridWithColumns(3, widget.NewLabel("other settings "), widget.NewLabel("info messages timeout: "), infoTimeout)

	saveSettingsButton := widget.NewButton("save settings", func() {

		newInfoTimeout, err := time.ParseDuration(infoTimeout.Text)
		if err != nil {
			log.Printf("error occurred while parsing info messages timeout: %v", err)
			return
		}

		cfg.Username = usernameEntry.Text
		cfg.Host = remoteHost.Text
		cfg.Port = remotePort.Text
		cfg.Path = remotePath.Text
		cfg.Password = remotePass.Text
		cfg.LocalHost = localHost.Text
		cfg.LocalPort = localPort.Text
		cfg.LocalPath = localPath.Text
		cfg.LocalPassword = localPass.Text
		cfg.DBName = dbName.Text
		cfg.InfoTimeout = newInfoTimeout

		err = config.UpdateConfig(*cfg, cfgPath)

		if err != nil {
			log.Println(err)
			errDialog := dialog.NewInformation("settings", err.Error(), settings)
			errDialog.Show()
		} else {
			successDialog := dialog.NewInformation("settings", "config has been changed", settings)
			successDialog.Show()
			usernameEntry.SetText(cfg.Username)
			if cfg.LocalMode {
				passwordEntry.SetText(cfg.LocalPassword)
			} else {
				passwordEntry.SetText(cfg.Password)
			}
		}
		settings.Close()
	})

	buttonsRowCols := container.NewGridWithColumns(3, widget.NewLabel(""), widget.NewLabel(""), saveSettingsButton)

	settingsRows := container.NewGridWithRows(10,
		dbNameCols,
		usernameSettingsCols,
		remoteHostCols,
		remotePathCols,
		remotePassCols,
		localHostCols,
		localPathCols,
		localPassCols,
		infoTimeoutCols,
		buttonsRowCols)

	settings.SetContent(settingsRows)
	settings.Show()
}
func newViewerWindow(app fyne.App, repo *repository.Repository, cfg *config.Config) {
	log.Printf("Main window is started")
	var tableName string
	statData := newData()

	w, h := setResolution()
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
		log.Printf("%s selected", tableName)
		log.Printf("%s selected", formName)
		// get where columns names is located
		colNameLocation, _ := repo.GetColNameLocation(tableName)
		// get the tables' header
		statData[0], _ = repo.GetHeader(colNameLocation)

		lineSplit(statData)
		// set the headers' height
		headerHeight := rowHeightCount(statData[0])
		t.SetRowHeight(0, headerHeight)

		indicators, _ := repo.GetIndicatorNumbers(tableName)
		for m := 1; m < len(statData); m++ {
			statData[m] = []string{"empty", "empty", "empty", "empty", "empty", "empty", "empty", "empty", "empty", "empty", "empty", "empty", "empty", "empty"}
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
		err := exportToExcel(statData, tableName)
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

	window.SetContent(container.NewVBox(horizontalContent, scr, info))
	window.Show()
}

func rowHeightCount(rowToCount []string) float32 {
	count := 0
	for _, field := range rowToCount {
		q := strings.Count(field, "\n")
		if q > count {
			count = q
		}
		log.Printf("the filed to count is: %v, the field len is: %v\n", count, len(field))
	}
	log.Printf("quanity of the \\n is: %v\n", count)
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
				maxLen = maxLengthAfterSplit(statData[r][c])
			}
		}
		log.Printf("max text len for the %v column is: %v", c, maxLen)
		t.SetColumnWidth(c, float32(maxLen)*7)
	}
}

// lineSplit adds the "\n" to the headers' rows
func lineSplit(data [][]string) {
	every := 7
	for n, colName := range data[0] {
		var result string
		cnt := 0
		for i, char := range colName {
			result += string(char)
			cnt++
			if ((i+1) > every || (i+1) < 2*every) && char == ' ' && cnt > 8 {
				result += "\n"
				cnt = 0
			}
		}
		data[0][n] = strings.ReplaceAll(result, "|", "\n")
	}
}

func maxLengthAfterSplit(str string) int {
	substrings := strings.Split(str, "\n")
	maxLength := 0

	if strings.Contains(str, "\n") {
		for _, sub := range substrings {
			if len(sub) > maxLength {
				maxLength = len(sub)
			}
		}
	} else {
		maxLength = len(str)
	}

	return maxLength
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

func setResolution() (w, h float32) {
	width, height := robotgo.GetScreenSize()
	if width > 1920 && height > 1080 {
		w = 0.5 * float32(width)
		h = 0.4 * float32(height)
		return w, h
	}
	w = 0.8 * float32(width)
	h = 0.8 * float32(height)
	return w, h
}

func exportToExcel(data [][]string, tableName string) error {
	file := xlsx.NewFile()
	sheet, err := file.AddSheet("Sheet1")
	if err != nil {
		return errors.New(fmt.Sprintf("Error creating sheet: %v\n", err))
	}

	for _, row := range data {
		newRow := sheet.AddRow()
		for _, cell := range row {
			newCell := newRow.AddCell()
			if cell == "empty" {
				cell = ""
			}
			newCell.Value = cell
		}
	}

	log.Println("File saved successfully.")
	for s := 0; s < sheet.MaxCol; s++ {
		sheet.Col(s).Width = float64(30)
	}

	currentTime := time.Now()
	currentDateTime := currentTime.Format("2006-01-02_15-04-05")
	filename := tableName + "_" + currentDateTime + ".xlsx"
	err = file.Save(filename)
	if err != nil {
		return errors.New(fmt.Sprintf("Error saving file: %v\n", err))
	}
	return nil
}
