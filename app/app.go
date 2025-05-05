package app

import (
	"fmt"
	b "wishlist-bot/bot"
	"wishlist-bot/database"
	"wishlist-bot/scheduler"
	"wishlist-bot/service"
)

func Start() {
	db, dbErr := database.Init()
	if dbErr != nil {
		fmt.Errorf("Error initializing database: %v", dbErr)
	}
	userRepository := &database.UserRepositoryImpl{DB: db}
	userService := &service.UserServiceImpl{Repo: userRepository}
	bot := b.SetUp(userService)
	go func() {
		scheduler.StartScheduler(bot, userService)
	}()
	bot.Start()
	select {}
}
