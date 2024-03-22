package repository

import (
	"database/sql"
	"fmt"
	indicator_tables_viewer "indicator-tables-viewer"
)

func NewFirebirdDB(cfg *indicator_tables_viewer.Config) (*sql.DB, error) {
	connectionString := fmt.Sprintf("%s:%s@%s:%s/%s/%s",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Path, cfg.DBName)
	fmt.Printf("connection string is: %s\n", connectionString)
	db, err := sql.Open("firebirdsql", connectionString)

	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		fmt.Println("Error pinging database:", err)
		return nil, err
	}
	return db, nil
}
