package handler

import (
	"compress/flate"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/artem-silaev/shorturl/internal/app/config"
	mdlwr "github.com/artem-silaev/shorturl/internal/app/middleware"
	"github.com/artem-silaev/shorturl/internal/app/service"
)

func NewRouter(service service.URLShortener, config config.Config) http.Handler {
	h := NewHandler(service, config)
	r := chi.NewRouter()

	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(flate.DefaultCompression, ContentTypeJSON, ContentTypeText))
	r.Use(mdlwr.WithLogging)
	r.Use(mdlwr.Decompress)

	r.Get("/{shortURL}", h.HandleGet)
	r.Post("/", h.HandlePost)
	r.Post("/api/shorten", h.HandlePostJSON)
	r.NotFound(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "path not found", http.StatusBadRequest)
	}))

	return r
}
