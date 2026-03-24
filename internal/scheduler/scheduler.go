package scheduler

import (
	"fmt"
	"log"
	"log/slog"
	"strings"
	"time"
	"wishlist-bot/internal/config"
	constants "wishlist-bot/internal/constant"
	"wishlist-bot/internal/group"
	"wishlist-bot/internal/logger/sl"
	"wishlist-bot/internal/user"
	"wishlist-bot/internal/wishlist"

	"github.com/go-co-op/gocron"
	"gopkg.in/telebot.v4"
)

type Scheduler struct {
	bot *telebot.Bot
	us  user.Service
	ws  wishlist.Service
	gs  *group.Service
	s   *gocron.Scheduler
	cfg *config.Config
	log *slog.Logger
}

func New(bot *telebot.Bot, us user.Service, ws wishlist.Service, gs *group.Service, cfg *config.Config, log *slog.Logger) *Scheduler {
	location, _ := time.LoadLocation(constants.LOCATION)
	return &Scheduler{
		bot: bot,
		us:  us,
		ws:  ws,
		gs:  gs,
		s:   gocron.NewScheduler(location),
		cfg: cfg,
		log: log,
	}
}

func (sch *Scheduler) Start() {
	const op = "Scheduler.Start"

	scheduleTime := sch.cfg.Scheduler.NotifTimeDaily
	if scheduleTime == "" {
		scheduleTime = "10:00"
	}

	scheduleTimeWeekly := sch.cfg.Scheduler.NotifTimeWeekly
	if scheduleTimeWeekly == "" {
		scheduleTimeWeekly = "9:55"
	}

	_, err := sch.s.Every(1).Days().At(scheduleTime).Do(sch.sendDailyNotifications)
	if err != nil {
		sch.log.Error(op, sl.Err(err))
	}

	_, err = sch.s.Every(1).Days().At(scheduleTimeWeekly).Do(sch.sendWeeklyNotifications)
	if err != nil {
		sch.log.Error(op, sl.Err(err))
	}

	_, err = sch.s.Every(1).Days().At("23:50").Do(sch.processBirthday)
	if err != nil {
		sch.log.Error(op, sl.Err(err))
	}

	_, err = sch.s.Every(1).Days().At("00:01").Do(sch.cleanupOldGroups)
	if err != nil {
		sch.log.Error(op, sl.Err(err))
	}

	sch.s.StartAsync()
}

func (sch *Scheduler) sendDailyNotifications() {
	const op = "Scheduler.SendDailyNotifications"

	users, err := sch.us.FindAllRegistered()
	if err != nil {
		sch.log.Error(op, sl.Err(err))
		_, sendErr := sch.bot.Send(telebot.ChatID(sch.cfg.AdminId), fmt.Sprintf("Ошибка в уведомлениях"))
		if sendErr != nil {
			log.Printf("Error sending daily notifications: %v", sendErr)
		}
	}

	birthdayTomorrow, others := splitUsersByBirthday(users, 1)
	response := makeResponse(birthdayTomorrow)

	if len(birthdayTomorrow) > 0 {
		for _, other := range others {
			_, err = sch.bot.Send(telebot.ChatID(other.ID), response)
			if err != nil {
				sch.log.Error(op, sl.Err(err))
				log.Printf("Failed to send to user %d: %v", other.ID, err)
			}
		}
	}
}

func (sch *Scheduler) sendWeeklyNotifications() {
	const op = "Scheduler.SendWeeklyNotifications"

	users, err := sch.us.FindAllRegistered()
	if err != nil {
		sch.log.Error(op, sl.Err(err))
		_, sendErr := sch.bot.Send(telebot.ChatID(sch.cfg.AdminId), fmt.Sprintf("Ошибка в уведомлениях"))
		if sendErr != nil {
			log.Printf("Error sending weekly notifications: %v", sendErr)
		}
		return
	}

	birthdayInWeek, others := splitUsersByBirthday(users, 7)

	for _, birthdayUser := range birthdayInWeek {
		g, err := sch.gs.CreateGroupForBirthday(birthdayUser.ID)
		if err != nil {
			sch.log.Error(op, sl.Err(err))
			continue
		}

		response := fmt.Sprintf("🎉 Через неделю (%s) день рождения у %s %s!\nХотите присоединиться к группе?",
			birthdayUser.Birthdate.Format("02.01.2006"), birthdayUser.Name, birthdayUser.Surname)

		for _, other := range others {
			markup := &telebot.ReplyMarkup{}
			btnJoin := markup.Data("✅ Присоединиться", fmt.Sprintf("jg|%d", g.ID))
			markup.Inline(markup.Row(btnJoin))

			_, err = sch.bot.Send(telebot.ChatID(other.ID), response, markup)
			if err != nil {
				sch.log.Error(op, sl.Err(err))
			}
		}
	}
}

func (sch *Scheduler) processBirthday() {
	const op = "Scheduler.ProcessBirthday"

	users, err := sch.us.FindAllRegistered()
	if err != nil {
		sch.log.Error(op, sl.Err(err))
		return
	}

	birthdayToday, _ := splitUsersByBirthday(users, 0)

	for _, bd := range birthdayToday {
		sch.ws.DeleteAll(bd.ID)

		g, err := sch.gs.FindByBirthdayUserID(bd.ID)
		if err == nil && g != nil {
			sch.gs.MarkGroupAsPassed(g.ID)
		}
	}
}

func (sch *Scheduler) cleanupOldGroups() {
	const op = "Scheduler.CleanupOldGroups"

	if err := sch.gs.CleanupOldGroups(); err != nil {
		sch.log.Error(op, sl.Err(err))
	}
}

func (sch *Scheduler) Stop() {
	sch.s.Stop()
}

func splitUsersByBirthday(users []user.User, daysBefore int) (birthdayUsers []user.User, others []user.User) {
	target := time.Now().AddDate(0, 0, daysBefore)
	targetMonth := target.Month()
	targetDay := target.Day()

	for _, u := range users {
		bdMonth := u.Birthdate.Month()
		bdDay := u.Birthdate.Day()

		if bdMonth == targetMonth && bdDay == targetDay {
			birthdayUsers = append(birthdayUsers, u)
		} else {
			others = append(others, u)
		}
	}

	return birthdayUsers, others
}

func makeResponse(users []user.User) string {
	var response strings.Builder
	now := time.Now().AddDate(0, 0, 1)
	response.WriteString(fmt.Sprintf("Доброе утро!\nЗавтра (%s) день рождения у:\n\n", now.Format("02.01.2006")))
	for _, u := range users {
		response.WriteString(fmt.Sprintf("(%s) %s %s\n\n", u.Username, u.Surname, u.Name))
	}
	return response.String()
}

func makeWeeklyResponse(users []user.User) string {
	var response strings.Builder
	response.WriteString("Доброе утро!\nЧерез неделю будет день рождения у:\n\n")
	for _, u := range users {
		response.WriteString(fmt.Sprintf("(%s) %s %s\nДень рождения: %s\n\n", u.Username, u.Surname, u.Name, u.Birthdate.Format("02.01.2006")))
	}
	return response.String()
}
