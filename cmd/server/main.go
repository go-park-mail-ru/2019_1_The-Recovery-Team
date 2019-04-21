package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/delivery/http/rest/profile/api"

	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/delivery/grpc/service/profile"

	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/delivery/grpc/service/session"

	"google.golang.org/grpc"

	_ "github.com/go-park-mail-ru/2019_1_The-Recovery-Team/docs"
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/infrastructure/repository/memory/game"
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/usecase"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
)

// @title Sad Islands API
// @version 1.0
// @description This is a super game.

// @host 104.248.28.45
// @BasePath /api/v1

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
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

	gameInteractor := usecase.NewGameInteractor(game.NewGameRepo(logger))
	go gameInteractor.Run()

	profileApi := api.NewApi(&profileManager, &sessionManager, logger)
	profileApi.Router.Handler("GET", "/swagger/:file", httpSwagger.WrapHandler)
	profileApi.Router.ServeFiles("/upload/*filepath", http.Dir("upload"))

	log.Print(http.ListenAndServe(":"+port, profileApi.Router))
}
