package bot

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"wishlist-bot/internal/logger/sl"

	constants "wishlist-bot/internal/constant"
	"wishlist-bot/internal/fsm"
	"wishlist-bot/internal/user"
	"wishlist-bot/internal/wishlist"

	"gopkg.in/telebot.v4"
)

type AdminHandler struct {
	us    user.Service
	ws    wishlist.Service
	state fsm.StateStore
	log   *slog.Logger
}

func NewAdminHandler(us user.Service, ws wishlist.Service, state fsm.StateStore, log *slog.Logger) *AdminHandler {
	return &AdminHandler{
		us:    us,
		ws:    ws,
		state: state,
		log:   log,
	}
}

func (h AdminHandler) SendRegistered(c telebot.Context) error {
	const op = "AdminHandler.SendRegistered"

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
	total, err := h.us.FindAllRegistered()
	if err != nil {
		h.log.Error(op, sl.Err(err))
		return c.Send("Ошибка при извлечении юзеров")
	}
	for _, user := range total {
		if user.ID == id {
			continue
		}
		return c.Send(telebot.ChatID(user.ID), message)
	}
	return c.Send("Все гуд")
}

func (h AdminHandler) SendUnregistered(c telebot.Context) error {
	const op = "AdminHandler.SendUnregistered"

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
	total, err := h.us.FindAllUnregistered()
	if err != nil {
		return c.Send("Ошибка при извлечении юзеров")
	}
	for _, user := range total {
		if user.ID == id {
			continue
		}
		_, err = c.Bot().Send(telebot.ChatID(user.ID), message)
		if err != nil {
			h.log.Error(op, sl.Err(err))
			return err
		}
	}
	return c.Send("Все гуд")
}

func (h AdminHandler) ShowUnregistered(c telebot.Context) error {
	const op = "AdminHandler.ShowUnregistered"

	id, _ := strconv.ParseInt(os.Getenv("ADMIN_ID"), 10, 64)
	if c.Chat().ID != id {
		return c.Send("Ты не админ, нельзя такое делать!!!")
	}
	users, err := h.us.FindAllUnregistered()
	if err != nil {
		h.log.Error(op, sl.Err(err))
		return c.Send("Ошибка при извлечении пользователей")
	}
	if users == nil {
		h.log.Error(op, sl.Err(err))
		return c.Send("Незарегистрированных пользователей нет!")
	}
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("Юзеры, которые не прошли полную регистрацию: \n\n"))
	for _, user := range users {
		builder.WriteString(fmt.Sprintf("- (%s) %s %s %s\n", user.Username, user.Name, user.Surname, user.Birthdate.Format("02.01.2006")))
	}
	return c.Send(builder.String())
}

func (h AdminHandler) SendMessage(c telebot.Context, ID string) error {
	const op = "AdminHandler.SendMessage"

	id, _ := strconv.ParseInt(os.Getenv("ADMIN_ID"), 10, 64)
	if c.Chat().ID != id {
		return c.Send("Ты не админ, нельзя такое делать!!!")
	}
	userID, err := strconv.ParseInt(ID, 10, 64)
	if err != nil {
		h.log.Error(op, sl.Err(err))
		return c.Send("Невозможно отправить сообщение")
	}
	_, sErr := c.Bot().Send(telebot.ChatID(userID), c.Text())
	if sErr != nil {
		h.log.Error(op, sl.Err(sErr))
		return c.Send("Ошибка отправления")
	}
	return nil
}
