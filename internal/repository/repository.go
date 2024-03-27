package repository

import (
	"database/sql"
)

type Viewing interface {
	GetTable() ([]string, error)
	GetHeader(tableName string) ([]string, error)
	GetIndicatorMaket(tableName string) ([]string, error)
}

type Repository struct {
	Viewing
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Viewing: NewViewingFirebird(db),
	}
}
