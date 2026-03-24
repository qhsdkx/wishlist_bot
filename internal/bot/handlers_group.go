package bot

import (
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"
	constants "wishlist-bot/internal/constant"
	"wishlist-bot/internal/group"
	"wishlist-bot/internal/logger/sl"
	"wishlist-bot/internal/user"

	"gopkg.in/telebot.v4"
)

type GroupHandler struct {
	groupService    *group.Service
	userService     user.Service
	wishlistHandler *WishlistHandler
	log             *slog.Logger
}

func NewGroupHandler(groupService *group.Service, userService user.Service, wishlistHandler *WishlistHandler, log *slog.Logger) *GroupHandler {
	return &GroupHandler{
		groupService:    groupService,
		userService:     userService,
		wishlistHandler: wishlistHandler,
		log:             log,
	}
}

func (h *GroupHandler) ShowGroups(c telebot.Context) error {
	const op = "GroupHandler.ShowGroups"

	groups, err := h.groupService.GetGroupsForUser(c.Chat().ID)
	if err != nil {
		h.log.Error(op, sl.Err(err))
		return c.Edit("Ошибка при загрузке групп", MainMenu())
	}

	if len(groups) == 0 {
		return c.Edit("Пока нет активных групп. Они появятся, когда кому-то из коллег приблизится день рождения", MainMenu())
	}

	markup := h.createGroupsMarkup(groups)
	return c.Edit("🎉 Группы:\n\nВыберите группу:", markup)
}

func (h *GroupHandler) ShowGroup(c telebot.Context, groupID string) error {
	const op = "GroupHandler.ShowGroup"

	id, err := strconv.ParseInt(groupID, 10, 64)
	if err != nil {
		h.log.Error(op, sl.Err(err))
		return c.Edit("Ошибка", MainMenu())
	}

	g, err := h.groupService.GetGroupByID(id)
	if err != nil {
		h.log.Error(op, sl.Err(err))
		return c.Edit("Ошибка", MainMenu())
	}
	if g == nil {
		return c.Edit("Группа не найдена", MainMenu())
	}

	birthdayUser, err := h.userService.FindByID(g.BirthdayUserID)
	if err != nil {
		h.log.Error(op, sl.Err(err))
		return c.Edit("Ошибка", MainMenu())
	}

	members, err := h.groupService.GetGroupMembers(id)
	if err != nil {
		h.log.Error(op, sl.Err(err))
		return c.Edit("Ошибка", MainMenu())
	}

	statusText := "⏳ День рождения ещё не наступил"
	if g.Status == group.GroupStatusPassed {
		statusText = "✅ День рождения прошёл"
	}

	msg := strings.Builder{}
	msg.WriteString(fmt.Sprintf("%s\n\n", statusText))
	msg.WriteString(fmt.Sprintf("👤 Именинник: %s %s\n", birthdayUser.Name, birthdayUser.Surname))
	msg.WriteString(fmt.Sprintf("📅 Дата: %s\n\n", birthdayUser.Birthdate.Format("02.01.2006")))
	msg.WriteString(fmt.Sprintf("👥 Участники (%d):\n", len(members)+1))
	msg.WriteString(fmt.Sprintf("• %s %s (именинник)\n", birthdayUser.Name, birthdayUser.Surname))

	for _, m := range members {
		u, _ := h.userService.FindByID(m.UserID)
		if u.Name != "" {
			msg.WriteString(fmt.Sprintf("• %s %s\n", u.Name, u.Surname))
		}
	}

	markup := h.createGroupDetailMarkup(id, c.Chat().ID == g.BirthdayUserID, h.groupService.IsMember(id, c.Chat().ID), g.Status)

	return c.Edit(msg.String(), markup)
}

func (h *GroupHandler) JoinGroup(c telebot.Context, groupID string) error {
	const op = "GroupHandler.JoinGroup"

	h.log.Info("JoinGroup called", "group_id", groupID, "user_id", c.Chat().ID)

	id, err := strconv.ParseInt(groupID, 10, 64)
	if err != nil {
		h.log.Error(op, sl.Err(err))
		return c.Edit("Ошибка парсинга ID группы", MainMenu())
	}

	if err := h.groupService.JoinGroup(id, c.Chat().ID); err != nil {
		h.log.Error(op, sl.Err(err))
		return c.Edit(err.Error(), MainMenu())
	}

	return h.ShowGroup(c, groupID)
}

func (h *GroupHandler) LeaveGroup(c telebot.Context, groupID string) error {
	const op = "GroupHandler.LeaveGroup"

	id, err := strconv.ParseInt(groupID, 10, 64)
	if err != nil {
		h.log.Error(op, sl.Err(err))
		return c.Respond(&telebot.CallbackResponse{Text: "Ошибка"})
	}

	if err := h.groupService.LeaveGroup(id, c.Chat().ID); err != nil {
		h.log.Error(op, sl.Err(err))
		return c.Respond(&telebot.CallbackResponse{Text: "Ошибка при выходе из группы", ShowAlert: true})
	}

	return c.Respond(&telebot.CallbackResponse{Text: "Вы покинули группу", ShowAlert: false})
}

func (h *GroupHandler) createGroupsMarkup(groups []group.Group) *telebot.ReplyMarkup {
	markup := &telebot.ReplyMarkup{}
	rows := make([]telebot.Row, 0, len(groups)+1)

	for _, g := range groups {
		statusIcon := "⏳"
		if g.Status == group.GroupStatusPassed {
			statusIcon = "✅"
		}
		btn := markup.Data(
			fmt.Sprintf("%s %s", statusIcon, g.Name),
			NewCallbackData(
				constants.SHOW_GROUP,
				"",
				strconv.FormatInt(g.ID, 10),
				"",
			).string(),
		)
		rows = append(rows, markup.Row(btn))
	}

	btnPrev := markup.Data("⬅ Назад", NewCallbackData(constants.BTN_PREV, "", "", "").string())
	rows = append(rows, markup.Row(btnPrev))

	markup.Inline(rows...)
	return markup
}

func (h *GroupHandler) createGroupDetailMarkup(groupID int64, isBirthdayUser, isMember bool, status string) *telebot.ReplyMarkup {
	markup := &telebot.ReplyMarkup{}
	rows := make([]telebot.Row, 0, 4)

	if !isBirthdayUser {
		if isMember {
			btnWishes := markup.Data(
				"🎁 Посмотреть пожелания",
				NewCallbackData(
					constants.SHOW_BD_WISHES,
					"",
					strconv.FormatInt(groupID, 10),
					"",
				).string(),
			)
			rows = append(rows, markup.Row(btnWishes))

			btnLeave := markup.Data(
				"🚪 Покинуть группу",
				NewCallbackData(
					constants.LEAVE_GROUP,
					"",
					strconv.FormatInt(groupID, 10),
					"",
				).string(),
			)
			rows = append(rows, markup.Row(btnLeave))
		} else if status == group.GroupStatusUpcoming {
			btnJoin := markup.Data(
				"✅ Присоединиться",
				NewCallbackData(
					constants.JOIN_GROUP,
					"",
					strconv.FormatInt(groupID, 10),
					"",
				).string(),
			)
			rows = append(rows, markup.Row(btnJoin))
		}
	}

	btnBack := markup.Data(
		"⬅ К списку групп",
		NewCallbackData(constants.BACK_TO_GROUPS, "", "", "").string(),
	)
	rows = append(rows, markup.Row(btnBack))

	markup.Inline(rows...)
	return markup
}

func (h *GroupHandler) ShowBirthdayWishes(c telebot.Context, groupID string) error {
	const op = "GroupHandler.ShowBirthdayWishes"

	id, err := strconv.ParseInt(groupID, 10, 64)
	if err != nil {
		h.log.Error(op, sl.Err(err))
		return c.Edit("Ошибка", MainMenu())
	}

	g, err := h.groupService.GetGroupByID(id)
	if err != nil || g == nil {
		h.log.Error(op, sl.Err(err))
		return c.Edit("Группа не найдена", MainMenu())
	}

	birthdayUser, err := h.userService.FindByID(g.BirthdayUserID)
	if err != nil {
		h.log.Error(op, sl.Err(err))
		return c.Edit("Ошибка", MainMenu())
	}

	wishes, err := h.wishlistHandler.service.FindAllByUserID(g.BirthdayUserID)
	if err != nil {
		h.log.Error(op, sl.Err(err))
		return c.Edit("Ошибка при загрузке пожеланий", MainMenu())
	}

	msg := strings.Builder{}
	msg.WriteString(fmt.Sprintf("🎁 Пожелания %s %s:\n\n", birthdayUser.Name, birthdayUser.Surname))

	if len(wishes) == 0 {
		msg.WriteString("(список пожеланий пуст)")
	} else {
		for _, w := range wishes {
			msg.WriteString(fmt.Sprintf("• %s\n", w.WishText))
		}
	}

	markup := &telebot.ReplyMarkup{}
	btnBack := markup.Data(
		"⬅ К группе",
		NewCallbackData(
			constants.SHOW_GROUP,
			"",
			groupID,
			"",
		).string(),
	)
	markup.Inline(markup.Row(btnBack))

	return c.Edit(msg.String(), markup)
}

func formatDaysUntilBirthday(birthdate time.Time) string {
	now := time.Now()
	thisYearBirthday := time.Date(now.Year(), birthdate.Month(), birthdate.Day(), 0, 0, 0, 0, now.Location())

	if thisYearBirthday.Before(now) {
		thisYearBirthday = thisYearBirthday.AddDate(1, 0, 0)
	}

	days := int(thisYearBirthday.Sub(now).Hours() / 24)

	switch {
	case days == 0:
		return "сегодня!"
	case days == 1:
		return "завтра!"
	case days <= 7:
		return fmt.Sprintf("через %d дней", days)
	default:
		return fmt.Sprintf("через %d дней", days)
	}
}
