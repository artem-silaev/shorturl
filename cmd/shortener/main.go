package main

import (
	"encoding/base64"
	"fmt"
	"github.com/go-chi/chi/v5"
	"io"
	"log"
	"net/http"
	"sync"
)

var (
	urlStore = make(map[string]string)
	mu       sync.Mutex
)

func shortenURL(url string) string {
	return base64.RawURLEncoding.EncodeToString([]byte(url))
}

func createShortURLHandler(w http.ResponseWriter, r *http.Request) {
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
	shortURL := fmt.Sprintf("http://localhost:8080/%s", id)

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

func ShortUrlRouter() chi.Router {
	r := chi.NewRouter()

	r.Post("/", createShortURLHandler)
	r.Get("/{url}", redirectHandler)
	return r
}

func main() {
	log.Fatal(http.ListenAndServe(":8080", ShortUrlRouter()))
}
