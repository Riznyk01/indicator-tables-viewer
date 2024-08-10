package main

import (
	"bufio"
	"fmt"
	"indicator-tables-viewer/internal/models"
	"os"
	"strings"
)

// filename string
func LoadTranslations(lng string) (models.Translations, error) {
	file, err := os.Open(fmt.Sprintf("lang_%s.txt", lng))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	translations := make(models.Translations)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		translations[parts[0]] = parts[1]
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return translations, nil
}
