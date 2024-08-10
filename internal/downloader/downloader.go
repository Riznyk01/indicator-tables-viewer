package downloader

import (
	"fmt"
	"indicator-tables-viewer/internal/models"
	"io"
	"log"
	"net/http"
	"os"
)

func Download(URL, fileName, destinationPath string, lang models.Translations) error {
	fc := "Download"
	fmt.Printf("%s downloading started: URL=%s, fileName=%s, destinationPath=%s\n",
		fc, URL, fileName, destinationPath)

	resp, err := http.Get(URL + fileName)
	if err != nil {
		log.Printf("%s %s: %v", lang["ErrOccur"], lang["DownUpd"], err)
		return err
	}
	defer resp.Body.Close()

	file, err := os.Create(destinationPath)
	if err != nil {
		log.Printf("%s %s: %v", lang["ErrOccur"], lang["FileCreating"], err)
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		log.Printf("%s copying file to the dir: %v", lang["ErrOccur"], err)
		return err
	}
	return nil
}
