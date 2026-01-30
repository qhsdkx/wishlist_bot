package bot

import (
	"errors"
	"gopkg.in/telebot.v4"
	"log"
	"time"
	"wishlist-bot/internal/config"
)

type Bot struct {
	tg     *telebot.Bot
	router HandlerRouter
}

func New(router HandlerRouter, cfg config.BotConfig) (*Bot, error) {
	if cfg.ApiKey == "" {
		return nil, errors.New("token is empty")
	}

	pref := telebot.Settings{
		Token:  cfg.ApiKey,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}
	tgBot, err := telebot.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &Bot{
		tg:     tgBot,
		router: router,
	}, nil
}

func (b *Bot) RegisterHandlers() {
	b.tg.Handle("/start", b.router.OnStart, b.router.checkSheluvssic())
	b.tg.Handle(telebot.OnText, b.router.OnText, b.router.checkSheluvssic())
	b.tg.Handle(telebot.OnCallback, b.router.OnCallback, b.router.checkSheluvssic())
	b.tg.Handle("/help", b.router.Help, b.router.checkSheluvssic())
}

func (b *Bot) Start() {
	b.tg.Start()
}

func (b *Bot) API() *telebot.Bot {
	return b.tg
}

func (b *Bot) Stop() {
	b.tg.Stop()
}
