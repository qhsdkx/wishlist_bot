package database

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"os"
	"strconv"
	"time"
)

var DB *sql.DB

func Init() (*sql.DB, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	port, _ := strconv.Atoi(os.Getenv("PORT"))
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("HOST"), port, os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))

	DB, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	conns, _ := strconv.Atoi(os.Getenv("MAX_DB_COONECTIONS"))

	DB.SetMaxOpenConns(conns)
	DB.SetMaxIdleConns(5)
	DB.SetConnMaxLifetime(time.Minute * 30)

	return DB, DB.Ping()
}

func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}

func beginTransaction(db *sql.DB) *sql.Tx {
	tx, err := db.Begin()
	if err != nil {
		fmt.Errorf("error beginning transaction: %v", err)
	}
	return tx
}

func rollbackTransaction(tx *sql.Tx) {
	err := tx.Rollback()
	if err != nil {
		fmt.Errorf("error rolling back transaction: %v", err)
	}
}

func commitTransaction(tx *sql.Tx) {
	err := tx.Commit()
	if err != nil {
		fmt.Errorf("error committing transaction: %v", err)
	}
}
