package bot

import (
	constants "wishlist-bot/internal/constant"
	u "wishlist-bot/internal/user"

	"gopkg.in/telebot.v4"
)

func updateUserListPage(c telebot.Context, page int, userService u.Service, mode string) error {
	users, pagination, err := userService.FindAll(page, constants.USERS_PER_PAGE, mode)
	if err != nil {
		return c.Respond(&telebot.CallbackResponse{
			Text: "Ошибка обновления списка",
		})
	}

	markup := createUserListMarkup(users, pagination, mode)
	if mode == constants.SHOW_USERS {
		return c.Edit("Список пользователей:\n", markup)
	}
	return c.Edit("Список пользователей:\n", markup)
}

func createBackButton() *telebot.ReplyMarkup {
	markup := &telebot.ReplyMarkup{}
	backBtn := markup.Data("Назад", constants.BACK_TO_LIST)
	markup.Inline(markup.Row(backBtn))
	return markup
}

func checkSheluvssic() telebot.MiddlewareFunc {
	return func(next telebot.HandlerFunc) telebot.HandlerFunc {
		return func(c telebot.Context) error {
			if c.Chat().ID == 420845081 || c.Chat().ID == 6466693361 {
				return c.Send("@sheluvssic и его фейк ловит бээээээээу в свой вазилиновый сракатан\nБББББББЭЭЭЭЭЭЭЭЭЭЭУУУУУУУУУУУУУУУУУ\n\nСТРИПКЛУБ БЛЭКЛИСТ ПАЦАНЧИКИ")
			}
			return next(c)
		}
	}
}

func onError(c telebot.Context) error {
	delete(states, c.Chat().ID)
	return c.Send("Неизвестное состояние. Пожалуйста, начните заново с команды /start.")
}
