package chat

import (
	"net/http"

	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/delivery/grpc/service/session"
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/delivery/http/rest/api/chat/handler"
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/delivery/http/rest/middleware"
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/usecase"
	"github.com/julienschmidt/httprouter"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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
				logger, middleware.CORSMiddleware(handler.Connect(chatManager, sessionManager, logger)))))

	//Metrics routes
	router.GET("/metrics", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		promhttp.Handler().ServeHTTP(w, r)
	})

	return &Api{
		Router: router,
	}
}
