package app

import (
	"log/slog"
	"wishlist-bot/internal/bot"
	"wishlist-bot/internal/config"
	"wishlist-bot/internal/fsm"
	"wishlist-bot/internal/logger/sl"
	"wishlist-bot/internal/scheduler"
	"wishlist-bot/internal/user"
	"wishlist-bot/internal/wishlist"
	"wishlist-bot/pkg/database"
)

type App struct {
	cfg       *config.Config
	log       *slog.Logger
	bot       *bot.Bot
	scheduler *scheduler.Scheduler
}

func New(cfg *config.Config, log *slog.Logger) *App {
	return &App{
		cfg: cfg,
		log: log,
	}
}

func (a *App) MustStart() {
	db := database.MustInit(a.cfg.Database)

	ur := user.NewRepository(db)
	wr := wishlist.NewRepository(db)

	us := user.NewService(ur)
	ws := wishlist.NewService(wr)

	states := fsm.NewInMemoryStateStore()

	uRouter := bot.NewUserHandler(us, states)
	wRouter := bot.NewWishlistHandler(ws, states)
	aRouter := bot.NewAdminHandler(us, ws, states)
	mainRouter := bot.NewHandlerRouter(uRouter, wRouter, aRouter, states)

	botApi, err := bot.New(mainRouter, a.cfg.Bot)
	if err != nil {
		a.log.Error("Error creating botApi", sl.Err(err))
		panic(err)
	}

	a.bot = botApi

	a.scheduler = scheduler.New(a.bot.API(), us, ws, a.cfg)

	go a.scheduler.Start()

	a.bot.RegisterHandlers()

	go a.bot.Start()

	select {}
}

func (a *App) Stop() {
	const op = "app.Stop"
	a.log.Info(op)

	if a.scheduler != nil {
		a.scheduler.Stop()
	}
	if a.bot != nil {
		a.bot.Stop()
	}
}
