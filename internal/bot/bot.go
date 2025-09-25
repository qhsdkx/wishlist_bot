package bot

import (
	"fmt"
	"wishlist-bot/internal/user"

	"gopkg.in/telebot.v4"
)

type Bot struct {
	tg    *telebot.Bot
	users *user.Service
}

func NewBot(tg *telebot.Bot, users *user.Service) *Bot {
	return &Bot{tg: tg, users: users}
}

func (b *Bot) RegisterHandlers() {
	b.tg.Handle("/start", func(c telebot.Context) error {
		u, err := b.users.GetUser(c.Chat().ID)
		if err != nil {
			return c.Send("Привет! Давай зарегистрируемся?")
		}
		return c.Send(fmt.Sprintf("С возвращением, %s!", u.Username))
	})
}
