package database

import (
	"database/sql"
	"fmt"
	"strconv"
	"wishlist-bot/internal/config"

	_ "github.com/lib/pq"
)

func MustInit(cfg config.DatabaseConfig) *sql.DB {
	port, _ := strconv.Atoi(cfg.Port)
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, port, cfg.User,
		cfg.Password, cfg.Name)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	db.SetMaxOpenConns(cfg.MaxDBConns)
	db.SetMaxIdleConns(cfg.MaxDBConns)

	return db
}
