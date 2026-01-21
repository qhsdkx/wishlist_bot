package bot

import (
	"log/slog"
	"wishlist-bot/internal/config"
)

type App struct {
	Bot *Bot
	cfg *config.Config
	log *slog.Logger
}
