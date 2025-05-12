package bot

import (
	"gopkg.in/telebot.v4"
	constants "wishlist-bot/constant"
	sv "wishlist-bot/service"
)

func updateUserListPage(c telebot.Context, page int, userService sv.UserService) error {
	users, pagination, err := userService.FindAll(page, constants.USERS_PER_PAGE)
	if err != nil {
		return c.Respond(&telebot.CallbackResponse{
			Text: "Ошибка обновления списка",
		})
	}

	markup := createUserListMarkup(users, pagination)
	_, err = bot.Edit(c.Message(), "Список пользователей:", markup)
	if err != nil {
		return c.Respond(&telebot.CallbackResponse{
			Text: "Ошибка обновления",
		})
	}

	return c.Respond()
}

func createBackButton() *telebot.ReplyMarkup {
	markup := &telebot.ReplyMarkup{}
	backBtn := markup.Data("Назад", constants.BACK_TO_LIST)
	markup.Inline(markup.Row(backBtn))
	return markup
}

func checkDeleted(service sv.UserService) telebot.MiddlewareFunc {
	return func(next telebot.HandlerFunc) telebot.HandlerFunc {
		return func(c telebot.Context) error {
			if service.CheckIfDeleted(c.Chat().ID) && c.Callback().Data[1:] != constants.BTN_RESTORE_USER {
				_, err := bot.Edit(c.Message(), "Вы удалены. Доступны следующие действия", deletedSelector)
				return err
			}
			return next(c)
		}
	}
}

func onError(c telebot.Context) error {
	delete(states, c.Chat().ID)
	return c.Send("Неизвестное состояние. Пожалуйста, начните заново с команды /start.")
}
