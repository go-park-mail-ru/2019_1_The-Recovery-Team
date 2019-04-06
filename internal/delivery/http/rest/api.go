package rest

import (
	"sadislands/internal/delivery/http/rest/handler"
	"sadislands/internal/delivery/http/rest/middleware"
	"sadislands/internal/usecase"

	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

type Api struct {
	Router *httprouter.Router
}

func NewRestApi(
	profileInteractor *usecase.ProfileInteractor,
	sessionInteractor *usecase.SessionInteractor,
	gameInteractor *usecase.GameInteractor,
	logger *zap.Logger,
) *Api {
	router := httprouter.New()

	//Profile routes
	router.GET("/api/v1/profiles",
		middleware.LoggerMiddleware(
			logger, middleware.RecoverMiddleware(
				logger, middleware.CORSMiddleware(handler.GetProfiles(profileInteractor)))),
	)
	router.POST("/api/v1/profiles",
		middleware.LoggerMiddleware(
			logger, middleware.RecoverMiddleware(
				logger, middleware.CORSMiddleware(handler.PostProfile(profileInteractor, sessionInteractor)))),
	)
	router.GET("/api/v1/profiles/:id",
		middleware.LoggerMiddleware(
			logger, middleware.RecoverMiddleware(
				logger, middleware.Authentication(
					sessionInteractor, middleware.CORSMiddleware(handler.GetProfile(profileInteractor))))),
	)
	router.PUT("/api/v1/profiles/:id",
		middleware.LoggerMiddleware(
			logger, middleware.RecoverMiddleware(
				logger, middleware.Authentication(
					sessionInteractor, middleware.CORSMiddleware(handler.PutProfile(profileInteractor))))),
	)
	router.PUT("/api/v1/profiles/:id/password",
		middleware.LoggerMiddleware(
			logger, middleware.RecoverMiddleware(
				logger, middleware.Authentication(
					sessionInteractor, middleware.CORSMiddleware(handler.PutProfilePassword(profileInteractor))))),
	)
	router.PUT("/api/v1/avatars",
		middleware.LoggerMiddleware(
			logger, middleware.RecoverMiddleware(
				logger, middleware.Authentication(
					sessionInteractor, middleware.CORSMiddleware(handler.PutAvatar(profileInteractor))))),
	)

	//Session routes
	router.GET("/api/v1/sessions",
		middleware.LoggerMiddleware(
			logger, middleware.RecoverMiddleware(
				logger, middleware.Authentication(
					sessionInteractor, middleware.CORSMiddleware(handler.GetSession())))),
	)
	router.POST("/api/v1/sessions",
		middleware.LoggerMiddleware(
			logger, middleware.RecoverMiddleware(
				logger, middleware.CORSMiddleware(handler.PostSession(profileInteractor, sessionInteractor)))),
	)
	router.DELETE("/api/v1/sessions",
		middleware.LoggerMiddleware(
			logger, middleware.RecoverMiddleware(
				logger, middleware.Authentication(
					sessionInteractor, middleware.CORSMiddleware(handler.DeleteSession(sessionInteractor))))),
	)

	//Scoreboard routes
	router.GET("/api/v1/scores",
		middleware.LoggerMiddleware(
			logger, middleware.RecoverMiddleware(
				logger, middleware.CORSMiddleware(handler.GetScoreboard(profileInteractor)))),
	)

	//Game routes
	router.GET("/api/v1/search",
		middleware.Authentication(
			sessionInteractor, handler.Search(profileInteractor, gameInteractor)),
	)

	return &Api{
		Router: router,
	}
}
