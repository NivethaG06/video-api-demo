package repository

import (
	"database/sql"
	"fmt"
	"os"
)

func GetConnection() (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"))
	return sql.Open("postgres", psqlInfo)
}

type PostgresDB struct {
	Sql_db *sql.DB
}

func LoadConnection() (*PostgresDB, error) {
	db, err := GetConnection()
	if err != nil {
		return nil, fmt.Errorf("failed to open the repository: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to the repository: %w", err)
	}

	return &PostgresDB{Sql_db: db}, nil
}
