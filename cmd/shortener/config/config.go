package config

import (
	"flag"
	"os"
)

// Config хранит параметры конфигурации
type Config struct {
	Address string // адрес запуска HTTP-сервера
	BaseURL string // базовый адрес результирующего сокращённого URL
}

// InitConfig инициализирует конфигурацию из аргументов командной строки
func InitConfig() *Config {
	address := flag.String("a", "localhost:8080", "Address for the HTTP server")
	baseURL := flag.String("b", "http://localhost:8080", "Base URL for the shortened URL")

	flag.Parse()

	if os.Getenv(`SERVER_ADDRESS`) != "" {
		*address = os.Getenv(`SERVER_ADDRESS`)
	}

	if os.Getenv("BASE_URL") != "" {
		*baseURL = os.Getenv(`BASE_URL`)
	}

	config := &Config{
		Address: *address,
		BaseURL: *baseURL,
	}

	return config
}
