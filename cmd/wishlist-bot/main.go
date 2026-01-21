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

	go app.Start(cfg, log)

	stop := make(chan os.Signal, 1)

	signal.Notify(stop, os.Interrupt)

	<-stop

	app.Stop()
}
