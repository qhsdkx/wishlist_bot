package bot

import (
	consta "wishlist-bot/internal/constant"

	"gopkg.in/telebot.v4"
)

func MainMenu() *telebot.ReplyMarkup {
	menu := &telebot.ReplyMarkup{}
	btnMe := menu.Data("Редактировать мои данные", consta.BTN_ME)
	btnWishlist := menu.Data("Список желаний", consta.BTN_WISHLIST)
	btnAllUsers := menu.Data("Показать всех пользователей", consta.BTN_ALL_USERS)
	btnDeleteMe := menu.Data("Удалить меня в базе", consta.BTN_DELETE_ME)
	btnHelp := menu.Data("Помощь", consta.BTN_HELP)

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
	btnBirthdate := m.Data("Дата рождения", consta.BTN_EDIT_BIRTHDATE)
	btnSurname := m.Data("Фамилия", consta.BTN_EDIT_SURNAME)
	btnName := m.Data("Имя", consta.BTN_EDIT_NAME)
	btnUsername := m.Data("Никнейм в телеграме", consta.BTN_EDIT_USERNAME)
	btnPrev := m.Data("⬅", consta.BTN_PREV)

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
	btnShowWishlist := m.Data("Показать мои пожелания", consta.BTN_SHOW_ALL_WISHLIST)
	btnRegWishlist := m.Data("Внести пожелания", consta.BTN_REGISTER_WISHLIST)
	btnDeleteWish := m.Data("Удалить пожелание", consta.DELETE_WISH)
	btnPrev := m.Data("⬅", consta.BTN_PREV)

	m.Inline(
		m.Row(btnShowWishlist),
		m.Row(btnRegWishlist),
		m.Row(btnDeleteWish),
		m.Row(btnPrev),
	)
	return m
}

func BackOnlyMenu() *telebot.ReplyMarkup {
	m := &telebot.ReplyMarkup{}
	btnPrev := m.Data("⬅", consta.BTN_PREV)
	m.Inline(
		m.Row(btnPrev),
	)
	return m
}

func RegisterOnlyMenu() *telebot.ReplyMarkup {
	m := &telebot.ReplyMarkup{}
	btnReg := m.Data("Регистрация", consta.BTN_REGISTER)
	btnHelp := m.Data("Помощь", consta.BTN_HELP)

	m.Inline(
		m.Row(btnReg),
		m.Row(btnHelp),
	)
	return m
}
