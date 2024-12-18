package handler

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/artem-silaev/shorturl/internal/app/config"
	"github.com/artem-silaev/shorturl/internal/app/mock"
	"github.com/artem-silaev/shorturl/internal/app/service"
	"github.com/artem-silaev/shorturl/internal/app/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandle(t *testing.T) {
	mockService := mock.NewMockShortenerService()
	mockService.On("ShortenURL", "https://ya.ru").Return("shortURL", nil)
	mockService.On("GetOriginalURL", "shortURL").Return("https://ya.ru", nil)
	config := config.DefaultConfig()
	r := NewRouter(mockService, config)

	tests := []struct {
		name   string
		method string
		path   string
		body   io.Reader
		want   int
	}{
		{
			name:   "test 1",
			method: http.MethodDelete,
			path:   "/",
			body:   nil,
			want:   http.StatusMethodNotAllowed,
		},
		{
			name:   "test 2",
			method: http.MethodPatch,
			path:   "/",
			body:   nil,
			want:   http.StatusMethodNotAllowed,
		},
		{
			name:   "test 3",
			method: http.MethodPost,
			path:   "/",
			body:   strings.NewReader(`https://ya.ru`),
			want:   http.StatusCreated,
		},
		{
			name:   "test 4",
			method: http.MethodGet,
			path:   "/shortURL",
			body:   nil,
			want:   http.StatusTemporaryRedirect,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, tt.path, tt.body)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, request)
			result := w.Result()
			defer result.Body.Close()
			assert.Equal(t, tt.want, result.StatusCode)
			if result.StatusCode == http.StatusTemporaryRedirect {
				location := result.Header.Get("Location")
				assert.Equal(t, location, "https://ya.ru")
			}

		})

	}

}

func TestHandlePost(t *testing.T) {
	config := config.DefaultConfig()
	service := service.NewShortenerService(config.FileStoragePath)
	h := NewHandler(service, config)
	type want struct {
		status int
		isURL  bool
	}
	tests := []struct {
		name        string
		path        string
		body        string
		contentType string
		want        want
	}{
		{
			name:        "test 1",
			body:        `https://ya.ru`,
			path:        `/`,
			contentType: ContentTypeText,
			want: want{
				status: http.StatusCreated,
				isURL:  true,
			},
		},
		{
			name:        "test 2",
			body:        ``,
			path:        `/`,
			contentType: ContentTypeText,
			want: want{
				status: http.StatusBadRequest,
				isURL:  false,
			},
		},
		{
			name:        "test 3",
			body:        `//ya.ru`,
			path:        `/`,
			contentType: ContentTypeText,
			want: want{
				status: http.StatusBadRequest,
				isURL:  false,
			},
		},
		{
			name:        "test 4",
			body:        `https://ya.ru`,
			path:        `/a/b/c`,
			contentType: ContentTypeText,
			want: want{
				status: http.StatusBadRequest,
				isURL:  false,
			},
		},
		{
			name:        "test 5",
			body:        `https://ya.ru`,
			path:        `//`,
			contentType: ContentTypeText,
			want: want{
				status: http.StatusBadRequest,
				isURL:  false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, tt.path, strings.NewReader(tt.body))
			request.Header.Add(ContentType, tt.contentType)

			w := httptest.NewRecorder()
			h := http.HandlerFunc(h.HandlePost)
			h(w, request)

			result := w.Result()
			shortURL, err := io.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)

			assert.Equal(t, tt.want.status, result.StatusCode)
			assert.Equal(t, tt.want.isURL, utils.IsURL(string(shortURL)))
		})
	}

}

func TestHandleGet(t *testing.T) {
	config := config.DefaultConfig()
	service := service.NewShortenerService(config.FileStoragePath)
	longURL := `https://ya.ru`
	serverAddr := `http://localhost:8080/`
	shortURL, _ := service.ShortenURL(longURL)
	r := NewRouter(service, config)

	fmt.Println("=====>", serverAddr+shortURL)

	type want struct {
		status   int
		location string
	}
	tests := []struct {
		name string
		path string
		want want
	}{
		{
			name: "test 1",
			path: serverAddr + shortURL,
			want: want{
				status:   http.StatusTemporaryRedirect,
				location: longURL,
			},
		},
		{
			name: "test 2",
			path: serverAddr + "qwerty",
			want: want{
				status:   http.StatusNotFound,
				location: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, tt.path, nil)

			w := httptest.NewRecorder()
			r.ServeHTTP(w, request)

			result := w.Result()
			defer result.Body.Close()
			assert.Equal(t, tt.want.status, result.StatusCode)
			if result.StatusCode == http.StatusTemporaryRedirect {
				assert.Equal(t, tt.want.location, result.Header.Get("Location"))
			}
		})
	}
}

func TestHandler_HandlePostJSON(t *testing.T) {
	config := config.DefaultConfig()
	service := service.NewShortenerService(config.FileStoragePath)
	h := NewHandler(service, config)
	type want struct {
		status int
		isURL  bool
	}
	tests := []struct {
		name        string
		path        string
		body        string
		contentType string
		want        want
	}{
		{
			name:        "test 1",
			body:        `{"url": "https://practicum.yandex.ru"} `,
			contentType: ContentTypeJSON,
			want: want{
				status: http.StatusCreated,
			},
		},
	}

	serverAddr := `http://localhost:8080/`
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, serverAddr+"api/shorten", strings.NewReader(tt.body))
			request.Header.Add(ContentType, tt.contentType)

			w := httptest.NewRecorder()
			h := http.HandlerFunc(h.HandlePostJSON)
			h(w, request)

			result := w.Result()
			_, err := io.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)

			assert.Equal(t, tt.want.status, result.StatusCode)
		})
	}
}
