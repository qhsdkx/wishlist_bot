package bot

import (
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	constants "wishlist-bot/internal/constant"
	"wishlist-bot/internal/fsm"
	"wishlist-bot/internal/logger/sl"
	"wishlist-bot/internal/wishlist"

	"gopkg.in/telebot.v4"
)

type WishlistHandler struct {
	service wishlist.Service
	states  fsm.StateStore
	log     *slog.Logger
}

func NewWishlistHandler(service wishlist.Service, states fsm.StateStore, log *slog.Logger) *WishlistHandler {
	return &WishlistHandler{service: service, states: states, log: log}
}

func (h *WishlistHandler) Show(c telebot.Context) error {
	const op = "WishlistHandler.Show"

	wishes, err := h.service.FindAllByUserID(c.Chat().ID)
	if err != nil {
		h.log.Error(op, sl.Err(err))
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
	const op = "WishlistHandler.Awaiting"

	count, err := h.service.FindCountByUserID(c.Chat().ID)
	if err != nil {
		h.log.Error(op, sl.Err(err))
		return c.Send("Ошибка при сохранении. Напишите админу @qhsdkx")
	}

	if count > 15 {
		h.log.Error(op, slog.String("count", strconv.Itoa(count)+" больше 15"))
		return c.Send("Временно доступно не больше 15 пожеланий. Удалите старые и добавьте новые, если хотите их изменить")
	}

	text := c.Text()
	splits := strings.Split(text, ",")
	if len(splits)-count == 0 {
		h.log.Error(op, slog.String("count", strconv.Itoa(count)+" больше 15"))
		return c.Send("Временно доступно не больше 15 пожеланий. Удалите старые и добавьте новые, если хотите их изменить")
	}

	wishes := make([]wishlist.Wish, 0, len(splits))
	for _, s := range splits {
		wishes = append(wishes, wishlist.Wish{
			UserID:   c.Chat().ID,
			WishText: strings.TrimSpace(s),
		})
	}

	if errSave := h.service.SaveAll(wishes); errSave != nil {
		h.log.Error(op, sl.Err(errSave))
		return c.Send(fmt.Sprintf("Ошибка сохранения: %+v", errSave), MainMenu())
	}

	h.states.Delete(c.Chat().ID)
	return c.Send("Список пожеланий успешно сохранен", MainMenu())
}

func (h *WishlistHandler) Delete(c telebot.Context) error {
	const op = "WishlistHandler.Delete"

	wishes, err := h.service.FindAllByUserID(c.Chat().ID)
	if err != nil {
		h.log.Error(op, sl.Err(err))
		return err
	}
	markup := h.createWishlistMarkup(wishes)
	return c.Edit("Выберите пожелание, которое хотите удалить", markup)
}

func (h *WishlistHandler) AwaitingDelete(c telebot.Context, cb CallbackData) error {
	const op = "WishlistHandler.AwaitingDelete"
	wishID, err := strconv.ParseInt(cb.Id(), 10, 64)
	if err != nil {
		h.log.Error(op, sl.Err(err))
	}

	if err := h.service.DeleteByID(wishID); err != nil {
		h.log.Error(op, sl.Err(err))
		return c.Send("Ошибка при удалении", WishlistMenu())
	}

	return c.Edit("Пожелание успешно удалено", WishlistMenu())
}

func (h *WishlistHandler) HandleDeleteWish(c telebot.Context) error {
	const op = "WishlistHandler.HandleDeleteWish"

	idStr := c.Callback().Data
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.log.Error(op, sl.Err(err))
		return c.Respond(&telebot.CallbackResponse{Text: "Некорректный ID"})
	}

	err = h.service.DeleteByID(id)
	if err != nil {
		h.log.Error(op, sl.Err(err))
		return c.Respond(&telebot.CallbackResponse{Text: "Ошибка удаления"})
	}

	return c.Respond(&telebot.CallbackResponse{Text: "Удалено ✅", ShowAlert: false})
}

func (h *WishlistHandler) createWishlistMarkup(wishes []wishlist.Wish) *telebot.ReplyMarkup {
	markup := &telebot.ReplyMarkup{}
	rows := make([]telebot.Row, 0, len(wishes)+1)

	for _, wish := range wishes {
		btn := markup.Data(
			fmt.Sprintf(wish.WishText),
			NewCallbackData(
				constants.DELETE_CHOOSED_WISH,
				"",
				strconv.FormatInt(wish.ID, 10),
				"",
			).string())
		rows = append(rows, markup.Row(btn))
	}

	btnPrev := markup.Data("⬅", NewCallbackData(constants.BTN_PREV, "", "", "").string())
	rows = append(rows, markup.Row(btnPrev))

	markup.Inline(rows...)
	return markup
}
