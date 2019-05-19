package middleware

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/pkg/metric"

	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/delivery/grpc/service/session"

	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

const (
	ProfileID = 0
	SessionID = 1
)

var allowedOrigins = map[string]bool{
	"http://127.0.0.1:5000":           true,
	"http://127.0.0.1:8080":           true,
	"http://localhost:5000":           true,
	"http://localhost:8080":           true,
	"http://104.248.28.45":            true,
	"https://104.248.28.45":           true,
	"https://sadislands.now.sh":       true,
	"http://sadislands.ru":            true,
	"https://sadislands.ru":           true,
	"https://hackathon.sadislands.ru": true,
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
				Domain:   ".sadislands.ru",
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

// Session middleware passes session and profile id if they exists
func SessionMiddleware(sessionManager *session.SessionClient, next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		ctx := r.Context()
		cookie, err := r.Cookie("session_id")
		if err != nil {
			next(w, r.WithContext(ctx), ps)
			return
		}

		sessionId := &session.SessionId{
			Id: cookie.Value,
		}

		profileId, err := (*sessionManager).Get(context.Background(), sessionId)

		if err != nil {
			next(w, r.WithContext(ctx), ps)
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
		domain := r.Header.Get("Origin")
		if allowedOrigins[domain] {
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Origin", domain)
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
		//start := time.Now()
		log := logger.With(zap.String("method", r.Method),
			zap.String("remote_address", r.RemoteAddr),
			zap.String("url", r.URL.Path))
		ctx := r.Context()
		ctx = context.WithValue(ctx, "logger", log)
		next(w, r.WithContext(ctx), ps)
		//log.Info("Finish",
		//	zap.Duration("work_time", time.Since(start)),
		//)
	}
}

func AccessHitsMiddleware(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		start := time.Now()
		log := r.Context().Value("logger").(*zap.Logger)
		lrw := NewLoggingResponseWriter(w)
		next(lrw, r, ps)

		// Write hits metric
		if metric.AccessHits != nil {
			metric.AccessHits.With(prometheus.Labels{
				"path":        r.URL.Path,
				"method":      r.Method,
				"status_code": strconv.Itoa(lrw.statusCode),
			}).Inc()
		}

		log.Info("Finish with status code",
			zap.Int("status_code", lrw.statusCode),
			zap.Duration("work_time", time.Since(start)))
	}
}

func OauthMiddleware(clientId, clientSecret string, next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		ctx := r.Context()
		ctx = context.WithValue(ctx, "clientId", clientId)
		ctx = context.WithValue(ctx, "clientSecret", clientSecret)
		next(w, r.WithContext(ctx), ps)
	}
}
