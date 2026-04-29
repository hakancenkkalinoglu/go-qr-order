package middleware

import (
	"go-qr-order/internal/utils"
	"net/http"
	"strings"
)

func RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "Unauthorized: Token not found!", http.StatusUnauthorized)
			return
		}

		if !strings.HasPrefix(token, "Bearer") {
			http.Error(w, "Unauthorized: Invalid Token format!", http.StatusUnauthorized)
			return
		}

		cleanToken := strings.TrimPrefix(token, "Bearer ")

		_, err := utils.ValidateToken(cleanToken)
		if err != nil {
			http.Error(w, "Unauthorized: Token invalid or expired!", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}
