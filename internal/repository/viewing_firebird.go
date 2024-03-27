package repository

import (
	"database/sql"
	"fmt"
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

type IndicatorData struct {
	TABL   string
	STRTAB string
	NZAP   string
	SHSTR  string
	FORMA  string
	STROKA string
	P1     string
	P2     string
	P3     string
	P4     string
	P5     string
	P6     string
	P7     string
	P8     string
	P9     string
	P10    string
	P11    string
	P12    string
	P13    string
	P14    string
}

func (v *ViewingFirebird) GetTable() ([]string, error) {
	var tableNames []string
	rows, err := v.db.Query("SELECT tabl, nazv FROM STABLES WHERE tabl LIKE  'P%'")
	if err != nil {
		log.Println("Error querying database:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var column1 string
		var column2 string
		if err := rows.Scan(&column1, &column2); err != nil {
			log.Println("Error scanning rows:", err)
			return nil, err
		}
		tableNames = append(tableNames, fmt.Sprintf("%s %s", column1, column2))
	}
	if err := rows.Err(); err != nil {
		log.Println("Error iterating over rows:", err)
		return nil, err
	}
	return tableNames, nil
}
func (v *ViewingFirebird) GetHeader(tableName string) ([]string, error) {
	var tableCols []string
	query := "SELECT NAME FROM SCOL WHERE TABL = ? AND COL != 0"
	rows, err := v.db.Query(query, tableName)

	if err != nil {
		log.Println("Error querying database:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var column string
		if err := rows.Scan(&column); err != nil {
			log.Println("Error scanning rows:", err)
			return nil, err
		}
		tableCols = append(tableCols, column)
	}
	if err := rows.Err(); err != nil {
		log.Println("Error iterating over rows:", err)
		return nil, err
	}
	return tableCols, nil
}

func (v *ViewingFirebird) GetIndicatorMaket(tableName string) ([]string, error) {
	fc := "GetIndicatorMaket"
	var data IndicatorData
	var queryCols string
	var indicatorNums []string

	log.Printf("table cols for query %s", queryCols)
	query := "SELECT * FROM PMAKET WHERE TABL = ?"
	rows, err := v.db.Query(query, tableName)

	if err != nil {
		log.Printf("%s Error querying database: %v", fc, err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var column string
		if err := rows.Scan(&data.TABL,
			&data.STRTAB,
			&data.NZAP,
			&data.SHSTR,
			&data.FORMA,
			&data.STROKA,
			&data.P1,
			&data.P2,
			&data.P3,
			&data.P4,
			&data.P5,
			&data.P6,
			&data.P7,
			&data.P8,
			&data.P9,
			&data.P10,
			&data.P11,
			&data.P12,
			&data.P13,
			&data.P14,
		); err != nil {
			log.Printf("%s Error scanning rows: %v", fc, err)
			return nil, err
		}
		indicatorNums = append(indicatorNums, column)
	}
	if err := rows.Err(); err != nil {
		log.Printf("%s Error iterating over rows: %v", fc, err)
		return nil, err
	}
	return indicatorNums, nil
}
