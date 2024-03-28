package main

import (
	"fyne.io/fyne/v2/app"
	_ "github.com/nakagami/firebirdsql"
	"indicator-tables-viewer/internal/config"
	"indicator-tables-viewer/internal/repository"
	"indicator-tables-viewer/internal/ui"
	"log"
)

func main() {
	a := app.New()
	//a.SetIcon(resourceIconPng)

	cfg := config.NewConfig()
	db, err := repository.NewFirebirdDB(cfg)
	if err != nil {
		log.Println(err)
	}
	repo := repository.NewRepository(db)

	v := ui.NewViewer(repo)
	v.LoadUI(a)
	//v.Run()
}
