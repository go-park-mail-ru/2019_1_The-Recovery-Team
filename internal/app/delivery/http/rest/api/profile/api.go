package profile

import (
	"net/http"

	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/delivery/grpc/service/profile"
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/delivery/grpc/service/session"
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/delivery/http/rest/api/profile/handler"
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/delivery/http/rest/middleware"
	"github.com/julienschmidt/httprouter"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

type Api struct {
	Router *httprouter.Router
}

func NewApi(
	profileManager *profile.ProfileClient,
	sessionManager *session.SessionClient,
	logger *zap.Logger,
	clientId string,
	clienSecret string,
) *Api {
	router := httprouter.New()

	//Profile routes
	router.GET("/api/v1/profiles",
		middleware.LoggerMiddleware(
			logger, middleware.RecoverMiddleware(
				logger, middleware.AccessHitsMiddleware(
					middleware.CORSMiddleware(handler.GetProfiles(profileManager))))),
	)
	router.POST("/api/v1/profiles",
		middleware.LoggerMiddleware(
			logger, middleware.RecoverMiddleware(
				logger, middleware.AccessHitsMiddleware(
					middleware.CORSMiddleware(handler.PostProfile(profileManager, sessionManager))))),
	)
	router.GET("/api/v1/profiles/:id",
		middleware.LoggerMiddleware(
			logger, middleware.RecoverMiddleware(
				logger, middleware.AccessHitsMiddleware(
					middleware.SessionMiddleware(sessionManager,
						middleware.CORSMiddleware(handler.GetProfile(profileManager)))))),
	)
	router.PUT("/api/v1/profiles/:id",
		middleware.LoggerMiddleware(
			logger, middleware.RecoverMiddleware(
				logger, middleware.AccessHitsMiddleware(
					middleware.Authentication(sessionManager,
						middleware.CORSMiddleware(handler.PutProfile(profileManager)))))),
	)
	router.PUT("/api/v1/profiles/:id/password",
		middleware.LoggerMiddleware(
			logger, middleware.RecoverMiddleware(
				logger, middleware.AccessHitsMiddleware(
					middleware.Authentication(sessionManager,
						middleware.CORSMiddleware(handler.PutProfilePassword(profileManager)))))),
	)
	router.PUT("/api/v1/avatars",
		middleware.LoggerMiddleware(
			logger, middleware.RecoverMiddleware(
				logger, middleware.AccessHitsMiddleware(
					middleware.Authentication(sessionManager,
						middleware.CORSMiddleware(handler.PutAvatar(profileManager)))))),
	)

	//Session routes
	router.GET("/api/v1/sessions",
		middleware.LoggerMiddleware(
			logger, middleware.RecoverMiddleware(
				logger, middleware.AccessHitsMiddleware(
					middleware.Authentication(sessionManager,
						middleware.CORSMiddleware(handler.GetSession()))))),
	)
	router.POST("/api/v1/sessions",
		middleware.LoggerMiddleware(
			logger, middleware.RecoverMiddleware(
				logger, middleware.AccessHitsMiddleware(
					middleware.CORSMiddleware(handler.PostSession(profileManager, sessionManager))))),
	)
	router.DELETE("/api/v1/sessions",
		middleware.LoggerMiddleware(
			logger, middleware.RecoverMiddleware(
				logger, middleware.AccessHitsMiddleware(
					middleware.Authentication(sessionManager,
						middleware.CORSMiddleware(handler.DeleteSession(sessionManager)))))),
	)

	//Scoreboard routes
	router.GET("/api/v1/scores",
		middleware.LoggerMiddleware(
			logger, middleware.RecoverMiddleware(
				logger, middleware.AccessHitsMiddleware(
					middleware.CORSMiddleware(handler.GetScoreboard(profileManager))))),
	)

	//Oauth routes
	router.GET("/api/v1/oauth/redirect",
		middleware.LoggerMiddleware(
			logger, middleware.RecoverMiddleware(
				logger, middleware.AccessHitsMiddleware(
					middleware.OauthMiddleware(
						clientId, clienSecret, handler.PostProfileOauth(profileManager, sessionManager))))),
	)

	//Metrics routes
	router.GET("/metrics", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		promhttp.Handler().ServeHTTP(w, r)
	})

	return &Api{
		Router: router,
	}
}
