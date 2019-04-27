package chat

import (
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/delivery/grpc/service/session"
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/delivery/http/rest/api/chat/handler"
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/delivery/http/rest/middleware"
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/usecase"
	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

type Api struct {
	Router *httprouter.Router
}

func NewApi(
	chatManager *usecase.ChatInteractor,
	sessionManager *session.SessionClient,
	logger *zap.Logger,
) *Api {
	router := httprouter.New()

	//Chat routes
	router.GET("/api/v1/chat.ws",
		middleware.LoggerMiddleware(
			logger, middleware.RecoverMiddleware(
				logger, middleware.Authentication(
					sessionManager, middleware.CORSMiddleware(handler.Connect(chatManager))))))

	return &Api{
		Router: router,
	}
}
