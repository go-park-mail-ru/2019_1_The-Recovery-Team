package main

import (
	"log"
	"net/http"
	"os"

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
		port = "9090"
	}

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal("Logger creation error")
	}
	defer logger.Sync()

	profileConn, err := grpc.Dial(
		"localhost:9091",
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatal("Profile service connection refused")
	}
	defer profileConn.Close()

	sessionConn, err := grpc.Dial(
		"localhost:8081",
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatal("Authentication service connection refused")
	}
	defer sessionConn.Close()

	sessionManager := session.NewSessionClient(sessionConn)
	profileManager := profile.NewProfileClient(profileConn)
	gameManager := usecase.NewGameInteractor(game.NewGameRepo(logger))

	go gameManager.Run()

	gameApi := api.NewApi(&profileManager, &sessionManager, gameManager, logger)

	log.Print(http.ListenAndServe(":"+port, gameApi.Router))
}
