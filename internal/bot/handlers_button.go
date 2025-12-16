package bot

import (
	consta "wishlist-bot/internal/constant"

	"gopkg.in/telebot.v4"
)

func MainMenu() *telebot.ReplyMarkup {
	menu := &telebot.ReplyMarkup{}
	btnMe := menu.Data("Редактировать мои данные", NewCallbackData(consta.BTN_ME, "", "", "").string())
	btnWishlist := menu.Data("Список желаний", NewCallbackData(consta.BTN_WISHLIST, "", "", "").string())
	btnAllUsers := menu.Data("Показать всех пользователей", NewCallbackData(consta.BTN_ALL_USERS, "", "", "").string())
	btnDeleteMe := menu.Data("Удалить меня в базе", NewCallbackData(consta.BTN_DELETE_ME, "", "", "").string())
	btnHelp := menu.Data("Помощь", NewCallbackData(consta.BTN_HELP, "", "", "").string())

	menu.Inline(
		menu.Row(btnMe),
		menu.Row(btnWishlist),
		menu.Row(btnAllUsers),
		menu.Row(btnDeleteMe),
		menu.Row(btnHelp),
	)
	return menu
}

func EditMenu() *telebot.ReplyMarkup {
	m := &telebot.ReplyMarkup{}
	btnBirthdate := m.Data("Дата рождения", NewCallbackData(consta.BTN_EDIT_BIRTHDATE, "", "", "").string())
	btnSurname := m.Data("Фамилия", NewCallbackData(consta.BTN_EDIT_SURNAME, "", "", "").string())
	btnName := m.Data("Имя", NewCallbackData(consta.BTN_EDIT_NAME, "", "", "").string())
	btnUsername := m.Data("Никнейм в телеграме", NewCallbackData(consta.BTN_EDIT_USERNAME, "", "", "").string())
	btnPrev := m.Data("⬅", NewCallbackData(consta.BTN_PREV, "", "", "").string())

	m.Inline(
		m.Row(btnBirthdate),
		m.Row(btnName),
		m.Row(btnSurname),
		m.Row(btnUsername),
		m.Row(btnPrev),
	)
	return m
}

func WishlistMenu() *telebot.ReplyMarkup {
	m := &telebot.ReplyMarkup{}
	btnRegWishlist := m.Data("Внести пожелания", NewCallbackData(consta.BTN_REGISTER_WISHLIST, "", "", "").string())
	btnDeleteWish := m.Data("Удалить пожелание", NewCallbackData(consta.DELETE_WISH, "", "", "").string())
	btnPrev := m.Data("⬅", NewCallbackData(consta.BTN_PREV, "", "", "").string())

	m.Inline(
		m.Row(btnRegWishlist),
		m.Row(btnDeleteWish),
		m.Row(btnPrev),
	)
	return m
}

func BackOnlyMenu() *telebot.ReplyMarkup {
	m := &telebot.ReplyMarkup{}
	btnPrev := m.Data("⬅", NewCallbackData(consta.BTN_PREV, "", "", "").string())
	m.Inline(
		m.Row(btnPrev),
	)
	return m
}

func RegisterOnlyMenu() *telebot.ReplyMarkup {
	m := &telebot.ReplyMarkup{}
	btnReg := m.Data("Регистрация", NewCallbackData(consta.BTN_REGISTER, "", "", "").string())
	btnHelp := m.Data("Помощь", consta.BTN_HELP)

	m.Inline(
		m.Row(btnReg),
		m.Row(btnHelp),
	)
	return m
}
