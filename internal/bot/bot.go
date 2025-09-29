package bot

import (
	"wishlist-bot/internal/user"

	"gopkg.in/telebot.v4"
)

type Bot struct {
	tg     *telebot.Bot
	users  *user.Service
	router HandlerRouter
}

func NewBot(tg *telebot.Bot, users *user.Service, router HandlerRouter) *Bot {
	return &Bot{tg: tg, users: users, router: router}
}

func (b *Bot) RegisterHandlers() {
	b.tg.Handle("/start", b.router.OnStart, checkSheluvssic())
	b.tg.Handle(telebot.OnCallback, b.router.OnCallback, checkSheluvssic())
	b.tg.Handle(telebot.OnText, b.router.OnText, checkSheluvssic())
	b.tg.Handle("/help", b.router.userHandler.Help, checkSheluvssic())
}
