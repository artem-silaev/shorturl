package main

import (
	"log"
	"net/http"

	"github.com/artem-silaev/shorturl/internal/app/config"
	"github.com/artem-silaev/shorturl/internal/app/handler"
	"github.com/artem-silaev/shorturl/internal/app/logger"
	_ "github.com/artem-silaev/shorturl/internal/app/middleware"
	"github.com/artem-silaev/shorturl/internal/app/service"
)

func main() {
	cfg := config.InitConfig()
	logger.Init()
	shortenerService := service.NewShortenerService(cfg.FileStoragePath)
	r := handler.NewRouter(shortenerService, cfg)
	err := http.ListenAndServe(cfg.Address, r)
	if err != nil {
		log.Fatal(err)
	}
}
