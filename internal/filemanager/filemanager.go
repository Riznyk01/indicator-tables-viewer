package filemanager

import (
	"archive/zip"
	"errors"
	"fmt"
	"github.com/tealeg/xlsx"
	"io"
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

func MakeDirIfNotExist(path string) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		err := os.MkdirAll(path, 0755)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	return nil
}
func Unzip(zipFile string, dest string) error {

	r, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		path := dest + "\\" + f.Name

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			f, err := os.OpenFile(
				path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer f.Close()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
