package main

import (
	"os"
	"os/signal"
	"wishlist-bot/internal/app"
	"wishlist-bot/internal/config"
	"wishlist-bot/internal/logger"
)

func main() {
	cfg := config.MustLoad()

	log := logger.InitializeLogger()

	mainApp := app.New(cfg, log)

	go mainApp.MustStart()

	stop := make(chan os.Signal, 1)

	signal.Notify(stop, os.Interrupt)

	<-stop

	mainApp.Stop()
}
