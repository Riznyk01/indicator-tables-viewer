package repository

import (
	"database/sql"
	"indicator-tables-viewer/internal/models"
	"log"
)

type ViewingFirebird struct {
	db *sql.DB
}

func NewViewingFirebird(db *sql.DB) *ViewingFirebird {
	return &ViewingFirebird{
		db: db,
	}
}

func (v *ViewingFirebird) GetTables() ([]models.Table, error) {
	var tableArr []models.Table

	rows, err := v.db.Query("SELECT tabl, nazv FROM STABLES WHERE tabl LIKE  'P%'")
	if err != nil {
		log.Println("Error querying database:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var columns models.Table
		if err := rows.Scan(&columns.Ident, &columns.Name); err != nil {
			log.Println("Error scanning rows:", err)
			return nil, err
		}
		tableArr = append(tableArr, columns)
	}
	if err := rows.Err(); err != nil {
		log.Println("Error iterating over rows:", err)
		return nil, err
	}
	return tableArr, nil
}
