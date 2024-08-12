package repository

import (
	"database/sql"
	"fmt"
	"indicator-tables-viewer/internal/models"
	"log"
	"strconv"
	"time"
)

const (
	empty = " "
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
	rows, err := v.db.Query("SELECT tabl, nazv FROM STABLES WHERE tabl LIKE  'P%' ORDER BY tabl")
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

// GetColNameLocation ...
func (v *ViewingFirebird) GetColNameLocation(tableName string) (string, error) {
	fc := "GetColNameLocation"
	var nameLocation string

	query := "SELECT NSH FROM STABLES WHERE TABL = ?"
	err := v.db.QueryRow(query, tableName).Scan(&nameLocation)
	if err != nil {
		log.Printf("%s error occurred while querying database: %v\n", fc, err)
		return "", err
	}
	return nameLocation, nil
}

// GetHeader ...
func (v *ViewingFirebird) GetHeader(tableName string) ([]string, error) {
	fc := "GetHeader"
	var tableCols []string
	query := "SELECT NAME FROM SCOL WHERE TABL = ? AND COL != 0 ORDER BY COL"
	rows, err := v.db.Query(query, tableName)

	if err != nil {
		log.Printf("%s error occurred while querying database: %v", fc, err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var columnStr string
		var column sql.NullString
		if err := rows.Scan(&column); err != nil {
			log.Printf("%s error occurred while scanning rows: %v", fc, err)
			return nil, err
		}
		if !column.Valid {
			columnStr = "name is empty"
		} else {
			columnStr = column.String
		}
		tableCols = append(tableCols, columnStr)
	}
	if err := rows.Err(); err != nil {
		log.Printf("%s error occurred while iterating over rows: %v", fc, err)
		return nil, err
	}
	fmt.Printf("%s: header is fetched: %s\n", fc, tableCols)
	return tableCols, nil
}

// GetIndicatorMaket ...
func (v *ViewingFirebird) GetIndicatorNumbers(tableName string) ([]models.IndicatorData, error) {
	fc := "GetIndicatorNumbers"
	var indicatorsRowsFirebird []IndicatorDataFirebird

	query := "SELECT * FROM PMAKET WHERE TABL = ?"
	rows, err := v.db.Query(query, tableName)
	if err != nil {
		log.Printf("%s: %v\n", fc, err)
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
	log.Printf("%s: the indicators' numbers fetched from the DB are: %v\n", fc, indicatorsRowsFirebird)
	IndicatorNumbers := handleNullString(indicatorsRowsFirebird)
	return IndicatorNumbers, nil
}

// GetIndicator ...
func (v *ViewingFirebird) GetIndicator(shifr, npokaz, indicatorsRow, decodingRow, decodingTable string) string {
	fc := "GetIndicator"

	if npokaz == empty {
		return empty
	}
	var mnojitNum int
	var ngrup, chislit, znamet, mnojitStr string
	query := "SELECT NGRUP, CHISLIT, ZNAMET, MNOJIT FROM SPPOK WHERE SHIFR = ? AND NPOKAZ = ?"
	err := v.db.QueryRow(query, shifr, npokaz).Scan(&ngrup, &chislit, &znamet, &mnojitNum)
	if err != nil {
		log.Printf("%s: %v\n", fc, err)
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

	if mnojitStr != "1" {
		mnojitStr = fmt.Sprintf("множник %s", mnojitStr)
	} else {
		mnojitStr = ""
	}

	if znamet != "1" {
		znamet = fmt.Sprintf("знаменник %s", znamet)
	} else {
		znamet = ""
	}

	return fmt.Sprintf("ГР%s №%s рядок табл. показника %s\nрядок розшифровки. %s з табл. %s\nчисельник %s\n%s\n%s", ngrup, npokaz, indicatorsRow, decodingRow, decodingTable, chislit, mnojitStr, znamet)
}

func handleNullString(indicatorsNumbers []IndicatorDataFirebird) []models.IndicatorData {
	fc := "handleNullString"
	processedData := make([]models.IndicatorData, len(indicatorsNumbers))

	for i, numbersRow := range indicatorsNumbers {
		log.Printf("%s: %v", fc, numbersRow)
		processedRow := models.IndicatorData{}
		processedRow.TABL = numbersRow.TABL
		processedRow.NZAP = numbersRow.NZAP
		if !numbersRow.STRTAB.Valid {
			processedRow.STRTAB = empty
		} else {
			processedRow.STRTAB = numbersRow.STRTAB.String
		}
		if !numbersRow.SHSTR.Valid {
			processedRow.SHSTR = empty
		} else {
			processedRow.SHSTR = numbersRow.SHSTR.String
		}
		if !numbersRow.FORMA.Valid {
			processedRow.FORMA = empty
		} else {
			processedRow.FORMA = numbersRow.FORMA.String
		}
		if !numbersRow.STROKA.Valid {
			processedRow.STROKA = empty
		} else {
			processedRow.STROKA = numbersRow.STROKA.String
		}
		if !numbersRow.P1.Valid {
			processedRow.P1 = empty
		} else {
			processedRow.P1 = numbersRow.P1.String
		}
		if !numbersRow.P2.Valid {
			processedRow.P2 = empty
		} else {
			processedRow.P2 = numbersRow.P2.String
		}
		if !numbersRow.P3.Valid {
			processedRow.P3 = empty
		} else {
			processedRow.P3 = numbersRow.P3.String
		}
		if !numbersRow.P4.Valid {
			processedRow.P4 = empty
		} else {
			processedRow.P4 = numbersRow.P4.String
		}
		if !numbersRow.P5.Valid {
			processedRow.P5 = empty
		} else {
			processedRow.P5 = numbersRow.P5.String
		}
		if !numbersRow.P6.Valid {
			processedRow.P6 = empty
		} else {
			processedRow.P6 = numbersRow.P6.String
		}
		if !numbersRow.P7.Valid {
			processedRow.P7 = empty
		} else {
			processedRow.P7 = numbersRow.P7.String
		}
		if !numbersRow.P8.Valid {
			processedRow.P8 = empty
		} else {
			processedRow.P8 = numbersRow.P8.String
		}
		if !numbersRow.P9.Valid {
			processedRow.P9 = empty
		} else {
			processedRow.P9 = numbersRow.P9.String
		}
		if !numbersRow.P10.Valid {
			processedRow.P10 = empty
		} else {
			processedRow.P10 = numbersRow.P10.String
		}
		if !numbersRow.P11.Valid {
			processedRow.P11 = empty
		} else {
			processedRow.P11 = numbersRow.P11.String
		}
		if !numbersRow.P12.Valid {
			processedRow.P12 = empty
		} else {
			processedRow.P12 = numbersRow.P12.String
		}
		if !numbersRow.P13.Valid {
			processedRow.P13 = empty
		} else {
			processedRow.P13 = numbersRow.P13.String
		}
		if !numbersRow.P14.Valid {
			processedRow.P14 = empty
		} else {
			processedRow.P14 = numbersRow.P14.String
		}
		processedData[i] = processedRow
	}
	log.Printf("%s: the indicators' numbers after processing NULLs are: %v\n", fc, processedData)
	return processedData
}

func (v *ViewingFirebird) UpdateDBCorrectionDate(currentTime time.Time) error {
	fc := "UpdateDBCorrectionDate"

	query := "UPDATE V_PARAM_VALUES SET F_DATE_VAL = ? WHERE F_PARAM_NAME = ?"

	_, err := v.db.Exec(query, currentTime, "GBackDate")
	if err != nil {
		log.Printf("%s: %v\n", fc, err)
		return err
	}
	return nil
}
