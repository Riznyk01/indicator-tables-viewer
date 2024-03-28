package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"indicator-tables-viewer/internal/repository"
	"log"
)

type Viewer struct {
	window fyne.Window
	repo   *repository.Repository
}

func NewViewer(repo *repository.Repository) *Viewer {
	return &Viewer{
		repo: repo,
	}
}

func (v *Viewer) LoadUI(app fyne.App) {
	header := []string{"column1", "column2", "column3", "column2", "column3", "column2", "column3"}

	v.window = app.NewWindow("Indicator tables viewer")
	v.window.Resize(fyne.NewSize(1000, 600))

	tableDescription := widget.NewTable(
		func() (int, int) {
			return len(header), 1 // header
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("123")
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			label := o.(*widget.Label)
			label.Text = header[i.Row]
		},
	)
	tablesList, _ := v.repo.GetTable()
	dropdown := widget.NewSelect(tablesList, func(selected string) {
		tableName := selected[:7]
		//header, err = repo.Viewing.GetHeader(tableName)
		//indicatorNumbers, _ := repo.Viewing.GetIndicatorMaket(tableName)
		log.Println(tableName)
		//log.Printf("tmp message table header fetched: %s\n", header)
		//log.Printf("tmp message fetched the indicators' numbers: %s\n", indicatorNumbers)
		//table.Refresh()
	})

	content := container.NewVBox(
		widget.NewLabel("Select an indicators table for view:"),
		dropdown,
	)

	scr := container.NewVScroll(tableDescription)
	scr.SetMinSize(fyne.NewSize(1000, 400))

	v.window.SetContent(container.NewVBox(content, scr))
	v.window.ShowAndRun()
}
