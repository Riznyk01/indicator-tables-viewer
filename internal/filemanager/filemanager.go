package filemanager

import (
	"errors"
	"fmt"
	"github.com/tealeg/xlsx"
	"log"
	"os"
	"path/filepath"
	"time"
)

func ExportToExcel(data [][]string, tableName string, exportPath string) error {
	file := xlsx.NewFile()
	sheet, err := file.AddSheet("Sheet1")
	if err != nil {
		return errors.New(fmt.Sprintf("Error creating sheet: %v\n", err))
	}

	for _, row := range data {
		newRow := sheet.AddRow()
		for _, cell := range row {
			newCell := newRow.AddCell()
			newCell.Value = cell
			newCell.GetStyle().Alignment.WrapText = true
		}
	}

	log.Println("File saved successfully.")
	for s := 0; s < sheet.MaxCol; s++ {
		sheet.Col(s).Width = float64(30)
	}

	currentTime := time.Now()
	currentDateTime := currentTime.Format("2006-01-02_15-04-05")
	filename := tableName + "_" + currentDateTime + ".xlsx"

	var fullPath string

	if exportPath == "" {
		fullPath = filename
	} else {
		fullPath = filepath.Join(exportPath, filename)
	}

	err = file.Save(fullPath)
	if err != nil {
		return errors.New(fmt.Sprintf("error occurred while saving xls file: %v\n", err))
	}
	return nil
}

// CheckLogFileSize ...
func CheckLogFileSize(logFilePath string, size int64) error {
	fileInfo, err := os.Stat(logFilePath)
	if err != nil {
		return err
	}
	if fileInfo.Size() > size {
		if err = os.Remove(logFilePath); err != nil {
			return err
		}
	}
	return nil
}
