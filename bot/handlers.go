package bot

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	constants "wishlist-bot/constant"
	sv "wishlist-bot/service"

	"gopkg.in/telebot.v4"
)

func setUpHandlers(bot *telebot.Bot, userService sv.UserService, wishlistService sv.WishService) {
	bot.Handle(constants.ON_START, func(c telebot.Context) error {
		exists := userService.ExistsById(c.Chat().ID)
		if exists != nil {
			userDto := sv.UserDto{ID: c.Chat().ID, Name: c.Chat().FirstName, Surname: c.Chat().LastName, Username: c.Chat().Username}
			saved := userService.Save(userDto)
			if saved == nil {
				return c.Send("Приветствуем, "+c.Chat().Username+". Этот бот был создан с целью внесения данных о работниках ЦЦР\n\n"+
					"Для дополнительной информации нажмите кнопку \"Помощь\", для внесения остальных данных нажмите \"Регистрация\"", menu)
			}
			return c.Send(fmt.Sprintf("Ошибка сохранения ваших данных. Напишите @qhsdkx %+v", saved))
		}
		return c.Send("Приветствуем, "+c.Chat().Username+". Вы нажали кнопку старта. Выберите действие", menu)
	})

	bot.Handle(constants.ON_HELP, func(c telebot.Context) error {
		return c.Send("Это бот для работников ЦЦР. Сюда вы можете внести данные о своих пожелниях на день рождения для своих коллег")
	}, checkSheluvssic())

	bot.Handle(telebot.OnCallback, func(c telebot.Context) error {
		var page, id string
		callback := c.Callback().Data[1:]
		if strings.Contains(callback, "_") {
			id = strings.Split(callback, "_")[1]
		}
		if strings.Contains(callback, "|") {
			page = strings.Split(callback, "|")[1]
		}
		switch callback {
		case constants.BTN_REGISTER:
			return onButtonRegister(c, userService)
		case constants.BTN_HELP:
			return onButtonHelp(c)
		case constants.BTN_WISHLIST:
			return onButtonWishlist(c, userService)
		case constants.BTN_ALL_USERS:
			return handleUserList(c, userService)
		case constants.BTN_PREV:
			return onButtonPrev(c)
		case constants.BTN_DELETE_ME:
			return onDeleteMe(c, userService)
		case constants.BTN_ME:
			return onButtonMyData(c, userService)
		case constants.BTN_EDIT_NAME:
			return onEditName(c)
		case constants.BTN_EDIT_SURNAME:
			return onEditSurname(c)
		case constants.BTN_EDIT_BIRTHDATE:
			return onEditBirthdate(c)
		case constants.BTN_EDIT_USERNAME:
			return onEditUserName(c)
		case constants.BTN_PREV_PAGE + "|" + page:
			return onButtonPrevAndBack(c, userService)
		case constants.BTN_NEXT_PAGE + "|" + page:
			return onButtonPrevAndBack(c, userService)
		case constants.USER_DATA_PREFIX + id:
			return onUserData(c, wishlistService, userService)
		case constants.BACK_TO_LIST:
			return handleUserList(c, userService)
		case constants.BTN_SHOW_ALL_WISHLIST:
			return onShowWishlist(c, wishlistService)
		case constants.BTN_REGISTER_WISHLIST:
			return onButtonRegWishList(c)
		case constants.DELETE_WISH:
			return onButtonDeleteWish(c)
		}
		return c.Respond()
	}, checkSheluvssic())

	bot.Handle(telebot.OnText, func(c telebot.Context) error {
		userState, exists := states[c.Chat().ID]
		if !exists {
			return c.Send("Пожалуйста, начните с команды /start")
		}

		switch userState {
		case constants.AWAITING_BIRTHDATE:
			return onAwaitingBirthdate(c, userService)
		case constants.AWAITING_NAME:
			return onAwaitingName(c, userService)
		case constants.AWAITING_SURNAME:
			return onAwaitingSurname(c, userService)

		case constants.AWAITING_NEW_NAME:
			return onAwaitingNewName(c, userService)
		case constants.AWAITING_NEW_SURNAME:
			return onAwaitingNewSurname(c, userService)
		case constants.AWAITING_NEW_BIRTHDATE:
			return onAwaitingNewBirthdate(c, userService)
		case constants.AWAITING_NEW_USERNAME:
			return onAwaitingNewUsername(c, userService)

		case constants.AWAITING_WISHES:
			return onAwaitingWishlist(c, wishlistService)
		case constants.DELETE_WISH:
			return onDeleteWish(c, wishlistService)
		default:
			return onError(c)
		}
	}, checkSheluvssic())

	bot.Handle(constants.SEND_MESSAGE_ADMIN+"_reg", func(c telebot.Context) error {
		id, _ := strconv.ParseInt(os.Getenv("ADMIN_ID"), 10, 64)
		if c.Chat().ID != id {
			return c.Send("Ты не админ, нельзя такое делать!!!")
		}
		if len(c.Text()) <= len(constants.SEND_MESSAGE_ADMIN+"_reg")+1 {
			return c.Send("Пустое не отправится")
		}
		message, found := strings.CutPrefix(c.Text(), constants.SEND_MESSAGE_ADMIN+"_reg ")
		if !found {
			return c.Send("Ошибка с сообщением")
		}
		total, err := userService.FindAllRegistered()
		if err != nil {
			return c.Send("Ошибка при извлечении юзеров")
		}
		for _, user := range total {
			if user.ID == id {
				continue
			}
			_, err = c.Bot().Send(telebot.ChatID(user.ID), message)
			if err != nil {
				return err
			}
		}
		return c.Send("Все гуд")
	})

	bot.Handle(constants.SEND_MESSAGE_ADMIN+"_unreg", func(c telebot.Context) error {
		id, _ := strconv.ParseInt(os.Getenv("ADMIN_ID"), 10, 64)
		if c.Chat().ID != id {
			return c.Send("Ты не админ, нельзя такое делать!!!")
		}
		if len(c.Text()) <= len(constants.SEND_MESSAGE_ADMIN+"_unreg")+1 {
			return c.Send("Пустое не отправится")
		}
		message, found := strings.CutPrefix(c.Text(), constants.SEND_MESSAGE_ADMIN+"_unreg ")
		if !found {
			return c.Send("Ошибка с сообщением")
		}
		total, err := userService.FindAllUnregistered()
		if err != nil {
			return c.Send("Ошибка при извлечении юзеров")
		}
		for _, user := range total {
			if user.ID == id {
				continue
			}
			_, err = c.Bot().Send(telebot.ChatID(user.ID), message)
			if err != nil {
				return err
			}
		}
		return c.Send("Все гуд")
	})

	bot.Handle(constants.SHOW_UNREGISTERED, func(c telebot.Context) error {
		id, _ := strconv.ParseInt(os.Getenv("ADMIN_ID"), 10, 64)
		if c.Chat().ID != id {
			return c.Send("Ты не админ, нельзя такое делать!!!")
		}
		users, err := userService.FindAllUnregistered()
		if err != nil {
			return c.Send("Ошибка при извлечении пользователей")
		}
		if users == nil {
			return c.Send("Незарегистрированных пользователей нет!")
		}
		var builder strings.Builder
		builder.WriteString(fmt.Sprintf("Юзеры, которые не прошли полную регистрацию: \n\n"))
		for _, user := range users {
			builder.WriteString(fmt.Sprintf("- (%s) %s %s %s\n", user.Username, user.Name, user.Surname, user.Birthdate.Format("02.01.2006")))
		}
		return c.Send(builder.String())
	})

}
