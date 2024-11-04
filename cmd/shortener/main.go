package main

import (
	"github.com/artem-silaev/shorturl/internal/app/config"
	"github.com/artem-silaev/shorturl/internal/app/handler"
	"github.com/artem-silaev/shorturl/internal/app/logger"
	_ "github.com/artem-silaev/shorturl/internal/app/middleware"
	"github.com/artem-silaev/shorturl/internal/app/service"
	"github.com/artem-silaev/shorturl/internal/app/storage"
	"log"
	"net/http"
)

func main() {
	cfg := config.InitConfig()
	logger.Init()
	service := service.NewShortenerService(cfg.FileStoragePath)
	storage := storage.NewStorage(cfg.FileStoragePath)
	storage.LoadURLs()
	r := handler.NewRouter(service, cfg)
	err := http.ListenAndServe(cfg.Address, r)
	if err != nil {
		log.Fatal(err)
	}
}
