package app

import (
	"fmt"
	"log"
	"log/slog"
	"wishlist-bot/internal/bot"
	"wishlist-bot/internal/config"
	"wishlist-bot/internal/fsm"
	"wishlist-bot/internal/scheduler"
	"wishlist-bot/internal/user"
	"wishlist-bot/internal/wishlist"
	"wishlist-bot/pkg/database"
)

type App struct {
	cfg    *config.Config
	log    *slog.Logger
	botApp *bot.App
}

func New(cfg *config.Config, log *slog.Logger) *App {
	return &App{
		cfg: cfg,
		log: log,
	}
}

func (a *App) MustStart() {
	db, err := database.Init()
	if err != nil {
		panic(fmt.Errorf("error with db connection: %w", err))
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

func (a *App) Stop() {
	const op = "app.Stop"
	a.log.Info(op)

}
