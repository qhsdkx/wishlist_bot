package bot

import (
	"fmt"
	"strconv"
	"strings"
	constants "wishlist-bot/internal/constant"
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
		return r.userHandler.EditUserName(c)

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

func (r *HandlerRouter) UserData(c telebot.Context) error {
	data := c.Callback().Data[1:]
	if !strings.HasPrefix(data, constants.USER_DATA_PREFIX) {
		return c.Respond()
	}
	userId, _ := strconv.ParseInt(data[len(constants.USER_DATA_PREFIX):], 10, 64)
	wishes, err := r.wishlistHandler.service.FindAllByUserId(userId)
	if err != nil {
		return c.Edit(fmt.Sprintf("Ошибка в поиске пожеланий у юзера с айди %d", userId), MainMenu())
	}
	user, err := r.userHandler.service.FindByID(userId)
	if err != nil {
		return c.Edit("Почему-то не смогли найти этого пользователя в базе. Возвращаем в начало", MainMenu())
	}

	var msg strings.Builder
	msg.WriteString(fmt.Sprintf("День рождения у пользователя: %s", user.Birthdate.Format("02.01.2006")))
	msg.WriteString("\n🎁 Список желаний:\n\n")
	for _, wish := range wishes {
		msg.WriteString(fmt.Sprintf("• %s\n", wish.WishText))
	}

	_, err = c.Bot().Edit(c.Message(), msg.String(), createBackButton())
	if err != nil {
		return c.Respond(&telebot.CallbackResponse{
			Text: "Ошибка отображения данных",
		})
	}

	return c.Respond()
}
