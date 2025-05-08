package bot

import (
	"gopkg.in/telebot.v4"
	"strconv"
	"strings"
	constants "wishlist-bot/constant"
	sv "wishlist-bot/service"
)

func setUpHandlers(bot *telebot.Bot, userService sv.UserService, wishlistService sv.WishService) {
	bot.Handle(constants.ON_START, func(c telebot.Context) error {
		exists := userService.ExistsById(c.Chat().ID)
		deleted := userService.CheckIfDeleted(c.Chat().ID)
		if !exists && !deleted {
			userDto := sv.UserDto{ID: c.Chat().ID, Name: c.Chat().FirstName, Surname: c.Chat().LastName, Username: c.Chat().Username}
			if userService.Save(userDto) {
				return c.Send("Приветствуем, "+c.Chat().Username+". Этот бот был создан с целью внесения данных о работниках ЦЦР (даты рождения и пожелания)\n"+
					"Для дополнительной информации нажмите кнопку \"Помощь\", для внесения остальных данных (на данный момент сохранен лишь ID никнейм и имя, доступные для бота) нажмите кнопку \"Регистрация\"", menu)
			}
			return c.Send("Ошибка при сохранении ваших данных")
		} else if deleted {
			return c.Send("Приветствуем, "+c.Chat().Username+". Спасибо, что вернулись. Выберите действие", deletedSelector)
		}
		return c.Send("Приветствуем, "+c.Chat().Username+". Вы нажали кнопку старта. Выберите действие", menu)
	})

	bot.Handle(constants.ON_HELP, func(c telebot.Context) error {
		return c.Send("Это бот для работников ЦЦР. Сюда вы можете внести данные о своих пожелниях на день рождения для своих коллег")
	}, CheckDeleted(userService))

	bot.Handle(telebot.OnCallback, func(c telebot.Context) error {
		callback := c.Callback().Data[1:]
		id := strconv.FormatInt(c.Chat().ID, 10)
		page := strings.Split(callback, "|")[1]
		switch callback {
		case constants.BTN_REGISTER:
			return onButtonRegister(c, userService)
		case constants.BTN_HELP:
			return onButtonHelp(c)
		case constants.BTN_WISHLIST:
			return onButtonWishlist(c)
		case constants.BTN_ALL_USERS:
			return HandleUserList(c, userService)
		case constants.BTN_PREV:
			return onButtonPrev(c)
		case constants.BTN_RESTORE_USER:
			return onRestoreUser(c, userService)
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
			return onUserData(c, wishlistService)
		case constants.BACK_TO_LIST:
			return HandleUserList(c, userService)
		}
		return c.Respond()
	}, CheckDeleted(userService))

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
		default:
			return onError(c)
		}
	}, CheckDeleted(userService))

}
