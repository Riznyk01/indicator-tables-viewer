package repository

import (
	"database/sql"
	"fmt"
	"indicator-tables-viewer/internal/models"
	"log"
	"strconv"
)

type ViewingFirebird struct {
	db *sql.DB
}

func NewViewingFirebird(db *sql.DB) *ViewingFirebird {
	return &ViewingFirebird{
		db: db,
	}
}

type IndicatorDataFirebird struct {
	TABL   string
	STRTAB sql.NullString
	NZAP   string
	SHSTR  sql.NullString
	FORMA  sql.NullString
	STROKA sql.NullString
	P1     sql.NullString
	P2     sql.NullString
	P3     sql.NullString
	P4     sql.NullString
	P5     sql.NullString
	P6     sql.NullString
	P7     sql.NullString
	P8     sql.NullString
	P9     sql.NullString
	P10    sql.NullString
	P11    sql.NullString
	P12    sql.NullString
	P13    sql.NullString
	P14    sql.NullString
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
		tableNames = append(tableNames, fmt.Sprintf("%s (%s)", column1, column2))
	}
	if err := rows.Err(); err != nil {
		log.Println("Error iterating over rows:", err)
		return nil, err
	}
	return tableNames, nil
}

// GetHeader ...
func (v *ViewingFirebird) GetHeader(tableName string) ([]string, error) {
	fc := "GetHeader"
	var tableCols []string
	query := "SELECT NAME FROM SCOL WHERE TABL = ? AND COL != 0"
	rows, err := v.db.Query(query, tableName)

	if err != nil {
		log.Printf("%s error occurred while querying database: %v", fc, err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var column string
		if err := rows.Scan(&column); err != nil {
			log.Printf("%s error occurred while scanning rows: %v", fc, err)
			return nil, err
		}
		tableCols = append(tableCols, column)
	}
	if err := rows.Err(); err != nil {
		log.Printf("%s error occurred while iterating over rows: %v", fc, err)
		return nil, err
	}
	return tableCols, nil
}

// GetIndicatorMaket ...
func (v *ViewingFirebird) GetIndicatorNumbers(tableName string) ([]models.IndicatorData, error) {
	fc := "GetIndicatorNumbers"
	var indicatorsRowsFirebird []IndicatorDataFirebird

	query := "SELECT * FROM PMAKET WHERE TABL = ?"
	rows, err := v.db.Query(query, tableName)

	if err != nil {
		log.Printf("%s error occurred while querying database: %v\n", fc, err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var data IndicatorDataFirebird
		if err = rows.Scan(&data.TABL,
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
			log.Printf("%s error occurred while scanning rows: %v\n", fc, err)
			return nil, err
		}
		indicatorsRowsFirebird = append(indicatorsRowsFirebird, data)
	}
	if err = rows.Err(); err != nil {
		log.Printf("%s error occurred while iterating over rows: %v\n", fc, err)
		return nil, err
	}
	log.Printf("the indicators' numbers fetched from the DB are: %v\n", indicatorsRowsFirebird)
	return handleNullString(indicatorsRowsFirebird), nil
}

// GetIndicator ...
func (v *ViewingFirebird) GetIndicator(shifr, npokaz string) string {
	fc := "GetIndicator"

	if npokaz == "" {
		return ""
	}
	var mnojitNum int
	var ngrup, chislit, znamet, mnojitStr string
	query := "SELECT NGRUP, CHISLIT, ZNAMET, MNOJIT FROM SPPOK WHERE SHIFR = ? AND NPOKAZ = ?"
	err := v.db.QueryRow(query, shifr, npokaz).Scan(&ngrup, &chislit, &znamet, &mnojitNum)
	if err != nil {
		log.Printf("%s error occurred while querying database: %v\n", fc, err)
		return ""
	}

	if mnojitNum == 1000 {
		mnojitStr = "1K"
	} else if mnojitNum == 10000 {
		mnojitStr = "10K"
	} else if mnojitNum == 100000 {
		mnojitStr = "100K"
	} else {
		mnojitStr = strconv.Itoa(mnojitNum)
	}

	return fmt.Sprintf("#%s (NGRUP: %s)\n (%s)*%s/\n(%s)", npokaz, ngrup, chislit, mnojitStr, znamet)
}

func handleNullString(indicatorsNumbers []IndicatorDataFirebird) []models.IndicatorData {
	processedData := make([]models.IndicatorData, len(indicatorsNumbers))

	for i, numbersRow := range indicatorsNumbers {
		processedRow := models.IndicatorData{}
		numbersRow.TABL = processedRow.TABL
		numbersRow.NZAP = processedRow.NZAP

		if !numbersRow.STRTAB.Valid {
			processedRow.STRTAB = ""
		} else {
			processedRow.STRTAB = numbersRow.STRTAB.String
		}

		if !numbersRow.SHSTR.Valid {
			processedRow.SHSTR = ""
		} else {
			processedRow.SHSTR = numbersRow.SHSTR.String
		}
		if !numbersRow.FORMA.Valid {
			processedRow.FORMA = ""
		} else {
			processedRow.FORMA = numbersRow.FORMA.String
		}
		if !numbersRow.STROKA.Valid {
			processedRow.STROKA = ""
		} else {
			processedRow.STROKA = numbersRow.STROKA.String
		}
		if !numbersRow.P1.Valid {
			processedRow.P1 = ""
		} else {
			processedRow.P1 = numbersRow.P1.String
		}
		if !numbersRow.P2.Valid {
			processedRow.P2 = ""
		} else {
			processedRow.P2 = numbersRow.P2.String
		}
		if !numbersRow.P3.Valid {
			processedRow.P3 = ""
		} else {
			processedRow.P3 = numbersRow.P3.String
		}
		if !numbersRow.P4.Valid {
			processedRow.P4 = ""
		} else {
			processedRow.P4 = numbersRow.P4.String
		}
		if !numbersRow.P5.Valid {
			processedRow.P5 = ""
		} else {
			processedRow.P5 = numbersRow.P5.String
		}
		if !numbersRow.P6.Valid {
			processedRow.P6 = ""
		} else {
			processedRow.P6 = numbersRow.P6.String
		}
		if !numbersRow.P7.Valid {
			processedRow.P7 = ""
		} else {
			processedRow.P7 = numbersRow.P7.String
		}
		if !numbersRow.P8.Valid {
			processedRow.P8 = ""
		} else {
			processedRow.P8 = numbersRow.P8.String
		}
		if !numbersRow.P9.Valid {
			processedRow.P9 = ""
		} else {
			processedRow.P9 = numbersRow.P9.String
		}
		if !numbersRow.P10.Valid {
			processedRow.P10 = ""
		} else {
			processedRow.P10 = numbersRow.P10.String
		}
		if !numbersRow.P11.Valid {
			processedRow.P11 = ""
		} else {
			processedRow.P11 = numbersRow.P11.String
		}
		if !numbersRow.P12.Valid {
			processedRow.P12 = ""
		} else {
			processedRow.P12 = numbersRow.P12.String
		}
		if !numbersRow.P13.Valid {
			processedRow.P13 = ""
		} else {
			processedRow.P13 = numbersRow.P13.String
		}
		if !numbersRow.P14.Valid {
			processedRow.P14 = ""
		} else {
			processedRow.P14 = numbersRow.P14.String
		}
		processedData[i] = processedRow
	}
	log.Printf("the indicators' numbers after processing NULLs are: %v\n", processedData)
	return processedData
}
