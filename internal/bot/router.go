package bot

import (
	"strings"
	"wishlist-bot/internal/fsm"

	"gopkg.in/telebot.v4"
)

type HandlerRouter struct {
	userHandler     *UserHandler
	wishlistHandler *WishlistHandler
	states          fsm.StateStore
}

func NewHandlerRouter(user *UserHandler, wishlist *WishlistHandler, states fsm.StateStore) *HandlerRouter {
	return &HandlerRouter{
		userHandler:     user,
		wishlistHandler: wishlist,
		states:          states,
	}
}

func (r *HandlerRouter) OnCallback(c telebot.Context) error {
	data := c.Callback().Data

	switch data {
	case "EDIT_NAME":
		return r.userHandler.EditName(c)
	case "EDIT_SURNAME":
		return r.userHandler.EditSurname(c)
	case "EDIT_BIRTHDATE":
		return r.userHandler.EditBirthdate(c)
	case "EDIT_USERNAME":
		return r.userHandler.EditUsername(c)

	// кнопки wishlist
	case "REGISTER_WISHES":
		return r.wishlistHandler.Register(c)
	case "DELETE_WISH":
		return r.wishlistHandler.Delete(c)
	case "SHOW_WISHES":
		return r.wishlistHandler.Show(c)

	case "MAIN_MENU":
		return c.Edit("Возвращаем в главное меню", MainMenu())
	}

	if strings.HasPrefix(data, "SEND_MESSAGE_ADMIN_") {
		idStr := strings.TrimPrefix(data, "SEND_MESSAGE_ADMIN_")
		r.states.Set(c.Chat().ID, "SEND_MESSAGE_ADMIN_"+idStr)
		return c.Edit("Введите сообщение для данного пользователя")
	}

	return c.Respond()
}

func (r *HandlerRouter) OnText(c telebot.Context) error {
	state, ok := r.states.Get(c.Chat().ID)
	if ok != nil {
		return c.Send("Пожалуйста, начните с /start")
	}

	switch state {
	case "AWAITING_NEW_NAME":
		return r.userHandler.AwaitingNewName(c)
	case "AWAITING_NEW_SURNAME":
		return r.userHandler.AwaitingNewSurname(c)
	case "AWAITING_NEW_BIRTHDATE":
		return r.userHandler.AwaitingNewBirthdate(c)
	case "AWAITING_NEW_USERNAME":
		return r.userHandler.AwaitingNewUsername(c)

	case "AWAITING_WISHES":
		return r.wishlistHandler.Awaiting(c)
	case "DELETE_WISH":
		return r.wishlistHandler.AwaitingDelete(c)

	default:
		return c.Send("Неизвестное состояние, возвращаем в главное меню", MainMenu())
	}

}

func (r *HandlerRouter) OnStart(c telebot.Context) error {
	menu := MainMenu()
	return c.Send("Привет! Выберите действие:", menu)
}
