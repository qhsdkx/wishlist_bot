package bot

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	constants "wishlist-bot/internal/constant"
	"wishlist-bot/internal/fsm"
	"wishlist-bot/internal/user"

	"gopkg.in/telebot.v4"
)

type UserHandler struct {
	service user.Service
	states  fsm.StateStore
}

func NewUserHandler(service user.Service, states fsm.StateStore) *UserHandler {
	return &UserHandler{service: service, states: states}
}

func (h *UserHandler) ShowProfile(c telebot.Context) error {
	user, err := h.service.FindByID(c.Chat().ID)
	if err != nil {
		return c.Edit(fmt.Sprintf("Невозможно найти юзера по ID %d", c.Chat().ID), MainMenu())
	}

	var msg strings.Builder
	if user.Status == "REGISTERED" {
		msg.WriteString(fmt.Sprintf("Ваши данные:\n\nНик: %s\n%s %s\nДата рождения: %s\n\n",
			user.Username, user.Surname, user.Name, user.Birthdate.Format("02.01.2006")))
		msg.WriteString("Кнопками ниже вы можете обновить данные")
		return c.Edit(msg.String(), EditMenu())
	}

	msg.WriteString(fmt.Sprintf("Вы не прошли полную регистрацию.\nИмя: %s\nНикнейм: %s", user.Name, user.Username))
	return c.Edit(msg.String(), MainMenu())
}

// func onButtonMyData(c telebot.Context, service sv.UserService) error {
// 	user, err := service.FindById(c.Chat().ID)
// 	if err != nil {
// 		return c.Edit(fmt.Sprintf("Невозиожно найти юзера по айди: %d", c.Chat().ID), menu)
// 	}
// 	var response strings.Builder
// 	if user.Status == constants.REGISTERED {
// 		response.WriteString("Ваши данные:\n\n")
// 		response.WriteString(fmt.Sprintf("Ник в телеграме: %s\n%s %s\nДата рождения: %s \n\n", user.Username, user.Surname, user.Name, user.Birthdate.Format("02.01.2006")))
// 		response.WriteString("Кнопками ниже вы можете обновить данные")
// 		if err = c.Edit(response.String(), wantEditSelector); err != nil {
// 			return c.Edit(fmt.Sprintf("Ошибка с редактированием собщения (%+v)", err), menu)
// 		}
// 		return c.Respond()

// 	}
// 	response.WriteString(fmt.Sprintf("Вы не прошли полную регистрацию, пока что в базе лишь ваши никнейм и имя, предоставленные телеграммом\n\nИмя: %s \nникнейм: %s", user.Name, user.Username))
// 	if _, err = c.Bot().Edit(c.Message(), response.String(), menu); err != nil {
// 		return c.Edit(fmt.Sprintf("Непредвиденная ошибка %v", err), menu)
// 	}
// 	return c.Respond()
// }

func (h *UserHandler) EditName(c telebot.Context) error {
	h.states.Set(c.Chat().ID, "AWAITING_NEW_NAME")
	return c.Edit("Введите новое имя")
}

// func onEditName(c telebot.Context) error {
// 	states[c.Chat().ID] = constants.AWAITING_NEW_NAME
// 	return c.Edit("Введите новое имя")
// }

func (h *UserHandler) AwaitingNewName(c telebot.Context) error {
	text := strings.TrimSpace(c.Text())
	if strings.Count(text, " ") > 0 {
		return c.Send("Введите только имя, без фамилии")
	}

	if err := h.service.UpdateName(text, c.Chat().ID); err != nil {
		return c.Send("Ошибка сохранения данных", MainMenu())
	}

	h.states.Delete(c.Chat().ID)
	return c.Send("Имя успешно обновлено", EditMenu())
}

// func onAwaitingNewName(c telebot.Context, service sv.UserService) error {
// 	if count := strings.Count(strings.TrimSpace(c.Text()), " "); count > 0 {
// 		return c.Send(fmt.Sprintf("Вероятно, вы ввели имя и фамилию сразу.\nВведите пожалуйста только имя"))
// 	}
// 	if err := service.UpdateName(c.Text(), c.Chat().ID); err == nil {
// 		delete(states, c.Chat().ID)
// 		return c.Send("Имя успешно обновлено. Можете продолжить обновление", wantEditSelector)
// 	}
// 	return c.Send("Ошибка сохранения данных", menu)
// }

func (h *UserHandler) EditSurname(c telebot.Context) error {
	h.states.Set(c.Chat().ID, constants.AWAITING_NEW_SURNAME)
	return c.Edit("Введите новую фамилию")
}

// func onEditSurname(c telebot.Context) error {
// 	states[c.Chat().ID] = constants.AWAITING_NEW_SURNAME
// 	return c.Edit("Введите новую фамилию")
// }

func (h *UserHandler) AwaitingNewSurname(c telebot.Context) error {
	text := strings.TrimSpace(c.Text())
	if strings.Count(text, " ") > 0 {
		return c.Send("Введите только фамилию, без имени")
	}

	if err := h.service.UpdateName(text, c.Chat().ID); err != nil {
		return c.Send("Ошибка сохранения данных", MainMenu())
	}

	h.states.Delete(c.Chat().ID)
	return c.Send("Фамилия успешно обновлена", EditMenu())
}

// func onAwaitingNewSurname(c telebot.Context, service sv.UserService) error {
// 	if count := strings.Count(strings.TrimSpace(c.Text()), " "); count > 0 {
// 		return c.Send(fmt.Sprintf("Вероятно, вы ввели два слова сразу.\nВведите пожалуйста только фамилию"))
// 	}
// 	if err := service.UpdateSurname(c.Text(), c.Chat().ID); err == nil {
// 		delete(states, c.Chat().ID)
// 		return c.Send("Фамилия успешно обновлена. Можете продолжить обновление", wantEditSelector)
// 	}
// 	return c.Send("Ошибка сохранения данных", menu)
// }

func (h *UserHandler) EditBirthdate(c telebot.Context) error {
	h.states.Set(c.Chat().ID, constants.AWAITING_NEW_BIRTHDATE)
	return c.Edit("Введите новый день рождения")
}

// func onEditBirthdate(c telebot.Context) error {
// 	states[c.Chat().ID] = constants.AWAITING_NEW_BIRTHDATE
// 	return c.Edit("Введите новый день рождения")
// }

func (h *UserHandler) AwaitingNewBirthdate(c telebot.Context) error {
	date, err := time.Parse("02.01.2006", c.Text())
	if err != nil {
		return c.Send("Неверный формат даты. Пожалуйста, используйте ДД.ММ.ГГГГ.")
	}
	if errUpdate := h.service.UpdateBirthdate(&date, c.Chat().ID); errUpdate == nil {
		h.states.Delete(c.Chat().ID)
		return c.Send("Дата успешно обновлена. Можете продолжить обновление", EditMenu())
	}
	return c.Send("Ошибка сохранения данных", MainMenu())
}

// func onAwaitingNewBirthdate(c telebot.Context, service sv.UserService) error {
// 	date, err := time.Parse("02.01.2006", c.Text())
// 	if err != nil {
// 		return c.Send("Неверный формат даты. Пожалуйста, используйте ДД.ММ.ГГГГ.")
// 	}
// 	if errUpdate := service.UpdateBirthdate(&date, c.Chat().ID); errUpdate == nil {
// 		delete(states, c.Chat().ID)
// 		return c.Send("Дата успешно обновлена. Можете продолжить обновление", wantEditSelector)
// 	}
// 	return c.Send("Ошибка сохранения данных", menu)
// }

func (h *UserHandler) EditUserName(c telebot.Context) error {
	h.states.Set(c.Chat().ID, constants.AWAITING_NEW_USERNAME)
	return c.Edit("Введите новый ник в телеграме (начинается с @). Проверьте его правильность, т.к. по нему можно перейти к вам в личные сообщения")
}

// func onEditUserName(c telebot.Context) error {
// 	states[c.Chat().ID] = constants.AWAITING_NEW_USERNAME
// 	return c.Edit("Введите новый ник в телеграме (начинается с @). Проверьте его правильность, т.к. по нему можно перейти к вам в личные сообщения")
// }

func (h *UserHandler) AwaitingNewUsername(c telebot.Context) error {
	if !strings.HasPrefix(c.Text(), "@") {
		return c.Send("Неверный формат. Никнейм начинается с \"@\". Попробуйте еще раз")
	}
	if count := strings.Count(strings.TrimSpace(c.Text()), " "); count > 0 {
		return c.Send(fmt.Sprintf("Вероятно, вы ввели два слова.\nВведите пожалуйста ваш ник одним словом"))
	}
	if err := h.service.UpdateUsername(c.Text(), c.Chat().ID); err == nil {
		h.states.Delete(c.Chat().ID)
		return c.Send("Никнейм успешно обновлен. Можете продолжить обновление", EditMenu())
	}
	return c.Send("Ошибка сохранения данных", MainMenu())
}

// func onAwaitingNewUsername(c telebot.Context, service sv.UserService) error {
// 	if !strings.HasPrefix(c.Text(), "@") {
// 		return c.Send("Неверный формат. Никнейм начинается с \"@\". Попробуйте еще раз")
// 	}
// 	if count := strings.Count(strings.TrimSpace(c.Text()), " "); count > 0 {
// 		return c.Send(fmt.Sprintf("Вероятно, вы ввели два слова.\nВведите пожалуйста ваш ник одним словом"))
// 	}
// 	if err := service.UpdateUsername(c.Text(), c.Chat().ID); err == nil {
// 		delete(states, c.Chat().ID)
// 		return c.Send("Никнейм успешно обновлен. Можете продолжить обновление", wantEditSelector)
// 	}
// 	return c.Send("Ошибка сохранения данных", menu)
// }

func (h *UserHandler) Register(c telebot.Context) error {
	if registered := h.service.CheckIfRegistered(c.Chat().ID); registered != nil {
		h.states.Set(c.Chat().ID, constants.AWAITING_BIRTHDATE)
		if err := c.Edit("Пожалуйста, введите дату рождения в формате ДД.ММ.ГГГГ"); err != nil {
			return err
		}
		return nil
	}
	if err := c.Edit("Вы уже зарегистрированный пользователь. Возвращаем в начало", MainMenu()); err != nil {
		return err
	}
	return nil
}

// func onButtonRegister(c telebot.Context, service sv.UserService) error {
// 	if registered := service.CheckIfRegistered(c.Chat().ID); registered != nil {
// 		states[c.Chat().ID] = constants.AWAITING_BIRTHDATE
// 		if err := c.Edit("Пожалуйста, введите дату рождения в формате ДД.ММ.ГГГГ"); err != nil {
// 			return err
// 		}
// 		return nil
// 	}
// 	if err := c.Edit("Вы уже зарегистрированный пользователь. Возвращаем в начало", menu); err != nil {
// 		return err
// 	}
// 	return nil
// }

func (h *UserHandler) Help(c telebot.Context) error {
	response := strings.Builder{}
	response.WriteString(fmt.Sprintf("Данная система была создана с целью помощи работникам ЦЦР (пока что 9-го департамента) следить за днями рождения коллег\n"))
	response.WriteString(fmt.Sprintf("ВНИМАНИЕ. ВСЕ предусмотренные уведомления приходят только в случае полной регистрации пользоваетеля по кнопке \"Регистрация\"\n"))
	response.WriteString(fmt.Sprintf("Короткая информация по возможностям бота:\n\n"))
	response.WriteString(fmt.Sprintf("• \"Редактировать мои данные\" - кнопка, представляющая возможность изменения введенных при регистрации данных. соответственно доступна только зарегистрированным\n"))
	response.WriteString(fmt.Sprintf("• \"Список желаний\" - дает возможность ввести ваши пожелания, удалить что-то либо посмотреть пожелания других\n"))
	response.WriteString(fmt.Sprintf("• \"Регистрация\" - необходима для регистрации вас в системе (ввод имени, фамилии и даты рождения)\n"))
	response.WriteString(fmt.Sprintf("• \"Показать всех пользователей\" - показывает всех ЗАРЕГИСТРИРОВАННЫХ пользователей и доступна только для них. По нажатию на кнопку с именем покажется день рождения человека и его пожелания\n"))
	response.WriteString(fmt.Sprintf("• \"Удалить меня в базе\" - полностью удаляет вас в базе. Далее необходимо следовать инструкции\n\n"))
	response.WriteString(fmt.Sprintf("Это было кратко описание основных возможностей бота. Так как обратной связи пока нет, то в случае возникающих проблем или предложений пишите разработчику @qhsdkx"))
	if _, err := c.Bot().Edit(c.Message(), response.String(), MainMenu()); err != nil {
		return err
	}
	return nil
}

// func onButtonHelp(c telebot.Context) error {
// 	response := strings.Builder{}
// 	response.WriteString(fmt.Sprintf("Данная система была создана с целью помощи работникам ЦЦР (пока что 9-го департамента) следить за днями рождения коллег\n"))
// 	response.WriteString(fmt.Sprintf("ВНИМАНИЕ. ВСЕ предусмотренные уведомления приходят только в случае полной регистрации пользоваетеля по кнопке \"Регистрация\"\n"))
// 	response.WriteString(fmt.Sprintf("Короткая информация по возможностям бота:\n\n"))
// 	response.WriteString(fmt.Sprintf("• \"Редактировать мои данные\" - кнопка, представляющая возможность изменения введенных при регистрации данных. соответственно доступна только зарегистрированным\n"))
// 	response.WriteString(fmt.Sprintf("• \"Список желаний\" - дает возможность ввести ваши пожелания, удалить что-то либо посмотреть пожелания других\n"))
// 	response.WriteString(fmt.Sprintf("• \"Регистрация\" - необходима для регистрации вас в системе (ввод имени, фамилии и даты рождения)\n"))
// 	response.WriteString(fmt.Sprintf("• \"Показать всех пользователей\" - показывает всех ЗАРЕГИСТРИРОВАННЫХ пользователей и доступна только для них. По нажатию на кнопку с именем покажется день рождения человека и его пожелания\n"))
// 	response.WriteString(fmt.Sprintf("• \"Удалить меня в базе\" - полностью удаляет вас в базе. Далее необходимо следовать инструкции\n\n"))
// 	response.WriteString(fmt.Sprintf("Это было кратко описание основных возможностей бота. Так как обратной связи пока нет, то в случае возникающих проблем или предложений пишите разработчику @qhsdkx"))
// 	if _, err := c.Bot().Edit(c.Message(), response.String(), menu); err != nil {
// 		return err
// 	}
// 	return nil
// }

func (h *UserHandler) Prev(c telebot.Context) error {
	h.states.Delete(c.Chat().ID)
	return c.Edit(c.Message(), "Возвращаем вас в начало", MainMenu())
}

// func onButtonPrev(c telebot.Context) error {
// 	delete(states, c.Chat().ID)
// 	if _, err := c.Bot().Edit(c.Message(), "Возвращаем вас в начало", menu); err != nil {
// 		return c.Edit("Непредвиденная ошибка. В начало", menu)
// 	}
// 	return nil
// }

func (h *UserHandler) AwaitingBirthdate(c telebot.Context) error {
	date, err := time.Parse("02.01.2006", c.Text())
	if err != nil {
		return c.Send("Неверный формат даты. Пожалуйста, используйте ДД.ММ.ГГГГ.")
	}
	if errUpdated := h.service.UpdateBirthdate(&date, c.Chat().ID); errUpdated == nil {
		h.states.Set(c.Chat().ID, constants.AWAITING_NAME)
		return c.Send("Дата успешно сохранена. Далее введите желаемое в системе имя")
	}
	return c.Send("Ошибка сохранения данных", MainMenu())
}

// func onAwaitingBirthdate(c telebot.Context, service sv.UserService) error {
// 	date, err := time.Parse("02.01.2006", c.Text())
// 	if err != nil {
// 		return c.Send("Неверный формат даты. Пожалуйста, используйте ДД.ММ.ГГГГ.")
// 	}
// 	if errUpdated := service.UpdateBirthdate(&date, c.Chat().ID); errUpdated == nil {
// 		states[c.Chat().ID] = constants.AWAITING_NAME
// 		return c.Send("Дата успешно сохранена. Далее введите желаемое в системе имя")
// 	}
// 	return c.Send("Ошибка сохранения данных", menu)
// }

func (h *UserHandler) AwaitingName(c telebot.Context) error {
	if count := strings.Count(strings.TrimSpace(c.Text()), " "); count > 0 {
		return c.Send(fmt.Print("Вероятно, вы ввели имя и фамилию сразу.\nВведите пожалуйста только имя"))
	}
	if err := h.service.UpdateName(c.Text(), c.Chat().ID); err == nil {
		h.states.Set(c.Chat().ID, constants.AWAITING_SURNAME)
		return c.Send("Имя успешно сохранено. Далее введите желаемую в системе фамилию")
	}
	return c.Send("Ошибка сохранения данных", MainMenu())
}

// func onAwaitingName(c telebot.Context, service sv.UserService) error {
// 	if count := strings.Count(strings.TrimSpace(c.Text()), " "); count > 0 {
// 		return c.Send(fmt.Print("Вероятно, вы ввели имя и фамилию сразу.\nВведите пожалуйста только имя"))
// 	}
// 	if err := service.UpdateName(c.Text(), c.Chat().ID); err == nil {
// 		states[c.Chat().ID] = constants.AWAITING_SURNAME
// 		return c.Send("Имя успешно сохранено. Далее введите желаемую в системе фамилию")
// 	}
// 	return c.Send("Ошибка сохранения данных", menu)
// }

func (h *UserHandler) AwaitingSurname(c telebot.Context) error {
	if count := strings.Count(strings.TrimSpace(c.Text()), " "); count > 0 {
		return c.Send(fmt.Print("Вероятно, вы ввели два слова.\nВведите пожалуйста только фамилию"))
	}
	if err := h.service.UpdateSurname(c.Text(), c.Chat().ID); err == nil {
		errUpdate := h.service.UpdateStatus(constants.REGISTERED, c.Chat().ID)
		if errUpdate != nil {
			return c.Send(fmt.Sprintf("Ошибка изменения статуса у юзера с айди %d", c.Chat().ID), MainMenu())
		}
		h.states.Delete(c.Chat().ID)
		return c.Send("Фамилия успешно сохранена. Возвращаем в начальное меню", MainMenu())
	}
	return c.Send("Ошибка сохранения данных", MainMenu())
}

// func onAwaitingSurname(c telebot.Context, service sv.UserService) error {
// 	if count := strings.Count(strings.TrimSpace(c.Text()), " "); count > 0 {
// 		return c.Send(fmt.Print("Вероятно, вы ввели два слова.\nВведите пожалуйста только фамилию"))
// 	}
// 	if err := service.UpdateSurname(c.Text(), c.Chat().ID); err == nil {
// 		errUpdate := service.UpdateStatus(constants.REGISTERED, c.Chat().ID)
// 		if errUpdate != nil {
// 			return c.Send(fmt.Sprintf("Ошибка изменения статуса у юзера с айди %d", c.Chat().ID), menu)
// 		}
// 		delete(states, c.Chat().ID)
// 		return c.Send("Фамилия успешно сохранена. Возвращаем в начальное меню", menu)
// 	}
// 	return c.Send("Ошибка сохранения данных", menu)
// }

func (h *UserHandler) DeleteMe(c telebot.Context) error {
	if err := h.service.Delete(c.Chat().ID); err != nil {
		err = c.Edit(fmt.Sprintf("Ошибка при удалении у юзера с айди %d", c.Chat().ID), MainMenu())
		return err
	}
	return c.Edit("Вы и ваши пожелания были успешно удалены из базы. Для того, чтобы начать с самого начала введите /start")
}

// func onDeleteMe(c telebot.Context, service sv.UserService) error {
// 	if err := service.Delete(c.Chat().ID); err != nil {
// 		err = c.Edit(fmt.Sprintf("Ошибка при удалении у юзера с айди %d", c.Chat().ID), menu)
// 		return err
// 	}
// 	return c.Edit("Вы и ваши пожелания были успешно удалены из базы. Для того, чтобы начать с самого начала введите /start")
// }

func (h *UserHandler) UserList(c telebot.Context, mode string) error {
	err := h.service.CheckIfRegistered(c.Chat().ID)
	if err != nil {
		err = c.Edit("Вы еще не зарегистрировались, чтобы просматривать пользователей. Пожалуйста, пройдите регистрацию", MainMenu())
		return err
	}
	users, pagination, err := h.service.FindAll(1, constants.USERS_PER_PAGE, mode)
	if err != nil {
		return c.Edit("Ошибка получения данных", MainMenu())
	}

	markup := h.createUserListMarkup(users, pagination, mode)
	if mode == constants.SHOW_USERS {
		return c.Edit("Список пользователей:", markup)
	}
	return c.Send("Список пользователей:", markup)
}

// func handleUserList(c telebot.Context, userService sv.UserService, mode string) error {
// 	err := userService.CheckIfRegistered(c.Chat().ID)
// 	if err != nil {
// 		err = c.Edit("Вы еще не зарегистрировались, чтобы просматривать пользователей. Пожалуйста, пройдите регистрацию", menu)
// 		return err
// 	}
// 	users, pagination, err := userService.FindAll(1, constants.USERS_PER_PAGE, mode)
// 	if err != nil {
// 		return c.Edit("Ошибка получения данных", menu)
// 	}

// 	markup := createUserListMarkup(users, pagination, mode)
// 	if mode == constants.SHOW_USERS {
// 		return c.Edit("Список пользователей:", markup)
// 	}
// 	return c.Send("Список пользователей:", markup)
// }

func (h *UserHandler) PrevAndBack(c telebot.Context, mode string) error {
	pageStr := strings.Split(c.Callback().Data, "|")[1]
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		return err
	}
	return h.updateUserListPage(c, page, mode)
}

// func onButtonPrevAndBack(c telebot.Context, userService sv.UserService, mode string) error {
// 	pageStr := strings.Split(c.Callback().Data, "|")[1]
// 	page, err := strconv.Atoi(pageStr)
// 	if err != nil {
// 		return err
// 	}
// 	return updateUserListPage(c, page, userService, mode)
// }

// func (h *UserHandler) UserData(c telebot.Context) error {
// 	data := c.Callback().Data[1:]
// 	if strings.HasPrefix(data, constants.USER_DATA_PREFIX) {
// 		userId, _ := strconv.ParseInt(data[len(constants.USER_DATA_PREFIX):], 10, 64)
// 		return h.showUserDetails(c, userId)
// 	}
// 	return c.Respond()
// }

// func onUserData(c telebot.Context, wishlistService sv.WishService, userService sv.UserService) error {
// 	data := c.Callback().Data[1:]
// 	if strings.HasPrefix(data, constants.USER_DATA_PREFIX) {
// 		userId, _ := strconv.ParseInt(data[len(constants.USER_DATA_PREFIX):], 10, 64)
// 		return showUserDetails(c, userId, wishlistService, userService)
// 	}
// 	return c.Respond()
// }

func (h *UserHandler) ChooseUser(c telebot.Context, id string) error {
	h.states.Set(c.Chat().ID, constants.SEND_MESSAGE_ADMIN+"_"+id)
	return c.Edit("Введите сообщение для данного пользователя")
}

// func onChooseUser(c telebot.Context, id string) error {
// 	states[c.Chat().ID] = constants.SEND_MESSAGE_ADMIN + "_" + id
// 	return c.Edit("Введите сообщение для данного пользователя")
// }

func (h *UserHandler) SendMessage(c telebot.Context) error {
	state, err := h.states.Get(c.Chat().ID)
	if err != nil {
		c.Send("Error with sending!")
	}
	userId, _ := strconv.ParseInt(state[:len(constants.SEND_MESSAGE_ADMIN+"_")], 10, 64)
	h.states.Delete(c.Chat().ID)
	err = c.Send(telebot.ChatID(userId), c.Text())
	if err != nil {
		return c.Send("Ошибка ", err)
	}
	return c.Send("Sent successfully!")
}

// func onSendMessage(c telebot.Context) error {
// 	userId, _ := strconv.ParseInt(states[c.Chat().ID][len(constants.SEND_MESSAGE_ADMIN+"_"):], 10, 64)
// 	delete(states, c.Chat().ID)
// 	_, err := c.Bot().Send(telebot.ChatID(userId), c.Text())
// 	if err != nil {
// 		return c.Send("Ошибка ", err)
// 	}
// 	return c.Send("Sent successfully!")
// }

func (h *UserHandler) createUserListMarkup(users []user.User, pagination *user.Pagination, mode string) *telebot.ReplyMarkup {
	markup := &telebot.ReplyMarkup{}
	rows := make([]telebot.Row, 0, len(users)+3)

	if mode == constants.SHOW_USERS {
		for _, user := range users {
			btn := markup.Data(
				fmt.Sprintf("%s %s (%s)", user.Name, user.Surname, user.Birthdate.Format("02.01.2006")),
				constants.USER_DATA_PREFIX+strconv.FormatInt(user.ID, 10),
			)
			rows = append(rows, markup.Row(btn))
		}
	}
	if mode == constants.SEND_MESSAGE_ADMIN {
		for _, user := range users {
			btn := markup.Data(
				fmt.Sprintf("%s %s (%s)", user.Name, user.Surname, user.Birthdate.Format("02.01.2006")),
				constants.SEND_MESSAGE_ADMIN+"_"+strconv.FormatInt(user.ID, 10),
			)
			rows = append(rows, markup.Row(btn))
		}
	}

	if pagination.TotalPages > 1 {
		var paginationRow []telebot.Btn
		if pagination.CurrentPage > 1 {
			prevBtn := markup.Data("⬅", constants.BTN_PREV_PAGE, strconv.Itoa(pagination.CurrentPage-1), mode)
			paginationRow = append(paginationRow, prevBtn)
		}

		if pagination.CurrentPage < pagination.TotalPages {
			nextBtn := markup.Data("➡", constants.BTN_NEXT_PAGE, strconv.Itoa(pagination.CurrentPage+1), mode)
			paginationRow = append(paginationRow, nextBtn)
		}

		rows = append(rows, markup.Row(paginationRow...))
	}
	rows = append(rows, markup.Row(markup.Data("В начало", constants.BTN_PREV)))

	markup.Inline(rows...)
	return markup
}
