package bot

import (
	"errors"
	"log"
	"os"
	"time"
	sv "wishlist-bot/service"

	"github.com/joho/godotenv"
	"gopkg.in/telebot.v4"
)

var bot *telebot.Bot

func newBot() (*telebot.Bot, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, errors.New("something went wrong with .env file")
	}
	pref := telebot.Settings{
		Token:  os.Getenv("API_KEY"),
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}
	bot, err = telebot.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return bot, nil
}

func SetUp(userService sv.UserService, wishlistService sv.WishService) *telebot.Bot {
	setUpButtons()
	bot, err := newBot()
	if err != nil {
		panic("Someting went wrong with bot connection\nMaybe u need check .env file or another settings")
	}
	setUpHandlers(bot, userService, wishlistService)
	return bot
}
