package scheduler

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
	constants "wishlist-bot/constant"
	sv "wishlist-bot/service"

	"github.com/go-co-op/gocron"
	"gopkg.in/telebot.v4"
)

func StartScheduler(bot *telebot.Bot, userService sv.UserService) {
	location, _ := time.LoadLocation(constants.LOCATION)
	s := gocron.NewScheduler(location)

	scheduleTime := os.Getenv("NOTIFICATION_TIME")
	if scheduleTime == "" {
		scheduleTime = "10:00"
	}

	scheduleTimeWeekly := os.Getenv("NOTIFICATION_TIME_WEEKLY")
	if scheduleTimeWeekly == "" {
		scheduleTimeWeekly = "9:55"
	}

	_, err := s.Every(1).Days().At(scheduleTime).Do(sendDailyNotifications, bot, userService)
	if err != nil {
		log.Fatalf("Error scheduling task: %v", err)
	}

	_, err = s.Every(1).Monday().At(scheduleTimeWeekly).Do(sendWeeklyNotifications, bot, userService)
	if err != nil {
		log.Fatalf("Error scheduling task: %v", err)
	}

	s.StartAsync()
}

func sendDailyNotifications(bot *telebot.Bot, userService sv.UserService) {
	id, _ := strconv.Atoi(os.Getenv("ADMIN_ID"))
	users, err := userService.FindAllRegistered()
	if err != nil {
		_, sendErr := bot.Send(telebot.ChatID(id), fmt.Sprintf("Ошибка в уведомлениях"))
		if sendErr != nil {
			log.Printf("Error sending daily notifications: %v", sendErr)
		}
	}

	birthdayTomorrow, others := splitUsersByBirthday(users, 1)
	response := makeResponse(birthdayTomorrow)

	if len(birthdayTomorrow) > 0 {
		for _, other := range others {
			_, err := bot.Send(telebot.ChatID(other.ID), response)
			if err != nil {
				log.Printf("Failed to send to user %d: %v", other.ID, err)
			}
		}
	}
}

func sendWeeklyNotifications(bot *telebot.Bot, userService sv.UserService) {
	id, _ := strconv.Atoi(os.Getenv("ADMIN_ID"))
	users, err := userService.FindAllRegistered()
	if err != nil {
		_, sendErr := bot.Send(telebot.ChatID(id), fmt.Sprintf("Ошибка в уведомлениях"))
		if sendErr != nil {
			log.Printf("Error sending daily notifications: %v", sendErr)
		}
	}

	birthdayInWeek, others := splitWeeklyUsersByBirthday(users)
	response := makeWeeklyResponse(birthdayInWeek)

	if len(birthdayInWeek) > 0 {
		for _, other := range others {
			_, err := bot.Send(telebot.ChatID(other.ID), response)
			if err != nil {
				log.Printf("Failed to send to user %d: %v", other.ID, err)
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
	now := time.Now().AddDate(0, 0, 1)
	response.WriteString(fmt.Sprintf("Доброе утро!\nЗавтра (%s)  день рождения у:\n\n", now.Format("02.01.2006")))
	for _, user := range users {
		response.WriteString(fmt.Sprintf("(%s) %s %s\n\n", user.Username, user.Surname, user.Name))
	}
	return response.String()
}

func splitWeeklyUsersByBirthday(users []sv.UserDto) (birthdayInWeek []sv.UserDto, others []sv.UserDto) {
	today := time.Now()
	weekLater := today.AddDate(0, 0, 7)

	for _, user := range users {
		bd := user.Birthdate
		currentYearBD := time.Date(today.Year(), bd.Month(), bd.Day(), 0, 0, 0, 0, time.UTC)

		if currentYearBD.Before(today) {
			currentYearBD = currentYearBD.AddDate(1, 0, 0)
		}

		if !currentYearBD.After(weekLater) {
			birthdayInWeek = append(birthdayInWeek, user)
		} else {
			others = append(others, user)
		}
	}

	return birthdayInWeek, others
}

func makeWeeklyResponse(users []sv.UserDto) string {
	var response strings.Builder
	response.WriteString("Доброе утро!\nЧерез неделю будет день рождения у:\n\n")
	for _, user := range users {
		response.WriteString(fmt.Sprintf("(%s) %s %s.\nДень рождения: %s\n\n", user.Username, user.Surname, user.Name, user.Birthdate.Format("02.01.2006")))
	}
	return response.String()
}
