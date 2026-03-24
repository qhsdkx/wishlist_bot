package bot

import (
	"errors"
	"fmt"
	"log"
	"log/slog"
	"strconv"
	"strings"
	constants "wishlist-bot/internal/constant"
	"wishlist-bot/internal/fsm"
	"wishlist-bot/internal/logger/sl"

	"gopkg.in/telebot.v4"
)

type HandlerRouter struct {
	userHandler     *UserHandler
	wishlistHandler *WishlistHandler
	adminHandler    *AdminHandler
	states          fsm.StateStore
	log             *slog.Logger
}

func NewHandlerRouter(user *UserHandler, wishlist *WishlistHandler, admin *AdminHandler, states fsm.StateStore, log *slog.Logger) HandlerRouter {
	return HandlerRouter{
		userHandler:     user,
		wishlistHandler: wishlist,
		adminHandler:    admin,
		states:          states,
		log:             log,
	}
}

func (r *HandlerRouter) OnCallback(c telebot.Context) error {
	callbackData := parseCallback(c.Callback().Data[1:])

	r.log.Info("Get callback data", "callback", callbackData)

	switch callbackData.action {
	case constants.BTN_EDIT_NAME:
		return r.userHandler.EditName(c)
	case constants.BTN_EDIT_SURNAME:
		return r.userHandler.EditSurname(c)
	case constants.BTN_EDIT_BIRTHDATE:
		return r.userHandler.EditBirthdate(c)
	case constants.BTN_EDIT_USERNAME:
		return r.userHandler.EditUserName(c)
	case constants.BTN_ME:
		return r.userHandler.ShowProfile(c)
	case constants.BTN_REGISTER:
		return r.userHandler.Register(c)
	case constants.BTN_HELP:
		return r.Help(c)
	case constants.BTN_WISHLIST:
		return r.wishlistHandler.Show(c)
	case constants.BTN_ALL_USERS:
		return r.userHandler.UserList(c, constants.SHOW_USERS)
	case constants.BTN_PREV:
		return r.userHandler.Prev(c)
	case constants.BTN_DELETE_ME:
		return r.userHandler.DeleteMe(c)
	case constants.BTN_PREV_PAGE:
		return r.userHandler.PrevAndBack(c, callbackData.mode, callbackData.page)
	case constants.BTN_NEXT_PAGE:
		return r.userHandler.PrevAndBack(c, callbackData.mode, callbackData.page)
	case constants.USER_DATA_PREFIX:
		return r.UserData(c, callbackData)
	case constants.BACK_TO_LIST:
		return r.userHandler.UserList(c, constants.SHOW_USERS)
	case constants.SEND_MESSAGE_ADMIN:
		return r.userHandler.UserList(c, constants.SEND_MESSAGE_ADMIN)

	// кнопки wishlist
	case constants.BTN_REGISTER_WISHLIST:
		return r.wishlistHandler.Register(c)
	case constants.DELETE_WISH:
		return r.wishlistHandler.Delete(c)
	case constants.DELETE_CHOOSED_WISH:
		return r.wishlistHandler.AwaitingDelete(c, callbackData)
	default:
		return r.Error(c)
	}
}

func (r *HandlerRouter) OnText(c telebot.Context) error {
	state, ok := r.states.Get(c.Chat().ID)
	r.log.Info("Get state", "state", state)
	if ok != nil {
		r.log.Error("Error get state", "err", ok)
		return c.Send("Пожалуйста, начните с /start")
	}

	id, err := parseID(state)
	if err != nil {
		r.log.Error("Error parse id", sl.Err(err))
		log.Printf("Can't parse id or another state: %v", err)
	}

	switch state {
	case constants.AWAITING_NAME:
		return r.userHandler.AwaitingName(c)
	case constants.AWAITING_SURNAME:
		return r.userHandler.AwaitingSurname(c)
	case constants.AWAITING_BIRTHDATE:
		return r.userHandler.AwaitingBirthdate(c)
	case constants.AWAITING_NEW_NAME:
		return r.userHandler.AwaitingNewName(c)
	case constants.AWAITING_NEW_SURNAME:
		return r.userHandler.AwaitingNewSurname(c)
	case constants.AWAITING_NEW_BIRTHDATE:
		return r.userHandler.AwaitingNewBirthdate(c)
	case constants.AWAITING_NEW_USERNAME:
		return r.userHandler.AwaitingNewUsername(c)

	case constants.AWAITING_WISHES:
		return r.wishlistHandler.Awaiting(c)
	case constants.SEND_MESSAGE_ADMIN:
		return r.adminHandler.SendMessage(c, id)

	default:
		return c.Send("Неизвестное состояние, возвращаем в главное меню", MainMenu())
	}

}

func (r *HandlerRouter) OnStart(c telebot.Context) error {
	const op = "HandlerRouter.OnStart"

	u, err := r.userHandler.service.FindByID(c.Chat().ID)
	if err != nil || u.Status != "REGISTERED" {
		var msg strings.Builder

		msg.WriteString("ВНИМАНИЕ!\n\n")
		msg.WriteString("Нажимая кнопку 'Регистрация', вы даете согласие на обработку ваших персональных данных:\n")
		msg.WriteString("• Имя\n")
		msg.WriteString("• Дата рождения\n\n")
		msg.WriteString("Данные будут использованы для создания вашей учетной записи в системе.")
		msg.WriteString("Для дальнейшего использования нужно зарегистрироваться")

		return c.Send(msg.String(), RegisterOnlyMenu())
	}
	return c.Send(fmt.Sprintf("С возвращением, %s!", u.Username), MainMenu())
}

func (r *HandlerRouter) UserData(c telebot.Context, cb CallbackData) error {
	const op = "HandlerRouter.UserData"

	data := c.Callback().Data[1:]
	if !strings.HasPrefix(data, constants.USER_DATA_PREFIX) {
		return c.Respond()
	}
	userId, _ := strconv.ParseInt(cb.id, 10, 64)
	r.log.Info("Get user id", "user_id", userId)
	wishes, err := r.wishlistHandler.service.FindAllByUserID(userId)

	if err != nil {
		r.log.Error(op, sl.Err(err))
		return c.Edit(fmt.Sprintf("Ошибка в поиске пожеланий у юзера с айди %d", userId), MainMenu())
	}
	r.log.Info("Get wishes", "wishes", wishes)

	user, err := r.userHandler.service.FindByID(userId)
	if err != nil {
		r.log.Error(op, sl.Err(err))
		return c.Edit("Почему-то не смогли найти этого пользователя в базе. Возвращаем в начало", MainMenu())
	}

	var msg strings.Builder
	msg.WriteString(fmt.Sprintf("День рождения у пользователя: %s", user.Birthdate.Format("02.01.2006")))
	msg.WriteString("\n🎁 Список желаний:\n\n")
	for _, wish := range wishes {
		msg.WriteString(fmt.Sprintf("• %s\n", wish.WishText))
	}

	_, err = c.Bot().Edit(c.Message(), msg.String(), r.createBackButton())
	if err != nil {
		r.log.Error(op, sl.Err(err))
		return c.Respond(&telebot.CallbackResponse{
			Text: "Ошибка отображения данных",
		})
	}

	return c.Respond()
}

func parseID(state string) (string, error) {
	if strings.HasPrefix(state, constants.SEND_MESSAGE_ADMIN+"_") {
		after, found := strings.CutPrefix(state, constants.SEND_MESSAGE_ADMIN+"_")
		if !found {
			return "", errors.New("Error with state")
		}
		return after, nil
	}
	return state, nil
}
