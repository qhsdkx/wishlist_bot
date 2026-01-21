package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
)

type Config struct {
	adminId   int64           `yaml:"admin-id"`
	database  DatabaseConfig  `yaml:"database"`
	bot       BotConfig       `yaml:"bot"`
	scheduler SchedulerConfig `yaml:"scheduler"`
}

type DatabaseConfig struct {
	host       string `yaml:"host"`
	port       string `yaml:"port"`
	user       string `yaml:"user"`
	password   string `yaml:"password"`
	name       string `yaml:"name"`
	maxDBConns int    `yaml:"max-db-connections"`
}

type BotConfig struct {
	apiKey string `yaml:"api-key"`
}

type SchedulerConfig struct {
	notifTimeDaily  string `yaml:"notification-time-daily"`
	notifTimeWeekly string `yaml:"notification-time-weekly"`
}

func MustLoad() *Config {
	configPath := fetchConfigPath()
	if configPath == "" {
		panic("config path is empty")
	}

	return MustLoadPath(configPath)
}

func MustLoadPath(configPath string) *Config {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exist: " + configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("cannot read config: " + err.Error())
	}

	return &cfg
}

func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
