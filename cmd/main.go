package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	_ "github.com/nakagami/firebirdsql"
	"indicator-tables-viewer/internal/config"
	"indicator-tables-viewer/internal/repository"
	"log"
)

func main() {

	a := app.New()
	window := a.NewWindow("Dropdown Example")
	window.Resize(fyne.NewSize(1000, 600))

	cfg := config.NewConfig()
	db, err := repository.NewFirebirdDB(cfg)
	if err != nil {
		log.Println(err)
	}
	repo := repository.NewRepository(db)
	tablesList, _ := repo.GetTable()

	dropdown := widget.NewSelect(tablesList, func(selected string) {
		tableName := selected[:7]
		header, _ := repo.Viewing.GetHeader(tableName)
		indicatorNumbers, _ := repo.Viewing.GetIndicatorMaket(tableName)
		log.Printf("tmp message table header fetched: %s\n", header)
		log.Printf("tmp message fetched the indicators' numbers: %s\n", indicatorNumbers)
	})
	content := container.NewVBox(
		widget.NewLabel("Select a table:"),
		dropdown,
	)
	window.SetContent(content)
	window.ShowAndRun()

}
