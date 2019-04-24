package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/delivery/grpc/service/profile"
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/delivery/grpc/service/session"
	gameApi "github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/delivery/http/rest/api/game"
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/infrastructure/repository/memory/game"
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/usecase"

	"google.golang.org/grpc/balancer/roundrobin"

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

	api := gameApi.NewApi(&profileManager, &sessionManager, gameManager, logger)

	log.Print(http.ListenAndServe(":"+port, api.Router))
}
