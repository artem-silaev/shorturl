package main

import (
	"github.com/artem-silaev/shorturl/internal/app/config"
	"github.com/artem-silaev/shorturl/internal/app/handler"
	"github.com/artem-silaev/shorturl/internal/app/service"
	"log"
	"net/http"
)

func main() {
	cfg := config.InitConfig()
	service := service.NewShortenerService()
	r := handler.NewRouter(service, cfg)
	err := http.ListenAndServe(cfg.Address, r)
	if err != nil {
		log.Fatal(err)
	}
}
