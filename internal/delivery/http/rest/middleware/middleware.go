package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/delivery/grpc/service/session"

	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

const (
	ProfileID = 0
	SessionID = 1
)

var allowedOrigins = map[string]interface{}{
	"http://127.0.0.1:5000":     struct{}{},
	"http://127.0.0.1:8080":     struct{}{},
	"http://localhost:5000":     struct{}{},
	"http://localhost:8080":     struct{}{},
	"http://104.248.28.45":      struct{}{},
	"https://104.248.28.45":     struct{}{},
	"https://sadislands.now.sh": struct{}{},
	"http://sadislands.ru":      struct{}{},
	"https://sadislands.ru":     struct{}{},
}

// Authentication middleware to check session.
// If cookie wasn't transmitted, expired or doesn't exist,
// returns status code unauthorized.
// Otherwise process request.
func Authentication(sessionManager *session.SessionClient, next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		ctx := r.Context()
		cookie, err := r.Cookie("session_id")
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		sessionId := &session.SessionId{
			Id: cookie.Value,
		}

		profileId, err := (*sessionManager).Get(context.Background(), sessionId)

		if err != nil {
			cookie := http.Cookie{
				Name:     "session_id",
				Path:     "/",
				MaxAge:   -1,
				HttpOnly: true,
			}
			http.SetCookie(w, &cookie)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx = context.WithValue(ctx, ProfileID, profileId.Id)
		ctx = context.WithValue(ctx, SessionID, cookie.Value)
		next(w, r.WithContext(ctx), ps)
	}
}

// CORSMiddleware implements CORS mechanism
func CORSMiddleware(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

		next(w, r, ps)
	}
}

// RecoverMiddleware catches panics
func RecoverMiddleware(logger *zap.Logger, next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error(err.(error).Error())
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()
		next(w, r, ps)
	}
}

// LoggerMiddleware write logs
func LoggerMiddleware(logger *zap.Logger, next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		start := time.Now()
		log := logger.With(zap.String("method", r.Method),
			zap.String("remote_address", r.RemoteAddr),
			zap.String("url", r.URL.Path))
		ctx := r.Context()
		ctx = context.WithValue(ctx, "logger", log)
		next(w, r.WithContext(ctx), ps)
		log.Info("Finish", zap.Duration("work_time", time.Since(start)))
	}
}
