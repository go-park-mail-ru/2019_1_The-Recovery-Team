package game

import (
	"context"
	"net/http"

	"google.golang.org/grpc/status"

	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/delivery/grpc/service/profile"
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/delivery/grpc/service/session"
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/delivery/http/rest/api/game/handler"
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

	//Metrics routes
	router.GET("/metrics", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		promhttp.Handler().ServeHTTP(w, r)
	})

	in := gameInteractor.UpdateRatingChan()
	go func() {
		for update := range in {
			r := &profile.UpdateRatingRequest{
				Winner: update.Winner,
				Loser:  update.Loser,
			}
			_, err := (*profileManager).UpdateRating(context.Background(), r)
			if err != nil {
				message := status.Convert(err).Message()
				logger.Error(message,
					zap.Uint64("winner_id", update.Winner),
					zap.Uint64("loser_id", update.Loser))
			}
		}
	}()
	go gameInteractor.Run()

	return &Api{
		Router: router,
	}
}
