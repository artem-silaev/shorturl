package middleware

import (
	"compress/gzip"
	"compress/zlib"
	"io"
	"net/http"
	"time"

	"github.com/artem-silaev/shorturl/internal/app/errors"
	"github.com/artem-silaev/shorturl/internal/app/logger"
)

type (
	responseData struct {
		status int
		size   int
	}

	loggingResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}
)

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

func WithLogging(h http.Handler) http.Handler {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		responseData := &responseData{
			status: 0,
			size:   0,
		}
		lw := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   responseData,
		}
		h.ServeHTTP(&lw, r)

		duration := time.Since(start)

		logger.Log.Info(
			"uri", r.RequestURI,
			"method", r.Method,
			"status", responseData.status, // получаем перехваченный код статуса ответа
			"duration", duration,
			"size", responseData.size, // получаем перехваченный размер ответа
		)
	}
	return http.HandlerFunc(logFn)
}

func Decompress(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var reader io.ReadCloser
		var err error

		encoding := r.Header.Get("Content-Encoding")
		switch encoding {
		case "gzip":
			reader, err = gzip.NewReader(r.Body)
		case "deflate":
			reader, err = zlib.NewReader(r.Body)
		case "":
			reader, err = r.Body, nil
		}

		if err != nil {
			logger.Log.Error(err.Error())
			http.Error(w, errors.ErrDecompress.Error(), http.StatusBadRequest)
			return
		}

		if reader == nil {
			logger.Log.Error("decompression not implemented")
			http.Error(w, "decompression not implemented", http.StatusInternalServerError)
			return
		}

		defer reader.Close()
		r.Body = reader

		next.ServeHTTP(w, r)

	})
}
