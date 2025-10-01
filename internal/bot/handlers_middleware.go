package bot

import (
	"fmt"
	"strings"
	constants "wishlist-bot/internal/constant"

	"gopkg.in/telebot.v4"
)

func (h *UserHandler) updateUserListPage(c telebot.Context, page int, mode string) error {
	users, pagination, err := h.service.FindAll(page, constants.USERS_PER_PAGE, mode)
	if err != nil {
		return c.Respond(&telebot.CallbackResponse{
			Text: "Ошибка обновления списка",
		})
	}

	markup := h.createUserListMarkup(users, pagination, mode)
	if mode == constants.SHOW_USERS {
		return c.Edit("Список пользователей:\n", markup)
	}
	return c.Edit("Список пользователей:\n", markup)
}

// func updateUserListPage(c telebot.Context, page int, userService u.Service, mode string) error {
// 	users, pagination, err := userService.FindAll(page, constants.USERS_PER_PAGE, mode)
// 	if err != nil {
// 		return c.Respond(&telebot.CallbackResponse{
// 			Text: "Ошибка обновления списка",
// 		})
// 	}

// 	markup := createUserListMarkup(users, pagination, mode)
// 	if mode == constants.SHOW_USERS {
// 		return c.Edit("Список пользователей:\n", markup)
// 	}
// 	return c.Edit("Список пользователей:\n", markup)
// }

func (r *HandlerRouter) createBackButton() *telebot.ReplyMarkup {
	markup := &telebot.ReplyMarkup{}
	backBtn := markup.Data("Назад", constants.BACK_TO_LIST)
	markup.Inline(markup.Row(backBtn))
	return markup
}

func (r *HandlerRouter) checkSheluvssic() telebot.MiddlewareFunc {
	return func(next telebot.HandlerFunc) telebot.HandlerFunc {
		return func(c telebot.Context) error {
			if c.Chat().ID == 420845081 || c.Chat().ID == 6466693361 {
				return c.Send("@sheluvssic и его фейк ловит бээээээээу в свой вазилиновый сракатан\nБББББББЭЭЭЭЭЭЭЭЭЭЭУУУУУУУУУУУУУУУУУ\n\nСТРИПКЛУБ БЛЭКЛИСТ ПАЦАНЧИКИ")
			}
			return next(c)
		}
	}
}

func (r *HandlerRouter) Error(c telebot.Context) error {
	r.states.Delete(c.Chat().ID)
	return c.Send("Неизвестное состояние. Пожалуйста, начните заново с команды /start.")
}

func (r *HandlerRouter) Help(c telebot.Context) error {
	response := strings.Builder{}
	response.WriteString(fmt.Sprintf("Данная система была создана с целью помощи работникам ЦЦР (пока что 9-го департамента) следить за днями рождения коллег\n"))
	response.WriteString(fmt.Sprintf("ВНИМАНИЕ. ВСЕ предусмотренные уведомления приходят только в случае полной регистрации пользоваетеля по кнопке \"Регистрация\"\n"))
	response.WriteString(fmt.Sprintf("Короткая информация по возможностям бота:\n\n"))
	response.WriteString(fmt.Sprintf("• \"Редактировать мои данные\" - кнопка, представляющая возможность изменения введенных при регистрации данных. соответственно доступна только зарегистрированным\n"))
	response.WriteString(fmt.Sprintf("• \"Список желаний\" - дает возможность ввести ваши пожелания, удалить что-то либо посмотреть пожелания других\n"))
	response.WriteString(fmt.Sprintf("• \"Регистрация\" - необходима для регистрации вас в системе (ввод имени, фамилии и даты рождения)\n"))
	response.WriteString(fmt.Sprintf("• \"Показать всех пользователей\" - показывает всех ЗАРЕГИСТРИРОВАННЫХ пользователей и доступна только для них. По нажатию на кнопку с именем покажется день рождения человека и его пожелания\n"))
	response.WriteString(fmt.Sprintf("• \"Удалить меня в базе\" - полностью удаляет вас в базе. Далее необходимо следовать инструкции\n\n"))
	response.WriteString(fmt.Sprintf("Это было краткое описание основных возможностей бота. В случае возникающих проблем или предложений пишите разработчику @qhsdkx"))
	return c.Edit(response.String(), MainMenu())
}

func (r *HandlerRouter) parseCallback(callback string) (page, id, mode string) {
	if strings.Contains(callback, "_") {
		id = strings.Split(callback, "_")[1]
	}
	if strings.Contains(callback, "|") {
		splitted := strings.Split(callback, "|")
		page = strings.Split(callback, "|")[1]
		if len(splitted) > 1 {
				mode = strings.Split(callback, "|")[2]
		}
	}
	return page, id, mode
}
