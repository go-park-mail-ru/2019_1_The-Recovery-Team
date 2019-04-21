package api

import (
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/delivery/grpc/service/profile"
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/delivery/grpc/service/session"
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/delivery/http/rest/middleware"
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/delivery/http/rest/profile/api/handler"
	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

type Api struct {
	Router *httprouter.Router
}

func NewApi(
	profileManager *profile.ProfileClient,
	sessionManager *session.SessionClient,
	logger *zap.Logger,
) *Api {
	router := httprouter.New()

	//Profile routes
	router.GET("/api/v1/profiles",
		middleware.LoggerMiddleware(
			logger, middleware.RecoverMiddleware(
				logger, middleware.CORSMiddleware(handler.GetProfiles(profileManager)))),
	)
	router.POST("/api/v1/profiles",
		middleware.LoggerMiddleware(
			logger, middleware.RecoverMiddleware(
				logger, middleware.CORSMiddleware(handler.PostProfile(profileManager, sessionManager)))),
	)
	router.GET("/api/v1/profiles/:id",
		middleware.LoggerMiddleware(
			logger, middleware.RecoverMiddleware(
				logger, middleware.Authentication(
					sessionManager, middleware.CORSMiddleware(handler.GetProfile(profileManager))))),
	)
	router.PUT("/api/v1/profiles/:id",
		middleware.LoggerMiddleware(
			logger, middleware.RecoverMiddleware(
				logger, middleware.Authentication(
					sessionManager, middleware.CORSMiddleware(handler.PutProfile(profileManager))))),
	)
	router.PUT("/api/v1/profiles/:id/password",
		middleware.LoggerMiddleware(
			logger, middleware.RecoverMiddleware(
				logger, middleware.Authentication(
					sessionManager, middleware.CORSMiddleware(handler.PutProfilePassword(profileManager))))),
	)
	router.PUT("/api/v1/avatars",
		middleware.LoggerMiddleware(
			logger, middleware.RecoverMiddleware(
				logger, middleware.Authentication(
					sessionManager, middleware.CORSMiddleware(handler.PutAvatar(profileManager))))),
	)

	//Session routes
	router.GET("/api/v1/sessions",
		middleware.LoggerMiddleware(
			logger, middleware.RecoverMiddleware(
				logger, middleware.Authentication(
					sessionManager, middleware.CORSMiddleware(handler.GetSession())))),
	)
	router.POST("/api/v1/sessions",
		middleware.LoggerMiddleware(
			logger, middleware.RecoverMiddleware(
				logger, middleware.CORSMiddleware(handler.PostSession(profileManager, sessionManager)))),
	)
	router.DELETE("/api/v1/sessions",
		middleware.LoggerMiddleware(
			logger, middleware.RecoverMiddleware(
				logger, middleware.Authentication(
					sessionManager, middleware.CORSMiddleware(handler.DeleteSession(sessionManager))))),
	)

	//Scoreboard routes
	router.GET("/api/v1/scores",
		middleware.LoggerMiddleware(
			logger, middleware.RecoverMiddleware(
				logger, middleware.CORSMiddleware(handler.GetScoreboard(profileManager)))),
	)

	return &Api{
		Router: router,
	}
}
