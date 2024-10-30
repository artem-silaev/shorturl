package handler

import (
	"net/http"

	"github.com/artem-silaev/shorturl/internal/app/config"
	"github.com/artem-silaev/shorturl/internal/app/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(service service.URLShortener, config config.Config) http.Handler {
	h := NewHandler(service, config)
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/{shortURL}", h.HandleGet)
	r.Post("/", h.HandlePost)
	r.NotFound(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "path not found", http.StatusBadRequest)
	}))

	return r
}
