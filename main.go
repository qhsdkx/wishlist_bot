package main

import (
	"database/sql"
	"fmt"
	"log"
	"wishlist-bot/bot"
	"wishlist-bot/database"
	"wishlist-bot/service"
)

func main() {
	var db *sql.DB = nil
	db, dbErr := database.Init()
	if dbErr != nil {
		fmt.Errorf("Error initializing database: %v", dbErr)
		database.Close()
	}
	userRepository := &database.UserRepositoryImpl{DB: db}
	userService := service.UserServiceImpl{Repo: userRepository}
	telebot, err := bot.NewBot()
	if err != nil {
		log.Fatal(err)
	}
	bot.Start(telebot, &userService)
}
