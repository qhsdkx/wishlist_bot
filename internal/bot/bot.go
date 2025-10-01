package bot

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gopkg.in/telebot.v4"
)

type Bot struct {
	tg     *telebot.Bot
	router HandlerRouter
}

func NewBot(router HandlerRouter) (*Bot, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, errors.New("something went wrong with .env file")
	}
	pref := telebot.Settings{
		Token:  os.Getenv("API_KEY"),
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}
	bot, err := telebot.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return &Bot{
		tg:     bot,
		router: router,
	}, nil
}

func (b *Bot) RegisterHandlers() {
	b.tg.Handle("/start", b.router.OnStart)
	b.tg.Handle(telebot.OnText, b.router.OnText)
	b.tg.Handle(telebot.OnCallback, b.router.OnCallback)
	b.tg.Handle("/help", b.router.Help)
}

func (b *Bot) Start() {
	b.tg.Start()
}
