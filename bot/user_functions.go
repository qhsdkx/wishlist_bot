package bot

import (
	"fmt"
	"gopkg.in/telebot.v4"
	"strconv"
	"strings"
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

func onEditName(c telebot.Context) error {
	states[c.Chat().ID] = constants.AWAITING_NEW_NAME
	return c.Send("Введите новое имя")
}

func onEditSurname(c telebot.Context) error {
	states[c.Chat().ID] = constants.AWAITING_NEW_SURNAME
	return c.Send("Введите новую фамилию")
}

func onEditBirthdate(c telebot.Context) error {
	states[c.Chat().ID] = constants.AWAITING_NEW_BIRTHDATE
	return c.Send("Введите новый день рождения")
}

func onEditUserName(c telebot.Context) error {
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
	if registered := !service.CheckIfRegistered(c.Chat().ID); registered {
		states[c.Chat().ID] = constants.AWAITING_BIRTHDATE
		if _, err := bot.Edit(c.Message(), "Пожалуйста, введите дату рождения в формате ДД.ММ.ГГГГ"); err != nil {
			return err
		}
		return nil
	}
	if _, err := bot.Edit(c.Message(), "Вы уже зарегистрированный пользователь. Возвращаем в начало", menu); err != nil {
		return err
	}
	return nil
}

func onButtonHelp(c telebot.Context) error {
	if _, err := bot.Edit(c.Message(), "Нажмите \"Регистрация\", чтобы начать ввод данных.", menu); err != nil {
		return err
	}
	return nil
}

func onButtonPrev(c telebot.Context) error {
	delete(states, c.Chat().ID)
	if _, err := bot.Edit(c.Message(), "Возвращаем вас в начало", menu); err != nil {
		return c.Send("Непредвиденная ошибка. В наччало", menu)
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

func onRestoreUser(c telebot.Context, service sv.UserService) error {
	service.Restore(c.Chat().ID)
	return c.Send("Вы успешно восстановлены в базе. Выбирайте дальнейшие действия", menu)
}

func onDeleteMe(c telebot.Context, service sv.UserService) error {
	service.Delete(c.Chat().ID)
	return c.Send("Вы были удалены из базы. Для доступных действий начните с команды /start")
}

func handleUserList(c telebot.Context, userService sv.UserService) error {
	users, pagination, err := userService.FindAll(1, constants.USERS_PER_PAGE)
	if err != nil {
		return c.Send("Ошибка получения данных")
	}

	markup := createUserListMarkup(users, pagination)
	return c.Edit("Список пользователей:", markup)
}

func onButtonPrevAndBack(c telebot.Context, userService sv.UserService) error {
	pageStr := strings.Split(c.Callback().Data, "|")[1]
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		return err
	}
	return updateUserListPage(c, page, userService)
}

func onUserData(c telebot.Context, wishlistService sv.WishService) error {
	data := c.Callback().Data[1:]
	if strings.HasPrefix(data, constants.USER_DATA_PREFIX) {
		userId, _ := strconv.ParseInt(data[len(constants.USER_DATA_PREFIX):], 10, 64)
		return showUserDetails(c, userId, wishlistService)
	}
	return c.Respond()
}

func createUserListMarkup(users []sv.UserDto, pagination *sv.Pagination) *telebot.ReplyMarkup {
	markup := &telebot.ReplyMarkup{}
	rows := make([]telebot.Row, 0, len(users)+3)

	for _, user := range users {
		btn := markup.Data(
			fmt.Sprintf("%s %s", user.Name, user.Surname),
			constants.USER_DATA_PREFIX+strconv.FormatInt(user.ID, 10),
		)
		rows = append(rows, markup.Row(btn))
	}

	if pagination.TotalPages > 1 {
		var paginationRow []telebot.Btn
		if pagination.CurrentPage > 1 {
			prevBtn := markup.Data("⬅", constants.BTN_PREV_PAGE, strconv.Itoa(pagination.CurrentPage-1))
			paginationRow = append(paginationRow, prevBtn)
		}

		if pagination.CurrentPage < pagination.TotalPages {
			nextBtn := markup.Data("➡", constants.BTN_NEXT_PAGE, strconv.Itoa(pagination.CurrentPage+1))
			paginationRow = append(paginationRow, nextBtn)
		}

		rows = append(rows, markup.Row(paginationRow...))
	}
	rows = append(rows, markup.Row(markup.Data("В начало", constants.BTN_PREV)))

	markup.Inline(rows...)
	return markup
}

func showUserDetails(c telebot.Context, userId int64, wishService sv.WishService) error {
	wishes := wishService.FindAllByUserId(userId)

	var msg strings.Builder
	msg.WriteString("🎁 Список желаний:\n\n")
	for _, wish := range wishes {
		msg.WriteString(fmt.Sprintf("• %s\n", wish.WishText))
	}

	_, err := bot.Edit(c.Message(), msg.String(), createBackButton())
	if err != nil {
		return c.Respond(&telebot.CallbackResponse{
			Text: "Ошибка отображения данных",
		})
	}

	return c.Respond()
}
