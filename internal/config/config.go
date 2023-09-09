package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Endpoints   []string      `yaml:"endpoints" env-required:"true"`
	Env         string        `yaml:"env" env-default:"local" env-required:"true"`
	StoragePath string        `yaml:"storage_path" env-required:"true"`
	DialTimeout time.Duration `yaml:"dial_timeout" env-default:"5s"`
	HTTPServer  `yaml:"http_server"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
	// User        string        `yaml:"user" env-required:"true"`
	// Password    string        `yaml:"password" env-required:"true" env:"HTTP_SERVER_PASSWORD"`
}

//Приставка  Must используется  если ф-ия возвращает панику а не ошибку(пример инициализация конфига при запуске приложении)

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	// если нет файла с конфигом то падаем с паникой
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}
