package repository

import (
	"database/sql"
	"fmt"
	"indicator-tables-viewer/internal/config"
	"log"
)

func NewFirebirdDB(cfg *config.Config, login, pass string, local bool, yearPeriod bool) (*sql.DB, string, error) {
	var connectionString, connectionStringWithoutPass, localPathToDb, dbDir string
	passwordReplacementText := "******"
	if yearPeriod {
		if local {
			dbDir = "\\" + cfg.LocalYearDbDir
		} else {
			dbDir = cfg.RemoteYearDbDir
		}
	} else {
		if local {
			dbDir = "\\" + cfg.LocalQuarterDbDir
		} else {
			dbDir = cfg.RemoteQuarterDbDir
		}
	}

	if cfg.Env == "dev" {
		localPathToDb = cfg.CodePath
	} else if cfg.Env == "prod" {
		localPathToDb = cfg.LocalPath
	}

	if local {
		connectionString = fmt.Sprintf("%s:%s@%s:%s/%s/%s",
			login, pass, cfg.LocalHost, cfg.LocalPort, localPathToDb+dbDir, cfg.DBName)

		connectionStringWithoutPass = fmt.Sprintf("%s:%s@%s:%s/%s/%s",
			login, passwordReplacementText, cfg.LocalHost, cfg.LocalPort, localPathToDb+dbDir, cfg.DBName)
	} else {
		connectionString = fmt.Sprintf("%s:%s@%s:%s/%s/%s",
			login, pass, cfg.Host, cfg.Port, cfg.RemotePathToDb+dbDir, cfg.DBName)
		connectionStringWithoutPass = fmt.Sprintf("%s:%s@%s:%s/%s/%s",
			login, passwordReplacementText, cfg.Host, cfg.Port, cfg.RemotePathToDb+dbDir, cfg.DBName)
	}

	log.Printf("connection string is: %s\n", connectionStringWithoutPass)
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
	return db, connectionStringWithoutPass, nil
}
