package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	_ "github.com/nakagami/firebirdsql"
	"indicator-tables-viewer/internal/config"
	"indicator-tables-viewer/internal/repository"
	"log"
	"os"
	"path/filepath"
	"strings"
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
	logFile, err := os.OpenFile("logfile.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("error occurred while opening logfile:", err)
	}
	defer logFile.Close()
	//log.SetOutput(logFile)
	log.SetOutput(os.Stdout)

	a := app.New()

	r, _ := LoadRecourseFromPath("cmd/viewer/data/Icon.png")
	a.SetIcon(r)

	sizer := newTermTheme()
	a.Settings().SetTheme(sizer)

	cfg := config.MustLoad()
	db, err := repository.NewFirebirdDB(cfg)
	if err != nil {
		log.Println(err)
	}
	repo := repository.NewRepository(db)

	w := newViewerWindow(a, repo, cfg)
	w.ShowAndRun()
}

func LoadRecourseFromPath(path string) (Resource, error) {

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

func newViewerWindow(app fyne.App, repo *repository.Repository, cfg *config.Config) fyne.Window {
	statData := newData()
	windowSize := fyne.NewSize(cfg.WindowWidth, cfg.WindowHeight)
	tableSize := fyne.NewSize(cfg.WindowWidth, cfg.WindowHeight)

	window := app.NewWindow("Indicator tables viewer")
	window.Resize(windowSize)

	tablesList, _ := repo.GetTable()

	t := newTable(statData)
	// set the header height
	t.SetRowHeight(0, cfg.HeaderHeight)
	// set the data rows height
	for row := 1; row < len(statData); row++ {
		t.SetRowHeight(row, cfg.RowHeight)
	}

	dropdown := widget.NewSelect(tablesList, func(selected string) {
		tableName := selected[:7]
		formName := "F" + selected[1:3]
		log.Printf("%s selected", tableName)
		log.Printf("%s selected", formName)
		statData[0], _ = repo.GetHeader(tableName)
		lineSplit(statData)

		indicators, _ := repo.GetIndicatorNumbers(tableName)
		for m := 1; m < len(statData); m++ {
			statData[m] = []string{"", "", "", "", "", "", "", "", "", "", "", "", "", ""}
		}

		for ind, _ := range indicators {
			statData[ind+1] = []string{
				repo.GetIndicator(formName, indicators[ind].P1),
				repo.GetIndicator(formName, indicators[ind].P2),
				repo.GetIndicator(formName, indicators[ind].P3),
				repo.GetIndicator(formName, indicators[ind].P4),
				repo.GetIndicator(formName, indicators[ind].P5),
				repo.GetIndicator(formName, indicators[ind].P6),
				repo.GetIndicator(formName, indicators[ind].P7),
				repo.GetIndicator(formName, indicators[ind].P8),
				repo.GetIndicator(formName, indicators[ind].P9),
				repo.GetIndicator(formName, indicators[ind].P10),
				repo.GetIndicator(formName, indicators[ind].P11),
				repo.GetIndicator(formName, indicators[ind].P12),
				repo.GetIndicator(formName, indicators[ind].P13),
				repo.GetIndicator(formName, indicators[ind].P14),
			}
		}

		setColumnWidth(t, statData)

		t.Refresh()
	})

	content := container.NewVBox(
		widget.NewLabel("select an indicators table for view:"),
		dropdown,
	)

	scr := container.NewVScroll(t)
	scr.SetMinSize(tableSize)

	window.SetContent(container.NewVBox(content, scr))
	return window
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
	rows := 200
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
	tbl := widget.NewTableWithHeaders(
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
