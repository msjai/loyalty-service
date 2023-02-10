package middleware

import (
	"compress/gzip"
	"context"
	"io"
	"net/http"
)

type ReaderContextKey string

const GZip = "gzip"

// Decompress .-
func Decompress(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var reader io.Reader
		if r.Header.Get("Content-Encoding") == GZip {
			gz, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			reader = gz
			defer gz.Close()
		} else {
			reader = r.Body
		}

		k := ReaderContextKey("reader")
		ctx := context.WithValue(r.Context(), k, reader)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
