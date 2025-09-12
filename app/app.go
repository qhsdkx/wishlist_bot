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
		fmt.Printf("Error initializing database: %v", dbErr)
	}
	wishlistRepository := database.NewWishlistRepository(db)
	wishlistService := service.NewWishService(wishlistRepository)
	userRepository := database.NewUserRepository(db)
	userService := service.NewUserService(userRepository)
	bot := b.SetUp(userService, wishlistService)
	scheduler := scheduler.NewScheduler(bot, userService, wishlistService)
	go func() {
		scheduler.StartScheduler()
	}()
	bot.Start()
	select {}
}
