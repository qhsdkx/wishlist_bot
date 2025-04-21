package bot

import (
	"gopkg.in/telebot.v4"
	consta "wishlist-bot/constant"
)

var (
	menu                 = &telebot.ReplyMarkup{OneTimeKeyboard: true}
	regSelector          = &telebot.ReplyMarkup{OneTimeKeyboard: true}
	wishlistSelector     = &telebot.ReplyMarkup{OneTimeKeyboard: true}
	deletedSelector      = &telebot.ReplyMarkup{OneTimeKeyboard: true}
	wantEditSelector     = &telebot.ReplyMarkup{OneTimeKeyboard: true}
	onlyRegisterSelector = &telebot.ReplyMarkup{OneTimeKeyboard: true}
)

var (
	btnHelp     = menu.Data("Помощь", consta.BTN_HELP)
	btnRegister = menu.Data("Регистрация", consta.BTN_REGISTER)
	btnWishlist = menu.Data("Список желаний", consta.BTN_WISHLIST)
	btnAllUsers = menu.Data("Показать всех пользователей", consta.BTN_ALL_USERS)
	btnRestore  = deletedSelector.Data("Восстановить меня в системе", consta.BTN_RESTORE_USER)
	btnDeleteMe = menu.Data("Удалить меня в базе", consta.BTN_DELETE_ME)
	btnMe       = menu.Data("Мои данные", consta.BTN_ME)
)

var (
	//btnShowWishlist = menu.Data("Показать мои пожелания", consta.BTN_SHOW_ALL_WISHLIST)
	//btnRegWishlist  = menu.Data("Внести пожелания", consta.BTN_REGISTER_WISHLIST)

	btnBirthdate = wantEditSelector.Data("Дата рождения", consta.BTN_EDIT_BIRTHDATE)
	btnSurname   = wantEditSelector.Data("Фамилия", consta.BTN_EDIT_SURNAME)
	btnName      = wantEditSelector.Data("Имя", consta.BTN_EDIT_NAME)
	btnUsername  = wantEditSelector.Data("Никнейм в телеграме", consta.BTN_EDIT_USERNAME)
	btnPrev      = wantEditSelector.Data("⬅", consta.BTN_PREV)
)

func setUpButtons() {
	menu.Inline(
		menu.Row(btnMe),
		menu.Row(btnWishlist),
		menu.Row(btnRegister),
		menu.Row(btnAllUsers),
		menu.Row(btnDeleteMe),
		menu.Row(btnHelp),
	)

	wantEditSelector.Inline(
		wantEditSelector.Row(btnBirthdate),
		wantEditSelector.Row(btnName),
		wantEditSelector.Row(btnSurname),
		wantEditSelector.Row(btnUsername),
		wantEditSelector.Row(btnPrev),
	)

	deletedSelector.Inline(
		deletedSelector.Row(btnRestore),
	)

	onlyRegisterSelector.Inline(
		menu.Row(btnRegister),
		menu.Row(btnPrev),
	)

	//wishlistSelector.Inline(
	//	wishlistSelector.Row(btnShowWishlist),
	//	wishlistSelector.Row(btnRegWishlist),
	//)
}
