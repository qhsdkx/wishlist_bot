package bot

import "gopkg.in/telebot.v4"

func (r *HandlerRouter) handleUnknown(c telebot.Context) error {
	return c.Send("Не понял команду 🤔")
}
