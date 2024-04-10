package repository

import (
	"database/sql"
	"indicator-tables-viewer/internal/models"
	"time"
)

type Viewing interface {
	GetTable() ([]string, error)
	GetHeader(tableName string) ([]string, error)
	GetIndicatorNumbers(tableName string) ([]models.IndicatorData, error)
	GetIndicator(shifr, npokaz, indicatorsRow, decodingRow, decodingTable string) string
	GetColNameLocation(tableName string) (string, error)
	UpdateDBCorrectionDate(currentTime time.Time) error
}

type Repository struct {
	Viewing
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Viewing: NewViewingFirebird(db),
	}
}
