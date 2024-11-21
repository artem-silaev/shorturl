package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/artem-silaev/shorturl/internal/app/config"
	e "github.com/artem-silaev/shorturl/internal/app/errors"
	"github.com/artem-silaev/shorturl/internal/app/service"
	"github.com/artem-silaev/shorturl/internal/app/utils"
	"github.com/go-chi/chi/v5"
)

const (
	ContentType     = "Content-Type"
	ContentTypeText = "text/plain"
	ContentTypeJSON = "application/json"
)

type Handler struct {
	service service.URLShortener
	config  config.Config
}

func NewHandler(service service.URLShortener, config config.Config) *Handler {
	return &Handler{
		service: service,
		config:  config,
	}
}

func (h *Handler) HandleGet(w http.ResponseWriter, r *http.Request) {
	shortURL := chi.URLParam(r, "shortURL")
	longURL, err := h.service.GetOriginalURL(shortURL)

	switch {
	case errors.Is(err, e.ErrInvalid):
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	case errors.Is(err, e.ErrNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	case errors.Is(err, e.ErrInternal) || err != nil:
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Location", longURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (h *Handler) HandlePost(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)

	if r.URL.Path != "/" || r.Body == http.NoBody || err != nil || !utils.IsURL(string(b)) {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	shortURL, err := h.service.ShortenURL(string(b))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set(ContentType, ContentTypeText)
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(h.config.BaseURL + shortURL))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) HandlePostJSON(w http.ResponseWriter, r *http.Request) {
	var req service.ShortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	shortURL, err := h.service.ShortenURL(req.URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := service.ShortenResponse{Result: h.config.BaseURL + shortURL}

	w.Header().Set(ContentType, ContentTypeJSON)
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
