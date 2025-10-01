package bot

import (
	"gopkg.in/telebot.v4"
)

type Bot struct {
	tg     *telebot.Bot
	router HandlerRouter
}

func NewBot(tg *telebot.Bot, router HandlerRouter) *Bot {
	return &Bot{tg: tg, router: router}
}

func (b *Bot) RegisterHandlers() {
	b.tg.Handle("/start", b.router.OnStart)
	b.tg.Handle(telebot.OnText, b.router.OnText)
	b.tg.Handle(telebot.OnCallback, b.router.OnCallback)
	b.tg.Handle("/help", b.router.Help)
}

func (b *Bot) Start() {
	b.tg.Start()
}
