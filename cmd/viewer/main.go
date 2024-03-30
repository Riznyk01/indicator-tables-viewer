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

	w := newTerminalWindow(a, repo)
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

func newTerminalWindow(app fyne.App, repo *repository.Repository) fyne.Window {
	statData := newData()
	windowSize := fyne.NewSize(1500, 800)
	tableSize := fyne.NewSize(1500, 800)
	var tableName string

	window := app.NewWindow("Indicator tables viewer")
	window.Resize(windowSize)

	tablesList, _ := repo.GetTable()

	t := newTable(statData)
	t.SetRowHeight(0, float32(140))
	dropdown := widget.NewSelect(tablesList, func(selected string) {
		tableName = selected[:7]
		log.Printf("%s selected", tableName)
		statData[0], _ = repo.GetHeader(tableName)
		lineSplit(statData)
		for l := 0; l < len(statData[0]); l++ {
			t.SetColumnWidth(l, float32(lenBefore(statData[0][l]))*6)
		}
		indicators, _ := repo.GetIndicatorMaket(tableName)
		for m := 1; m < len(statData); m++ {
			statData[m] = []string{"", "", "", "", "", "", "", "", "", "", "", "", "", ""}
		}

		for ind, _ := range indicators {
			statData[ind+1] = []string{
				indicators[ind].P1,
				indicators[ind].P2,
				indicators[ind].P3,
				indicators[ind].P4,
				indicators[ind].P5,
				indicators[ind].P6,
				indicators[ind].P7,
				indicators[ind].P8,
				indicators[ind].P9,
				indicators[ind].P10,
				indicators[ind].P11,
				indicators[ind].P12,
				indicators[ind].P13,
				indicators[ind].P14,
			}
		}
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

func lenBefore(str string) (length int) {
	index := strings.Index(str, "\n")
	if index != -1 {
		length = index
	} else {
		length = len(str)
	}
	return length
}

func newData() [][]string {
	rows := 100
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
