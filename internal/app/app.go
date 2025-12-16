package app

import (
	"log"
	"wishlist-bot/pkg/database"
	"wishlist-bot/internal/user"
	"wishlist-bot/internal/wishlist"
	"wishlist-bot/internal/bot"
	"wishlist-bot/internal/fsm"
	"wishlist-bot/internal/scheduler"
)

func Start() {
	db, err := database.Init()
	if err != nil {
		log.Printf("Error with db connection: %s", err)
	}
	ur := user.NewRepository(db)
	wr := wishlist.NewRepository(db)

	us := user.NewService(ur)
	ws := wishlist.NewService(wr)

	states := fsm.NewInMemoryStateStore()

	uRouter := bot.NewUserHandler(us, states)
	wRouter := bot.NewWishlistHandler(ws, states)
	aRouter := bot.NewAdminHandler(us, ws, states)

	mainRouter := bot.NewHandlerRouter(uRouter, wRouter, aRouter, states)
	bot, err := bot.NewBot(mainRouter)
	if err != nil {
		log.Printf("Error with create bot: %w", err)
	}

	sch := scheduler.NewScheduler(bot.API(), us, ws)
	go sch.StartScheduler()

	bot.RegisterHandlers()
	bot.Start()
	select {}
}
