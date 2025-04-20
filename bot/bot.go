package bot

import (
	"github.com/joho/godotenv"
	"gopkg.in/telebot.v4"
	"log"
	"os"
	"time"
	sv "wishlist-bot/service"
)

func NewBot() (*telebot.Bot, error) {
	err := godotenv.Load(".env")
	pref := telebot.Settings{
		Token:  os.Getenv("API_KEY"),
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}
	bot, err := telebot.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return bot, nil
}

func Start(bot *telebot.Bot, service sv.UserService) {
	setUpButtons()
	setUpHandlers(bot, service)
	bot.Start()
}
