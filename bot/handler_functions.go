package bot

import (
	"fmt"
	"gopkg.in/telebot.v4"
	"strings"
	"time"
	constants "wishlist-bot/constant"
	sv "wishlist-bot/service"
)

var states = make(map[int64]string)

func onButtonMyData(c telebot.Context, service sv.UserService) error {
	user := service.FindById(c.Chat().ID)
	var response strings.Builder
	if user.Status == constants.REGISTERED {
		response.WriteString("*Ваши данные:*\n\n")
		response.WriteString(fmt.Sprintf("Ник в телеграме: %s\n%s %s\nДата рождения:%s \n\n", user.Username, user.Surname, user.Name, user.Birthdate.Format("02.01.2006")))
		response.WriteString("Кнопками ниже вы можете обновить данные")
		if _, err := bot.Edit(c.Message(), response.String(), wantEditSelector, telebot.ModeMarkdown); err != nil {
			return err
		}
		return nil

	}
	response.WriteString(fmt.Sprintf("Вы не прошли полную регистрацию, пока что в базе лишь ваши никнейм и имя, предоставленные телеграммом\n\nИмя: %s \nникнейм: %s", user.Name, user.Username))
	if _, err := bot.Edit(c.Message(), response.String(), menu, telebot.ModeMarkdown); err != nil {
		return err
	}
	return nil
}

func onEditName(c telebot.Context, service sv.UserService) error {
	states[c.Chat().ID] = constants.AWAITING_NEW_NAME
	return c.Send("Введите новое имя")
}

func onEditSurname(c telebot.Context, service sv.UserService) error {
	states[c.Chat().ID] = constants.AWAITING_NEW_SURNAME
	return c.Send("Введите новую фамилию")
}

func onEditBirthdate(c telebot.Context, service sv.UserService) error {
	states[c.Chat().ID] = constants.AWAITING_NEW_BIRTHDATE
	return c.Send("Введите новый день рождения")
}

func onEditUserName(c telebot.Context, service sv.UserService) error {
	states[c.Chat().ID] = constants.AWAITING_NEW_USERNAME
	return c.Send("Введите новый ник в телеграме (начинается с @). Проверьте его правильность, т.к. по нему можно перейти к вам в личные сообщения")
}

func onAwaitingNewName(c telebot.Context, service sv.UserService) error {
	delete(states, c.Chat().ID)
	if service.UpdateName(c.Text(), c.Chat().ID) {
		return c.Send("Имя успешно обновлено. Можете продолжить обновление", wantEditSelector)
	}
	return c.Send("Ошибка сохранения данных")
}

func onAwaitingNewSurname(c telebot.Context, service sv.UserService) error {
	delete(states, c.Chat().ID)
	if service.UpdateSurname(c.Text(), c.Chat().ID) {
		return c.Send("Фамилия успешно обновлена. Можете продолжить обновление", wantEditSelector)
	}
	return c.Send("Ошибка сохранения данных")
}

func onAwaitingNewBirthdate(c telebot.Context, service sv.UserService) error {
	delete(states, c.Chat().ID)
	date, err := parseDate(c.Text())
	if err != nil {
		return c.Send("Неверный формат даты. Пожалуйста, используйте ДД.ММ.ГГГГ.")
	}
	if service.UpdateBirthdate(&date, c.Chat().ID) {
		return c.Send("Дата успешно обновлена. Можете продолжить обновление", wantEditSelector)
	}
	return c.Send("Ошибка сохранения данных")
}

func onAwaitingNewUsername(c telebot.Context, service sv.UserService) error {
	delete(states, c.Chat().ID)
	if !strings.HasPrefix(c.Text(), "@") {
		return c.Send("Неверный формат. Никнейм начинается с \"@\". Попробуйте еще раз")
	}
	if service.UpdateUsername(c.Text(), c.Chat().ID) {
		return c.Send("Никнейм успешно обновлен. Можете продолжить обновление", wantEditSelector)
	}
	return c.Send("Ошибка сохранения данных")
}

func onButtonRegister(c telebot.Context, service sv.UserService) error {
	if !service.CheckIfRegistered(c.Chat().ID) {
		states[c.Chat().ID] = constants.AWAITING_BIRTHDATE
		if _, err := bot.Edit(c.Message(), "Пожалуйста, введите дату рождения в формате ДД.ММ.ГГГГ"); err != nil {
			return err
		}
	}
	if _, err := bot.Edit(c.Message(), "Вы уже зарегистрированный пользователь. Возвращаем в начало", menu); err != nil {
		return err
	}
	return nil
}

func onButtonHelp(c telebot.Context, service sv.UserService) error {
	if _, err := bot.Edit(c.Message(), "Нажмите \"Регистрация\", чтобы начать ввод данных.", menu); err != nil {
		return err
	}
	return nil
}

func onButtonWishlist(c telebot.Context, service sv.UserService) error {
	states[c.Chat().ID] = constants.AWAITING_WISHES
	if _, err := bot.Edit(c.Message(), "Введите свои пожелания через запятую (Майбах, бананы, вилла в Италии)", wishlistSelector, telebot.ModeMarkdown); err != nil {
		return err
	}
	return nil
}

func onButtonPrev(c telebot.Context, service sv.UserService) error {
	delete(states, c.Chat().ID)
	if _, err := bot.Edit(c.Message(), "Возвращаем вас в начало", menu); err != nil {
		return err
	}
	return nil
}

func onButtonAllUsers(c telebot.Context, service sv.UserService) error {
	users := service.FindAll()
	var response strings.Builder

	response.WriteString("*Список пользователей:*\n\n")
	for _, user := range users {
		response.WriteString(fmt.Sprintf("Ник в телеграме: %s\n%s %s\nДата рождения: %s\n\n", user.Username, user.Surname, user.Name, user.Birthdate.Format("02.01.2006")))
	}
	if _, err := bot.Edit(c.Message(), response.String(), menu, telebot.ModeMarkdown); err != nil {
		return err
	}
	return nil
}

func onAwaitingBirthdate(c telebot.Context, service sv.UserService) error {
	date, err := parseDate(c.Text())
	if err != nil {
		return c.Send("Неверный формат даты. Пожалуйста, используйте ДД.ММ.ГГГГ.")
	}
	if service.UpdateBirthdate(&date, c.Chat().ID) {
		states[c.Chat().ID] = constants.AWAITING_NAME
		return c.Send("Дата успешно сохранена. Далее введите желаемое в системе имя")
	}
	return c.Send("Ошибка сохранения данных")
}

func onAwaitingName(c telebot.Context, service sv.UserService) error {
	if service.UpdateName(c.Text(), c.Chat().ID) {
		states[c.Chat().ID] = constants.AWAITING_SURNAME
		return c.Send("Имя успешно сохранено. Далее введите желаемую в системе фамилию")
	}
	return c.Send("Ошибка сохранения данных")
}

func onAwaitingSurname(c telebot.Context, service sv.UserService) error {
	if service.UpdateSurname(c.Text(), c.Chat().ID) {
		service.UpdateStatus(constants.REGISTERED, c.Chat().ID)
		delete(states, c.Chat().ID)
		return c.Send("Фамилия успешно сохранена. Возвращаем в начальное меню", menu)
	}
	return c.Send("Ошибка сохранения данных")
}

func onAwaitingWishlist(c telebot.Context, service sv.UserService) error {
	delete(states, c.Chat().ID)
	return c.Send("Ваш список желний успешно сохранен", menu)
}

func onRestoreUser(c telebot.Context, service sv.UserService) error {
	service.Restore(c.Chat().ID)
	return c.Send("Вы успешно восстановлены в базе. Выбирайте дальнейшие действия", menu)
}

func onDeleteMe(c telebot.Context, service sv.UserService) error {
	service.Delete(c.Chat().ID)
	return c.Send("Вы были удалены из базы. Для доступных действий начните с команды /start")
}

func onError(c telebot.Context) error {
	delete(states, c.Chat().ID)
	return c.Send("Неизвестное состояние. Пожалуйста, начните заново с команды /start.")
}

func parseDate(date string) (time.Time, error) {
	parsedDate, err := time.Parse("02.01.2006", date)
	return parsedDate, err
}

func CheckDeleted(service sv.UserService) telebot.MiddlewareFunc {
	return func(next telebot.HandlerFunc) telebot.HandlerFunc {
		return func(c telebot.Context) error {
			if service.CheckIfDeleted(c.Chat().ID) {
				_, err := bot.Edit(c.Message(), "Вы удалены...", deletedSelector)
				return err
			}
			return next(c)
		}
	}
}
