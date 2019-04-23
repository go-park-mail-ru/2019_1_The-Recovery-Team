package main

import (
	"log"
	"net/http"
	"os"

	"google.golang.org/grpc/balancer/roundrobin"

	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/delivery/grpc/service/profile"
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/delivery/grpc/service/session"
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/delivery/http/rest/game/api"
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/infrastructure/repository/memory/game"
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/usecase"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal("Logger creation error")
	}
	defer logger.Sync()

	profileConn, err := grpc.Dial("srv://consul/profile-service", grpc.WithInsecure(), grpc.WithBalancerName(roundrobin.Name))
	if err != nil {
		log.Fatal("Can't connect to profile service:", err)
	}
	defer profileConn.Close()

	sessionConn, err := grpc.Dial("srv://consul/session-service", grpc.WithInsecure(), grpc.WithBalancerName(roundrobin.Name))
	if err != nil {
		log.Fatal("Authentication service connection refused:", err)
	}
	defer sessionConn.Close()

	sessionManager := session.NewSessionClient(sessionConn)
	profileManager := profile.NewProfileClient(profileConn)
	gameManager := usecase.NewGameInteractor(game.NewGameRepo(logger))

	go gameManager.Run()

	gameApi := api.NewApi(&profileManager, &sessionManager, gameManager, logger)

	log.Print(http.ListenAndServe(":"+port, gameApi.Router))
}
