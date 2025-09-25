package bot

import "gopkg.in/telebot.v4"

func (b *Bot) handleHelp(c telebot.Context) error {
	return c.Send("Доступные команды: /start /help")
}
