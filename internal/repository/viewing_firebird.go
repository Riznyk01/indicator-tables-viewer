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

func (v *ViewingFirebird) GetTables() ([]models.Table, []string, error) {
	var tableArr []models.Table
	var tableNames []string

	rows, err := v.db.Query("SELECT tabl, nazv FROM STABLES WHERE tabl LIKE  'P%'")
	if err != nil {
		log.Println("Error querying database:", err)
		return nil, nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var columns models.Table
		if err := rows.Scan(&columns.Ident, &columns.Name); err != nil {
			log.Println("Error scanning rows:", err)
			return nil, nil, err
		}
		tableArr = append(tableArr, columns)
		tableNames = append(tableNames, columns.Name)
	}
	if err := rows.Err(); err != nil {
		log.Println("Error iterating over rows:", err)
		return nil, nil, err
	}
	return tableArr, tableNames, nil
}
