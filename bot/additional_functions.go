package bot

import (
	"fmt"
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

func checkSheluvssic() telebot.MiddlewareFunc {
	return func(next telebot.HandlerFunc) telebot.HandlerFunc {
		return func(c telebot.Context) error {
			if c.Chat().ID == 420845081 {
				return c.Send(fmt.Sprintf("@sheluvssic ловит бээээээээу в свой вазилиновый сракатан\nБББББББЭЭЭЭЭЭЭЭЭЭЭУУУУУУУУУУУУУУУУУ\n\n**СТРИПКЛУБ БЛЭКЛИСТ ПАЦАНЧИКИ**"))
			}
			return next(c)
		}
	}
}

func onError(c telebot.Context) error {
	delete(states, c.Chat().ID)
	return c.Send("Неизвестное состояние. Пожалуйста, начните заново с команды /start.")
}
