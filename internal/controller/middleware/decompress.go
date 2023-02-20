package middleware

import (
	"compress/gzip"
	"log"
	"net/http"
)

const GZip = "gzip"

// Decompress .-
func Decompress(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Encoding") == GZip {
			gz, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			r.Body = gz
			defer func(gz *gzip.Reader) {
				err = gz.Close()
				if err != nil {
					log.Printf("close failed: %v", err)
				}
			}(gz)
		}
		next.ServeHTTP(w, r)
	})
}
