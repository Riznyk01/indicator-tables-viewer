package repository

import (
	"database/sql"
	"fmt"
	"indicator-tables-viewer/internal/config"
	"log"
)

func NewFirebirdDB(cfg *config.Config, login, pass string) (*sql.DB, error) {
	connectionString := fmt.Sprintf("%s:%s@%s:%s/%s/%s",
		login, pass, cfg.Host, cfg.Port, cfg.Path, cfg.DBName)
	log.Printf("connection string is: %s\n", connectionString)
	db, err := sql.Open("firebirdsql", connectionString)

	if err != nil {
		log.Println("Error connecting to database:", err)
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		log.Println("Error pinging database:", err)
		return nil, err
	}
	return db, nil
}
