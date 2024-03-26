package repository

import (
	"database/sql"
	"indicator-tables-viewer/internal/models"
)

type Viewing interface {
	GetTables() ([]models.Table, []string, error)
}

type Repository struct {
	Viewing
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Viewing: NewViewingFirebird(db),
	}
}
