package bot

import (
	"fmt"
	"gopkg.in/telebot.v4"
	"strings"
	constants "wishlist-bot/constant"
	sv "wishlist-bot/service"
)

func onButtonDeleteWish(c telebot.Context) error {
	states[c.Chat().ID] = constants.DELETE_WISH
	return c.Send("Введите ваше пожелание, которое хотите удалить (точно также, как написано выше)")
}

func onDeleteWish(c telebot.Context, wishlistService sv.WishService) error {
	delete(states, c.Chat().ID)
	err := wishlistService.Delete(c.Text(), c.Chat().ID)
	if err != nil {
		return c.Send(fmt.Sprintf("Ошибка удаления. Возможно, вы ввели неккоректное название\nВозвращаем в начало"), menu)
	}
	return c.Send(fmt.Sprintf("Пожелание успешно удалено\nВыберите действие"), wishlistSelector)
}

func onShowWishlist(c telebot.Context, service sv.WishService) error {
	wishes, err := service.FindAllByUserId(c.Chat().ID)
	if err != nil {
		return c.Respond(&telebot.CallbackResponse{
			Text: fmt.Sprintf("Ошибка поиска пожеланий юзера с айди %d", c.Chat().ID),
		})
	}

	var msg strings.Builder
	msg.WriteString("🎁 Список желаний:\n\n")
	for _, wish := range wishes {
		msg.WriteString(fmt.Sprintf("• %s\n", wish.WishText))
	}

	_, err = bot.Edit(c.Message(), msg.String(), onlyBack)
	if err != nil {
		return c.Respond(&telebot.CallbackResponse{
			Text: "Ошибка отображения данных",
		})
	}

	return c.Respond()
}

func onAwaitingWishlist(c telebot.Context, wishlistService sv.WishService) error {
	delete(states, c.Chat().ID)
	splits := strings.Split(c.Text(), ",")
	var wishes []sv.WishDto
	for _, split := range splits {
		wish := sv.WishDto{WishText: strings.TrimSpace(split), UserId: c.Chat().ID}
		wishes = append(wishes, wish)
	}
	err := wishlistService.SaveAll(wishes)
	if err != nil {
		c.Send(fmt.Sprintf("Ошибка во время сохранения %+v", err))
	}
	return c.Send("Ваш список желаний успешно сохранен", menu)
}

func onButtonWishlist(c telebot.Context, userService sv.UserService) error {
	err := userService.CheckIfRegistered(c.Chat().ID)
	if err != nil {
		_, err = bot.Edit(c.Message(), "Вы еще не зарегистрировались, чтобы пользоваться функционалом списка желаний. Пожалуйста, пройдите регистрацию", menu)
		return err
	}
	if _, err := bot.Edit(c.Message(), "Что хотите сделать?", wishlistSelector); err != nil {
		return c.Send(fmt.Sprintf("Что-то пошло не так\nв начало"), menu)
	}
	return nil
}

func onButtonRegWishList(c telebot.Context) error {
	states[c.Chat().ID] = constants.AWAITING_WISHES
	if _, err := bot.Edit(c.Message(), "Пожалуйста, введите ваши пожелания через запятую (майбах, шевроле камара, консервированные анансы)"); err != nil {
		return c.Send(fmt.Sprintf("Что-то пошло не так\nв начало"), menu)
	}
	return nil
}
