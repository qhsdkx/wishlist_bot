package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
)

type Config struct {
	AdminId   int64           `yaml:"admin-id"`
	Database  DatabaseConfig  `yaml:"database"`
	Bot       BotConfig       `yaml:"bot"`
	Scheduler SchedulerConfig `yaml:"scheduler"`
}

type DatabaseConfig struct {
	Host       string `yaml:"host"`
	Port       string `yaml:"port"`
	User       string `yaml:"user"`
	Password   string `yaml:"password"`
	Name       string `yaml:"name"`
	MaxDBConns int    `yaml:"max-db-connections"`
}

type BotConfig struct {
	ApiKey string `yaml:"api-key"`
}

type SchedulerConfig struct {
	NotifTimeDaily  string `yaml:"notification-time-daily"`
	NotifTimeWeekly string `yaml:"notification-time-weekly"`
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
