package config

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env    string `yaml:"env" env:"ENV" env-default:"local"`
	DB     `yaml:"db"`
	Server `yaml:"server"`
}

type DB struct {
	User     string `yaml:"user"     env:"POSTGRES_USER"     env-default:"todo_user"`
	Password string `yaml:"password" env:"POSTGRES_PASSWORD" env-default:"todo_password"`
	Name     string `yaml:"name"     env:"POSTGRES_DB"       env-default:"todo_db"`
	Host     string `yaml:"host"     env:"POSTGRES_HOST"     env-default:"localhost"`
	Port     int    `yaml:"port"     env:"POSTGRES_PORT"     env-default:"5432"`
	SSLMode  string `yaml:"sslmode"  env:"POSTGRES_SSLMODE"  env-default:"disable"`
}

type Server struct {
	Port    int    `yaml:"port" env:"SERVER_PORT" env-default:"8080"`
	Address string `yaml:"address" env:"ADDRESS" env-default:"localhost"`
}

func MustLoad(path string) Config {
	if path == "" {
		log.Fatal("config path is empty")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", path)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		log.Fatalf("error reading config: %s", path)
	}

	return cfg
}
