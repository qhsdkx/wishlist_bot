package bot

import (
	"fmt"

	"gopkg.in/telebot.v4"
)

func makeDeleteButton(id int64, title string) telebot.InlineButton {
	return telebot.InlineButton{
		Unique: "del_wish",
		Text:   title,
		Data:   fmt.Sprintf("%d", id),
	}
}
