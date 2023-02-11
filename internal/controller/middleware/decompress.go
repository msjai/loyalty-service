package middleware

import (
	"compress/gzip"
	"context"
	"io"
	"net/http"
)

type KeyForReader string

var KeyReader KeyForReader = "reader"

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

		// var k KeyForReader
		// k = "reader"
		ctx := context.WithValue(r.Context(), KeyReader, reader)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
