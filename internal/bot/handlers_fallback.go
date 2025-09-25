package bot

import "gopkg.in/telebot.v4"

func (b *Bot) handleUnknown(c telebot.Context) error {
	return c.Send("Не понял команду 🤔")
}
