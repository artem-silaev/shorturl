package config

import (
	"flag"
	"os"
	"strings"
)

// Config хранит параметры конфигурации
type Config struct {
	Address         string // адрес запуска HTTP-сервера
	BaseURL         string // базовый адрес результирующего сокращённого URL
	FileStoragePath string
	DatabaseDSN     string
}

// InitConfig инициализирует конфигурацию из аргументов командной строки
func InitConfig() Config {
	address := flag.String("a", "localhost:8080", "Address for the HTTP server")
	baseURL := flag.String("b", "http://localhost:8080", "Base URL for the shortened URL")
	fileStoragePath := flag.String("f", "test.json", "File for JSON data")
	databaseDSN := flag.String("d", "", "Database DSN")

	flag.Parse()

	if os.Getenv(`SERVER_ADDRESS`) != "" {
		*address = os.Getenv(`SERVER_ADDRESS`)
	}

	if os.Getenv("BASE_URL") != "" {
		*baseURL = os.Getenv(`BASE_URL`)
	}

	if os.Getenv("FILE_STORAGE_PATH") != "" {
		*fileStoragePath = os.Getenv(`FILE_STORAGE_PATH`)
	}

	if os.Getenv("DATABASE_DSN") != "" {
		*databaseDSN = os.Getenv(`DATABASE_DSN`)
	}

	if !strings.HasSuffix(*baseURL, "/") {
		*baseURL += "/"
	}

	config := Config{
		Address:         *address,
		BaseURL:         *baseURL,
		FileStoragePath: *fileStoragePath,
		DatabaseDSN:     *databaseDSN,
	}

	return config
}

func DefaultConfig() Config {
	return Config{
		Address:         `localhost:8080`,
		BaseURL:         `http://localhost:8080/`,
		FileStoragePath: `default.json`,
		DatabaseDSN:     ``,
	}
}
