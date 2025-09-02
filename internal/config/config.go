package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"time"
)

type Config struct {
	DB       DatabaseConfig `yaml:"database"`
	HTTP     HTTPConfig     `yaml:"http"`
	Telegram TelegramConfig `yaml:"telegram"`
}

type DatabaseConfig struct {
	ServerAddress      string `yaml:"server_address"`
	Username           string `yaml:"username"`
	Password           string `yaml:"password"`
	Host               string `yaml:"host"`
	Port               string `yaml:"port"`
	Database           string `yaml:"name"`
	MaxOpenConnections int    `yaml:"max_open_conns"`
	MaxIdleConnections int    `yaml:"max_idle_conns"`
	ConnMaxLifetime    int    `yaml:"conn_max_lifetime"`
}

type HTTPConfig struct {
	Port         string        `yaml:"port"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
}

type TelegramConfig struct {
	BotToken string `yaml:"bot_token"`
	ChatID   string `yaml:"chat_id"`
}

func Load(path string) (*Config, error) {
	var cfg Config

	err := cleanenv.ReadConfig(path, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
