package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
	_ "github.com/nakagami/firebirdsql"
	indicator_tables_viewer "indicator-tables-viewer"
	"indicator-tables-viewer/internal/repository"
	"log"
)

func main() {

	a := app.New()
	w := a.NewWindow(" ")

	w.SetContent(widget.NewLabel(" "))
	w.Resize(fyne.NewSize(600, 400))
	w.Show()

	a.Run()

	cfg := indicator_tables_viewer.NewConfig()
	db, err := repository.NewFirebirdDB(cfg)
	if err != nil {
		log.Println(err)
	}
	repo := repository.NewRepository(db)
	tablesList, _ := repo.GetTables()

	fmt.Println(tablesList)
}
