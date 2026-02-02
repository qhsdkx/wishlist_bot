package bot

import (
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"
	constants "wishlist-bot/internal/constant"
	"wishlist-bot/internal/fsm"
	"wishlist-bot/internal/logger/sl"
	"wishlist-bot/internal/user"

	"gopkg.in/telebot.v4"
)

type UserHandler struct {
	service user.Service
	states  fsm.StateStore
	log     *slog.Logger
}

func NewUserHandler(service user.Service, states fsm.StateStore, log *slog.Logger) *UserHandler {
	return &UserHandler{service: service, states: states, log: log}
}

func (h *UserHandler) ShowProfile(c telebot.Context) error {
	const op = "UserHandler.ShowProfile"
	u, err := h.service.FindByID(c.Chat().ID)
	if err != nil {
		h.log.Info(op, err, "user", u)
		return c.Edit(fmt.Sprintf("Невозможно найти юзера по ID %d", c.Chat().ID), MainMenu())
	}

	var msg strings.Builder
	if u.Status == "REGISTERED" {
		msg.WriteString(fmt.Sprintf("Ваши данные:\n\nНик: %s\n%s %s\nДата рождения: %s\n\n",
			u.Username, u.Surname, u.Name, u.Birthdate.Format("02.01.2006")))
		msg.WriteString("Кнопками ниже вы можете обновить данные")
		return c.Send(msg.String(), EditMenu())
	}

	msg.WriteString(fmt.Sprintf("Вы не прошли полную регистрацию.\nИмя: %s\nНикнейм: %s", u.Name, u.Username))
	return c.Edit(msg.String(), MainMenu())
}

func (h *UserHandler) EditName(c telebot.Context) error {
	h.states.Set(c.Chat().ID, constants.AWAITING_NEW_NAME)
	return c.Edit("Введите новое имя")
}

func (h *UserHandler) AwaitingNewName(c telebot.Context) error {
	const op = "UserHandler.AwaitingNewName"

	text := strings.TrimSpace(c.Text())
	if strings.Count(text, " ") > 0 {
		return c.Send("Введите только имя, без фамилии")
	}

	if err := h.service.UpdateName(text, c.Chat().ID); err != nil {
		h.log.Info(op, err)
		return c.Send("Ошибка сохранения данных", MainMenu())
	}

	h.states.Delete(c.Chat().ID)
	u, errFind := h.service.FindByID(c.Chat().ID)
	var msg strings.Builder
	if errFind != nil {
		h.log.Info(op, errFind)
		return c.Send("Ошибка. В начало.", MainMenu())
	}
	if u.Status == "REGISTERED" {
		msg.WriteString("ДАННЫЕ БЫЛИ ОБНОВЛЕНЫ\n\n")
		msg.WriteString(fmt.Sprintf("Ваши данные:\n\nНик: %s\n%s %s\nДата рождения: %s\n\n",
			u.Username, u.Surname, u.Name, u.Birthdate.Format("02.01.2006")))
		msg.WriteString("Кнопками ниже вы можете обновить данные")
		return c.Send(msg.String(), EditMenu())
	}
	return c.Send("Ошибка. В начало.", MainMenu())
}

func (h *UserHandler) EditSurname(c telebot.Context) error {
	h.states.Set(c.Chat().ID, constants.AWAITING_NEW_SURNAME)
	return c.Edit("Введите новую фамилию")
}

func (h *UserHandler) AwaitingNewSurname(c telebot.Context) error {
	const op = "UserHandler.AwaitingNewSurname"

	text := strings.TrimSpace(c.Text())
	if strings.Count(text, " ") > 0 {
		return c.Send("Введите только фамилию, без имени")
	}

	if err := h.service.UpdateSurname(text, c.Chat().ID); err != nil {
		h.log.Info(op, err)
		return c.Send("Ошибка сохранения данных", MainMenu())
	}

	h.states.Delete(c.Chat().ID)
	u, errFind := h.service.FindByID(c.Chat().ID)
	var msg strings.Builder
	if errFind != nil {
		h.log.Info(op, errFind)
		return c.Send("Ошибка. В начало.", MainMenu())
	}
	if u.Status == "REGISTERED" {
		msg.WriteString("ДАННЫЕ БЫЛИ ОБНОВЛЕНЫ\n\n")
		msg.WriteString(fmt.Sprintf("Ваши данные:\n\nНик: %s\n%s %s\nДата рождения: %s\n\n",
			u.Username, u.Surname, u.Name, u.Birthdate.Format("02.01.2006")))
		msg.WriteString("Кнопками ниже вы можете обновить данные")
		return c.Send(msg.String(), EditMenu())
	}
	return c.Send("Ошибка. В начало.", MainMenu())
}

func (h *UserHandler) EditBirthdate(c telebot.Context) error {
	h.states.Set(c.Chat().ID, constants.AWAITING_NEW_BIRTHDATE)
	return c.Edit("Введите новый день рождения")
}

func (h *UserHandler) AwaitingNewBirthdate(c telebot.Context) error {
	const op = "UserHandler.AwaitingNewBirthdate"

	date, err := time.Parse("02.01.2006", c.Text())
	if err != nil {
		h.log.Info(op, err)
		return c.Send("Неверный формат даты. Пожалуйста, используйте ДД.ММ.ГГГГ.")
	}
	if errUpdate := h.service.UpdateBirthdate(&date, c.Chat().ID); errUpdate == nil {
		h.log.Info(op, errUpdate)
		h.states.Delete(c.Chat().ID)
		u, errFind := h.service.FindByID(c.Chat().ID)
		var msg strings.Builder
		if errFind != nil {
			h.log.Info(op, errFind)
			return c.Send("Ошибка. В начало.", MainMenu())
		}
		if u.Status == "REGISTERED" {
			msg.WriteString("ДАННЫЕ БЫЛИ ОБНОВЛЕНЫ\n\n")
			msg.WriteString(fmt.Sprintf("Ваши данные:\n\nНик: %s\n%s %s\nДата рождения: %s\n\n",
				u.Username, u.Surname, u.Name, u.Birthdate.Format("02.01.2006")))
			msg.WriteString("Кнопками ниже вы можете обновить данные")
			return c.Send(msg.String(), EditMenu())
		}
	}
	return c.Send("Ошибка сохранения данных", MainMenu())
}

func (h *UserHandler) EditUserName(c telebot.Context) error {
	h.states.Set(c.Chat().ID, constants.AWAITING_NEW_USERNAME)
	return c.Edit("Введите новый ник в телеграме (начинается с @). Проверьте его правильность, т.к. по нему можно перейти к вам в личные сообщения")
}

func (h *UserHandler) AwaitingNewUsername(c telebot.Context) error {
	const op = "UserHandler.AwaitingNewUsername"

	if !strings.HasPrefix(c.Text(), "@") {
		return c.Send("Неверный формат. Никнейм начинается с \"@\". Попробуйте еще раз")
	}
	if count := strings.Count(strings.TrimSpace(c.Text()), " "); count > 0 {
		return c.Send(fmt.Sprintf("Вероятно, вы ввели два слова.\nВведите пожалуйста ваш ник одним словом"))
	}
	if err := h.service.UpdateUsername(c.Text(), c.Chat().ID); err == nil {
		h.log.Info(op, err)
		h.states.Delete(c.Chat().ID)
		u, errFind := h.service.FindByID(c.Chat().ID)
		var msg strings.Builder
		if errFind != nil {
			h.log.Info(op, err)
			return c.Send("Ошибка. В начало.", MainMenu())
		}
		if u.Status == "REGISTERED" {
			msg.WriteString("ДАННЫЕ БЫЛИ ОБНОВЛЕНЫ\n\n")
			msg.WriteString(fmt.Sprintf("Ваши данные:\n\nНик: %s\n%s %s\nДата рождения: %s\n\n",
				u.Username, u.Surname, u.Name, u.Birthdate.Format("02.01.2006")))
			msg.WriteString("Кнопками ниже вы можете обновить данные")
			return c.Send(msg.String(), EditMenu())
		}
	}
	return c.Send("Ошибка сохранения данных", MainMenu())
}

func (h *UserHandler) Register(c telebot.Context) error {
	const op = "UserHandler.Register"

	if registered := h.service.CheckIfRegistered(c.Chat().ID); registered != nil {
		h.service.Save(user.User{
			ID:       c.Chat().ID,
			Name:     c.Chat().FirstName,
			Surname:  c.Chat().LastName,
			Username: "@" + c.Chat().Username,
			Status:   constants.ADDED,
		})
		h.states.Set(c.Chat().ID, constants.AWAITING_BIRTHDATE)
		if err := c.Edit("Пожалуйста, введите дату рождения в формате ДД.ММ.ГГГГ"); err != nil {
			h.log.Error(op, sl.Err(err))
			return err
		}
		return nil
	}
	if err := c.Edit("Вы уже зарегистрированный пользователь. Возвращаем в начало", MainMenu()); err != nil {
		h.log.Error(op, sl.Err(err))
		return err
	}
	return nil
}

func (h *UserHandler) Prev(c telebot.Context) error {
	h.states.Delete(c.Chat().ID)
	return c.Edit("Возвращаем вас в начало", MainMenu())
}

func (h *UserHandler) AwaitingBirthdate(c telebot.Context) error {
	const op = "UserHandler.AwaitingBirthdate"

	date, err := time.Parse("02.01.2006", c.Text())
	if err != nil {
		h.log.Info(op, err)
		return c.Send("Неверный формат даты. Пожалуйста, используйте ДД.ММ.ГГГГ.")
	}
	if errUpdated := h.service.UpdateBirthdate(&date, c.Chat().ID); errUpdated == nil {
		h.states.Set(c.Chat().ID, constants.AWAITING_NAME)
		return c.Send("Дата успешно сохранена. Далее введите желаемое в системе имя")
	}
	return c.Send("Ошибка сохранения данных", MainMenu())
}

func (h *UserHandler) AwaitingName(c telebot.Context) error {
	if count := strings.Count(strings.TrimSpace(c.Text()), " "); count > 0 {
		return c.Send("Вероятно, вы ввели имя и фамилию сразу. Введите пожалуйста только имя")
	}
	if err := h.service.UpdateName(c.Text(), c.Chat().ID); err == nil {
		h.states.Set(c.Chat().ID, constants.AWAITING_SURNAME)
		return c.Send("Имя успешно сохранено.\nДалее введите желаемую в системе фамилию")
	}
	return c.Send("Ошибка сохранения данных", MainMenu())
}

func (h *UserHandler) AwaitingSurname(c telebot.Context) error {
	const op = "UserHandler.AwaitingSurname"

	if count := strings.Count(strings.TrimSpace(c.Text()), " "); count > 0 {
		return c.Send("Вероятно, вы ввели два слова.\nВведите пожалуйста только фамилию")
	}
	if err := h.service.UpdateSurname(c.Text(), c.Chat().ID); err == nil {
		errUpdate := h.service.UpdateStatus(constants.REGISTERED, c.Chat().ID)
		if errUpdate != nil {
			h.log.Error(op, sl.Err(err))
			return c.Send(fmt.Sprintf("Ошибка изменения статуса у юзера с айди %d", c.Chat().ID), MainMenu())
		}
		h.states.Delete(c.Chat().ID)
		return c.Send("Фамилия успешно сохранена. Возвращаем в начальное меню", MainMenu())
	}
	return c.Send("Ошибка сохранения данных", MainMenu())
}

func (h *UserHandler) DeleteMe(c telebot.Context) error {
	const op = "UserHandler.DeleteMe"
	if err := h.service.Delete(c.Chat().ID); err != nil {
		h.log.Error(op, sl.Err(err))
		err = c.Edit(fmt.Sprintf("Ошибка при удалении у юзера с айди %d", c.Chat().ID), MainMenu())
		return err
	}
	return c.Edit("Вы и ваши пожелания были успешно удалены из базы. Для того, чтобы начать с самого начала введите /start")
}

func (h *UserHandler) UserList(c telebot.Context, mode string) error {
	const op = "UserHandler.UserList"

	err := h.service.CheckIfRegistered(c.Chat().ID)
	if err != nil {
		h.log.Error(op, sl.Err(err))
		err = c.Edit("Вы еще не зарегистрировались, чтобы просматривать пользователей. Пожалуйста, пройдите регистрацию", MainMenu())
		return err
	}
	users, pagination, err := h.service.FindAll(1, constants.USERS_PER_PAGE, mode)
	if err != nil {
		h.log.Error(op, sl.Err(err))
		return c.Edit("Ошибка получения данных", MainMenu())
	}

	markup := h.createUserListMarkup(users, pagination, mode)
	if mode == constants.SHOW_USERS {
		return c.Edit("Список пользователей:", markup)
	}
	return c.Send("Список пользователей:", markup)
}

func (h *UserHandler) PrevAndBack(c telebot.Context, mode, pageStr string) error {
	const op = "UserHandler.PrevAndBack"

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		h.log.Error(op, sl.Err(err))
		return err
	}
	return h.updateUserListPage(c, page, mode)
}

func (h *UserHandler) ChooseUser(c telebot.Context, id string) error {
	h.states.Set(c.Chat().ID, constants.SEND_MESSAGE_ADMIN+"_"+id)
	return c.Edit("Введите сообщение для данного пользователя")
}

func (h *UserHandler) SendMessage(c telebot.Context) error {
	const op = "UserHandler.SendMessage"

	state, err := h.states.Get(c.Chat().ID)
	if err != nil {
		h.log.Error(op, sl.Err(err))
		c.Send("Error with sending!")
	}
	userId, _ := strconv.ParseInt(state[:len(constants.SEND_MESSAGE_ADMIN+"_")], 10, 64)
	h.states.Delete(c.Chat().ID)
	err = c.Send(telebot.ChatID(userId), c.Text())
	if err != nil {
		h.log.Error(op, sl.Err(err))
		return c.Send("Ошибка ", err)
	}
	return c.Send("Sent successfully!")
}

func (h *UserHandler) createUserListMarkup(users []user.User, pagination *user.Pagination, mode string) *telebot.ReplyMarkup {
	markup := &telebot.ReplyMarkup{}
	rows := make([]telebot.Row, 0, len(users)+3)

	if mode == constants.SHOW_USERS {
		for _, user := range users {
			btn := markup.Data(
				fmt.Sprintf("%s %s (%s)", user.Name, user.Surname, user.Birthdate.Format("02.01.2006")),
				NewCallbackData(
					constants.USER_DATA_PREFIX,
					mode,
					strconv.FormatInt(user.ID, 10),
					strconv.Itoa(pagination.CurrentPage)).string(),
			)
			rows = append(rows, markup.Row(btn))
		}
	}
	if mode == constants.SEND_MESSAGE_ADMIN {
		for _, user := range users {
			btn := markup.Data(
				fmt.Sprintf("%s %s (%s)", user.Name, user.Surname, user.Birthdate.Format("02.01.2006")),
				NewCallbackData(
					constants.SEND_MESSAGE_ADMIN,
					mode,
					strconv.FormatInt(user.ID, 10),
					strconv.Itoa(pagination.CurrentPage)).string(),
			)
			rows = append(rows, markup.Row(btn))
		}
	}

	if pagination.TotalPages > 1 {
		var paginationRow []telebot.Btn
		if pagination.CurrentPage > 1 {
			prevBtn := markup.Data("⬅",
				NewCallbackData(
					constants.BTN_PREV_PAGE,
					mode,
					"",
					strconv.Itoa(pagination.CurrentPage-1)).string(),
			)
			paginationRow = append(paginationRow, prevBtn)
		}

		if pagination.CurrentPage < pagination.TotalPages {
			nextBtn := markup.Data("➡",
				NewCallbackData(
					constants.BTN_PREV_PAGE,
					mode,
					"",
					strconv.Itoa(pagination.CurrentPage+1)).string(),
			)
			paginationRow = append(paginationRow, nextBtn)
		}

		rows = append(rows, markup.Row(paginationRow...))
	}
	rows = append(rows, markup.Row(markup.Data("В начало",
		NewCallbackData(
			constants.BTN_PREV,
			"",
			"",
			"").string(),
	),
	))

	markup.Inline(rows...)
	return markup
}
