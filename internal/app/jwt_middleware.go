package app

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"strings"
)

var (
	ErrInvalidHeaderFormat = errors.New("invalid authorization header format")
	ErrInvalidToken        = errors.New("invalid token")
	ErrMissingAuthHeader   = errors.New("authorization header required")
)

func (a *App) ApiTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			respondWithErr(w, ErrMissingAuthHeader, http.StatusUnauthorized)
			return
		}

		tokenString := strings.Split(authHeader, " ")
		if len(tokenString) != 2 || tokenString[0] != "Bearer" {
			respondWithErr(w, ErrInvalidHeaderFormat, http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(tokenString[1], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return a.config.AuthJwtSecret(), nil
		})

		if err != nil || !token.Valid {
			respondWithErr(w, ErrInvalidToken, http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
