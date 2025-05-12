package scheduler

import (
	"fmt"
	"github.com/go-co-op/gocron"
	"gopkg.in/telebot.v4"
	"log"
	"os"
	"strings"
	"time"
	constants "wishlist-bot/constant"
	sv "wishlist-bot/service"
)

func StartScheduler(bot *telebot.Bot, userService sv.UserService) {
	location, _ := time.LoadLocation(constants.LOCATION)
	s := gocron.NewScheduler(location)
	scheduleTime := os.Getenv("NOTIFICATION_TIME")
	if scheduleTime == "" {
		scheduleTime = "10:00"
	}
	_, err := s.Every(1).Days().At(scheduleTime).Do(sendDailyNotifications, bot, userService)
	if err != nil {
		log.Fatalf("Error scheduling task: %v", err)
	}

	s.StartAsync()
}

func sendDailyNotifications(bot *telebot.Bot, userService sv.UserService) {
	users := userService.FindAllTotal()
	birthdayTomorrow, others := splitUsersByBirthday(users, 1)
	response := makeResponse(birthdayTomorrow)
	if len(birthdayTomorrow) > 0 {
		for _, other := range others {
			if other.Status == constants.REGISTERED {
				_, err := bot.Send(telebot.ChatID(other.ID), response, telebot.ModeMarkdown)
				if err != nil {
					log.Printf("Failed to send to user %d: %v", other.ID, err)
				}
			}
		}
	}
}

func splitUsersByBirthday(users []sv.UserDto, daysBefore int) (birthdayTomorrow []sv.UserDto, others []sv.UserDto) {
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

func makeResponse(users []sv.UserDto) string {
	var response strings.Builder
	response.WriteString("Доброе утро!\nЗавтра день рождения у:\n")
	for _, user := range users {
		response.WriteString(fmt.Sprintf("(%s) %s %s\n\n", user.Username, user.Surname, user.Name))
	}
	return response.String()
}
