package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"sync"
)

var (
	urlStore = make(map[string]string)
	mu       sync.Mutex
)

func shortenURL(url string) string {
	// Кодируем URL в base64 для создания уникального идентификатора
	return base64.RawURLEncoding.EncodeToString([]byte(url))
}

func createShortURLHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		//http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		redirectHandler(w, r)
	}

	body, err := io.ReadAll(r.Body)
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
	w.Header().Set("Content-Type", "text/plain")
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

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc(`/`, createShortURLHandler)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
