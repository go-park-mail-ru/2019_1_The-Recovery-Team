package middleware

import (
	"api/models"
	"context"
	"net/http"
	"time"
)

type key int

const (
	ProfileID key = iota
	SessionID
)

// MiddlewareWithEnv middleleware with env
type MiddlewareWithEnv func(*models.Env, http.HandlerFunc) http.HandlerFunc

// Authentication middleware to check authentication
func Authentication(env *models.Env, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		cookie, err := r.Cookie("session_id")
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		id, err := env.Sm.Get(cookie.Value)
		if err != nil {
			cookie.Expires = time.Now().AddDate(0, 0, -1)
			http.SetCookie(w, cookie)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx = context.WithValue(ctx, ProfileID, id)
		ctx = context.WithValue(ctx, SessionID, cookie.Value)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}