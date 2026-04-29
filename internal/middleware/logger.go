package middleware

import (
	"log"
	"net/http"
	"time"
)

func RequestLogger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next(w, r)
		duration := time.Since(start)
		log.Printf("[%s] %s - Süre: %v\n", r.Method, r.URL.Path, duration)

	}
}
