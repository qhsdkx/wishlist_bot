package database

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"os"
	"strconv"
)

var db *sql.DB

func Init() (*sql.DB, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	port, _ := strconv.Atoi(os.Getenv("PORT"))
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("HOST"), port, os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))

	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	conns, _ := strconv.Atoi(os.Getenv("MAX_DB_COONECTIONS"))

	db.SetMaxOpenConns(conns)
	db.SetMaxIdleConns(5)

	return db, db.Ping()
}
