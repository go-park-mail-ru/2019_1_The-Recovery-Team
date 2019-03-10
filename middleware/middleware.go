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

var allowedOrigins = map[string]interface{}{
	"http://127.0.0.1:5000":     struct{}{},
	"http://127.0.0.1:8080":     struct{}{},
	"http://localhost:5000":     struct{}{},
	"http://localhost:8080":     struct{}{},
	"https://sadislands.now.sh": struct{}{},
	"http://sadislands.ru":      struct{}{},
	"https://sadislands.ru":     struct{}{},
}

// MiddlewareWithEnv middleleware with env
type MiddlewareWithEnv func(*models.Env, http.HandlerFunc) http.HandlerFunc

// Authentication middleware to check authentication
func Authentication(env *models.Env, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
	})
}

// CORSMiddleware CORS middleware
func CORSMiddleware(env *models.Env, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		o := r.Header.Get("Origin")
		if _, ok := allowedOrigins[o]; ok {
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Origin", o)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Max-Age", "86400")
			w.Header().Set("Access-Control-Allow-Headers",
				"Content-Type, User-Agent, Cache-Control, Accept, X-Requested-With, If-Modified-Since, Origin")
		}

		if r.Method == http.MethodOptions {
			return
		}

		next.ServeHTTP(w, r)
	})
}

// func RecoverMiddleware(next http.HandlerFunc) http.HandlerFunc {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		defer func() {
// 			if err := recover(); err != nil {
// 				logger.Error("[PANIC]: ", err, " at ", string(debug.Stack()))
// 				w.WriteHeader(http.StatusInternalServerError)
// 			}
// 		}()
// 		next.ServeHTTP(w, r)
// 	})
// }

// func AccessLogMiddleware(next http.HandlerFunc) http.HandlerFunc {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		start := time.Now()
// 		next.ServeHTTP(w, r)

// 		logger.Infow(r.URL.Path,
// 			"method", r.Method,
// 			"remote_addr", r.RemoteAddr,
// 			"url", r.URL.Path,
// 			"work_time", time.Since(start).String(),
// 		)
// 	})
// }
