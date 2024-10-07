package main

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_createShortURLHandler(t *testing.T) {
	type want struct {
		code        int
		request     string
		response    string
		contentType string
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "positive test #1",
			want: want{
				code:        201,
				request:     `http://ya.ru`,
				response:    `http://localhost:8080/aHR0cDovL3lhLnJ1`,
				contentType: "text/plain",
			},
		},
		{
			name: "negative test #1",
			want: want{
				code:        400,
				request:     ``,
				response:    "Bad request\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(test.want.request))

			w := httptest.NewRecorder()
			createShortURLAction(w, request, `http://localhost:8080`)
			res := w.Result()
			assert.Equal(t, test.want.code, res.StatusCode)
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.Equal(t, test.want.response, string(resBody))
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}

func Test_shortenURL(t *testing.T) {
	type want struct {
		want  string
		value string
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "positive test #1",
			want: want{
				value: `http://pip111i.ru`,
				want:  `aHR0cDovL3BpcDExMWkucnU`,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require.Equal(t, test.want.want, shortenURL(test.want.value))
		})
	}
}

func Test_redirectHandler(t *testing.T) {
	type want struct {
		code        int
		request     string
		response    string
		contentType string
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "positive test #1",
			want: want{
				code:        http.StatusTemporaryRedirect,
				request:     `http://localhost:8080/aHR0cDovL3lhLnJ1`,
				response:    `http://ya.ru`,
				contentType: "text/plain",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			shortURL := shortenURL(test.want.response)

			mu.Lock()
			urlStore[shortURL] = test.want.response
			mu.Unlock()

			req := httptest.NewRequest(http.MethodGet, "/"+shortURL, nil)
			w := httptest.NewRecorder()

			redirectHandler(w, req)

			resp := w.Result()
			defer resp.Body.Close()
			require.Equal(t, test.want.code, resp.StatusCode)

			location := resp.Header.Get("Location")
			require.Equal(t, test.want.response, location)
		})
	}
}
