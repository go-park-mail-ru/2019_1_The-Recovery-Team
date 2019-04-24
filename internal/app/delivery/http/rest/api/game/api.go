package game

import (
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/delivery/grpc/service/profile"
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/delivery/grpc/service/session"
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/delivery/http/rest/api/game/handler"
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/delivery/http/rest/middleware"
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/usecase"
	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

type Api struct {
	Router *httprouter.Router
}

func NewApi(
	profileManager *profile.ProfileClient,
	sessionManager *session.SessionClient,
	gameInteractor *usecase.GameInteractor,
	logger *zap.Logger,
) *Api {
	router := httprouter.New()

	//Game routes
	router.GET("/api/v1/game.ws",
		middleware.LoggerMiddleware(
			logger, middleware.RecoverMiddleware(
				logger, middleware.Authentication(
					sessionManager, middleware.CORSMiddleware(handler.Search(profileManager, gameInteractor))))),
	)

	return &Api{
		Router: router,
	}
}
