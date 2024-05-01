package logger

import (
	"fmt"
	"log"
	"os"
)

// CheckLogFile checks and creates log file if not exist by LogFileCreate function
func CheckLogFile(filepath string) error {
	_, err := os.Stat(filepath)
	if err != nil {
		if os.IsNotExist(err) {
			err = LogFileCreate(filepath)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	return nil
}

// LogFileCreate creates log file
func LogFileCreate(filepath string) error {
	newFile, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	newFile.Close()
	return nil
}

// CheckLogFileSize ...
func CheckLogFileSize(pathForViewer, pathForLauncher string, size int64) error {
	fileInfo, err := os.Stat(pathForViewer)
	if err != nil {
		return err
	}

	if fileInfo.Size() > size {
		if err = os.Remove(pathForViewer); err != nil {
			return err
		}
		if err = os.Remove(pathForLauncher); err != nil {
			return err
		}
		err = LogFileCreate(pathForViewer)
		if err != nil {
			return err
		}
		err = LogFileCreate(pathForLauncher)
		if err != nil {
			return err
		}
	}
	return nil
}
func OpenLogFile(filepath string) (*os.File, error) {
	logFile, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("error occurred while opening logfile:", err)
		return nil, nil
	}
	return logFile, nil
}
