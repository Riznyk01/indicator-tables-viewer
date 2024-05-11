package repository

import (
	"database/sql"
	"fmt"
	"indicator-tables-viewer/internal/config"
	"log"
)

func NewFirebirdDB(cfg *config.Config, login, pass string) (*sql.DB, string, error) {
	var connectionString, connectionStringWithoutPass, localPathToDb, dbDir string
	passwordReplacementText := "******"
	if cfg.YearDB {
		if cfg.LocalMode {
			dbDir = "\\" + cfg.LocalYearDbDir
		} else {
			dbDir = cfg.RemoteYearDbDir
		}
	} else {
		if cfg.LocalMode {
			dbDir = "\\" + cfg.LocalQuarterDbDir
		} else {
			dbDir = cfg.RemoteQuarterDbDir
		}
	}

	localPathToDb = cfg.LocalPath

	if cfg.LocalMode {
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
