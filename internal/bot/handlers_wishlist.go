package bot

import (
	"fmt"
	"strconv"
	"strings"
	constants "wishlist-bot/internal/constant"
	"wishlist-bot/internal/fsm"
	"wishlist-bot/internal/wishlist"

	"gopkg.in/telebot.v4"
)

type WishlistHandler struct {
	service wishlist.Service
	states  fsm.StateStore
}

func NewWishlistHandler(service wishlist.Service, states fsm.StateStore) *WishlistHandler {
	return &WishlistHandler{service: service, states: states}
}

func (h *WishlistHandler) Show(c telebot.Context) error {
	wishes, err := h.service.FindAllByUserId(c.Chat().ID)
	if err != nil {
		return c.Edit("Ошибка при поиске ваших пожеланий", MainMenu())
	}

	var msg strings.Builder
	msg.WriteString("🎁 Ваши пожелания:\n\n")
	for _, w := range wishes {
		msg.WriteString(fmt.Sprintf("• %s\n", w.WishText))
	}

	return c.Edit(msg.String(), WishlistMenu())
}

func (h *WishlistHandler) Register(c telebot.Context) error {
	h.states.Set(c.Chat().ID, constants.AWAITING_WISHES)
	return c.Edit("Введите ваши пожелания через запятую")
}

func (h *WishlistHandler) Awaiting(c telebot.Context) error {
	text := c.Text()
	splits := strings.Split(text, ",")
	wishes := make([]wishlist.Wish, 0, len(splits))
	for _, s := range splits {
		wishes = append(wishes, wishlist.Wish{
			UserID:   c.Chat().ID,
			WishText: strings.TrimSpace(s),
		})
	}

	if err := h.service.SaveAll(wishes); err != nil {
		return c.Send(fmt.Sprintf("Ошибка сохранения: %+v", err), MainMenu())
	}

	h.states.Delete(c.Chat().ID)
	return c.Send("Список пожеланий успешно сохранен", MainMenu())
}

func (h *WishlistHandler) Delete(c telebot.Context) error {
	h.states.Set(c.Chat().ID, "DELETE_WISH")
	return c.Edit("Введите пожелание, которое хотите удалить")
}

func (h *WishlistHandler) AwaitingDelete(c telebot.Context) error {
	text := c.Text()
	h.states.Delete(c.Chat().ID)

	if err := h.service.Delete(text, c.Chat().ID); err != nil {
		return c.Send("Ошибка при удалении. Проверьте название", WishlistMenu())
	}

	return c.Send("Пожелание успешно удалено", WishlistMenu())
}

func (h *WishlistHandler) HandleDeleteWish(c telebot.Context) error {
	idStr := c.Callback().Data
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.Respond(&telebot.CallbackResponse{Text: "Некорректный ID"})
	}

	err = h.service.DeleteByID(id)
	if err != nil {
		return c.Respond(&telebot.CallbackResponse{Text: "Ошибка удаления"})
	}

	return c.Respond(&telebot.CallbackResponse{Text: "Удалено ✅", ShowAlert: false})
}
