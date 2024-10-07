package main

import (
	"encoding/base64"
	"fmt"
	"github.com/artem-silaev/shorturl/cmd/shortener/config"
	"github.com/go-chi/chi/v5"
	"io"
	"log"
	"net/http"
	"sync"
)

var (
	urlStore = make(map[string]string)
	cfg      *config.Config
	mu       sync.Mutex
)

func shortenURL(url string) string {
	return base64.RawURLEncoding.EncodeToString([]byte(url))
}

func createShortURLHandler(w http.ResponseWriter, r *http.Request) {
	createShortUrlAction(w, r, cfg.BaseURL)
}

func createShortUrlAction(w http.ResponseWriter, r *http.Request, baseURL string) {
	body, err := io.ReadAll(r.Body)

	w.Header().Set("Content-Type", "text/plain")
	if err != nil || len(body) == 0 {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	originalURL := string(body)

	id := shortenURL(originalURL)
	urlStore[id] = originalURL

	// Формируем сокращённый URL
	shortURL := fmt.Sprintf("%s/%s", baseURL, id)

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(shortURL))
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[1:]

	mu.Lock()
	originalURL, exists := urlStore[id]
	mu.Unlock()

	if !exists {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	w.Header().Set("Location", originalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func ShortURLRouter() chi.Router {
	r := chi.NewRouter()

	r.Post("/", createShortURLHandler)
	r.Get("/{url}", redirectHandler)
	return r
}

func main() {
	cfg = config.InitConfig()
	log.Fatal(http.ListenAndServe(cfg.Address, ShortURLRouter()))
}
