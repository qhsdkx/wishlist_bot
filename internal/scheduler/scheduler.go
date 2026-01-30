package scheduler

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
	"wishlist-bot/internal/config"
	constants "wishlist-bot/internal/constant"
	"wishlist-bot/internal/user"
	"wishlist-bot/internal/wishlist"

	"github.com/go-co-op/gocron"
	"gopkg.in/telebot.v4"
)

type Scheduler struct {
	bot *telebot.Bot
	us  user.Service
	ws  wishlist.Service
	s   *gocron.Scheduler
	cfg *config.Config
}

func New(bot *telebot.Bot, us user.Service, ws wishlist.Service, cfg *config.Config) *Scheduler {
	location, _ := time.LoadLocation(constants.LOCATION)
	return &Scheduler{
		bot: bot,
		us:  us,
		ws:  ws,
		s:   gocron.NewScheduler(location),
		cfg: cfg,
	}
}

func (sch *Scheduler) Start() {
	scheduleTime := os.Getenv(sch.cfg.Scheduler.NotifTimeDaily)
	if scheduleTime == "" {
		scheduleTime = "10:00"
	}

	scheduleTimeWeekly := os.Getenv(sch.cfg.Scheduler.NotifTimeWeekly)
	if scheduleTimeWeekly == "" {
		scheduleTimeWeekly = "9:55"
	}

	_, err := sch.s.Every(1).Days().At(scheduleTime).Do(sch.sendDailyNotifications)
	if err != nil {
		log.Fatalf("Error scheduling task: %v", err)
	}

	_, err = sch.s.Every(1).Days().At(scheduleTimeWeekly).Do(sch.sendWeeklyNotifications)
	if err != nil {
		log.Fatalf("Error scheduling task: %v", err)
	}

	_, err = sch.s.Every(1).Days().At("23:50").Do(sch.deleteWishes)

	sch.s.StartAsync()
}

func (sch *Scheduler) sendDailyNotifications() {
	users, err := sch.us.FindAllRegistered()
	if err != nil {
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
				log.Printf("Failed to send to user %d: %v", other.ID, err)
			}
		}
	}
}

func (sch *Scheduler) sendWeeklyNotifications() {
	id, _ := strconv.Atoi(os.Getenv("ADMIN_ID"))
	users, err := sch.us.FindAllRegistered()
	if err != nil {
		_, sendErr := sch.bot.Send(telebot.ChatID(id), fmt.Sprintf("Ошибка в уведомлениях"))
		if sendErr != nil {
			log.Printf("Error sending daily notifications: %v", sendErr)
		}
	}

	birthdayInWeek, others := splitUsersByBirthday(users, 7)
	response := makeWeeklyResponse(birthdayInWeek)

	if len(birthdayInWeek) > 0 {
		for _, other := range others {
			_, err = sch.bot.Send(telebot.ChatID(other.ID), response)
			if err != nil {
				log.Printf("Failed to send to user %d: %v", other.ID, err)
			}

		}
	}
}

func (sch *Scheduler) deleteWishes() {
	id, _ := strconv.Atoi(os.Getenv("ADMIN_ID"))
	users, err := sch.us.FindAllRegistered()
	if err != nil {
		_, sendErr := sch.bot.Send(telebot.ChatID(id), fmt.Sprintf("Ошибка в уведомлениях"))
		if sendErr != nil {
			log.Printf("Error sending daily notifications: %v", sendErr)
		}
	}

	birthdayTomorrow, _ := splitUsersByBirthday(users, 0)

	if len(birthdayTomorrow) > 0 {
		for _, bd := range birthdayTomorrow {
			sch.ws.DeleteAll(bd.ID)
		}
	}
}

func (sch *Scheduler) Stop() {
	sch.s.Stop()
}

func splitUsersByBirthday(users []user.User, daysBefore int) (birthdayTomorrow []user.User, others []user.User) {
	tomorrow := time.Now().AddDate(0, 0, daysBefore)
	tomorrowMonth := tomorrow.Month()
	tomorrowDay := tomorrow.Day()

	for _, user := range users {
		bdMonth := user.Birthdate.Month()
		bdDay := user.Birthdate.Day()

		if bdMonth == tomorrowMonth && bdDay == tomorrowDay {
			birthdayTomorrow = append(birthdayTomorrow, user)
		} else {
			others = append(others, user)
		}
	}

	return birthdayTomorrow, others
}

func makeResponse(users []user.User) string {
	var response strings.Builder
	now := time.Now().AddDate(0, 0, 1)
	response.WriteString(fmt.Sprintf("Доброе утро!\nЗавтра (%s)  день рождения у:\n\n", now.Format("02.01.2006")))
	for _, user := range users {
		response.WriteString(fmt.Sprintf("(%s) %s %s\n\n", user.Username, user.Surname, user.Name))
	}
	return response.String()
}

func makeWeeklyResponse(users []user.User) string {
	var response strings.Builder
	response.WriteString("Доброе утро!\nЧерез неделю будет день рождения у:\n\n")
	for _, user := range users {
		response.WriteString(fmt.Sprintf("(%s) %s %s.\nДень рождения: %s\n\n", user.Username, user.Surname, user.Name, user.Birthdate.Format("02.01.2006")))
	}
	return response.String()
}
