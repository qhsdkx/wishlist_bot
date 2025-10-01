package bot

import (
	"fmt"
	"os"
	"strconv"
	"strings"

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
}

func NewAdminHandler (us user.Service, ws wishlist.Service, state fsm.StateStore) *AdminHandler {
	return &AdminHandler{
		us: us,
		ws: ws,
		state: state,
	}
}

func (h AdminHandler) SendRegistered(c telebot.Context) error {
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
			return err
		}
	}
	return c.Send("Все гуд")
}

func (h AdminHandler) ShowUnregistered(c telebot.Context) error {
	id, _ := strconv.ParseInt(os.Getenv("ADMIN_ID"), 10, 64)
	if c.Chat().ID != id {
		return c.Send("Ты не админ, нельзя такое делать!!!")
	}
	users, err := h.us.FindAllUnregistered()
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
}

func (h AdminHandler) SendMessage(c telebot.Context, ID string) error {
	id, _ := strconv.ParseInt(os.Getenv("ADMIN_ID"), 10, 64)
	if c.Chat().ID != id {
		return c.Send("Ты не админ, нельзя такое делать!!!")
	}
	userID, err := strconv.ParseInt(ID, 10, 64)
	if err != nil {
		return c.Send("Невозможно отправить сообщение")
	}
	_, sErr := c.Bot().Send(telebot.ChatID(userID), c.Text())
	if sErr != nil {
		return c.Send("Ошибка отправления")
	}
	return nil
}
