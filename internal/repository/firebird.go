package repository

import (
	"database/sql"
	"fmt"
	"indicator-tables-viewer/internal/config"
	"log"
)

func NewFirebirdDB(cfg *config.Config, login, pass string, local bool) (*sql.DB, string, error) {
	var connectionString string
	var localPathToDb string

	if cfg.Env == "dev" {
		localPathToDb = cfg.CodePath
	} else if cfg.Env == "prod" {
		localPathToDb = cfg.LocalPath
	}

	if local {
		connectionString = fmt.Sprintf("%s:%s@%s:%s/%s/%s",
			login, pass, cfg.LocalHost, cfg.LocalPort, localPathToDb, cfg.DBName)
	} else {
		connectionString = fmt.Sprintf("%s:%s@%s:%s/%s/%s",
			login, pass, cfg.Host, cfg.Port, cfg.RemotePathToDb, cfg.DBName)
	}

	log.Printf("connection string is: %s\n", connectionString)
	db, err := sql.Open("firebirdsql", connectionString)

	if err != nil {
		log.Println("Error connecting to database:", err)
		return nil, "", err
	}

	err = db.Ping()
	if err != nil {
		log.Println("Error pinging database:", err)
		return nil, "", err
	}
	return db, connectionString, nil
}
