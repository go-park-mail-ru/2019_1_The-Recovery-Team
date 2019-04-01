package main

import (
	"log"
	"net/http"
	"os"
	"sadislands/internal/delivery/http/rest"
	"sadislands/internal/infrastructure/repository/memory/game"
	"sadislands/internal/infrastructure/repository/postgresql"
	"sadislands/internal/infrastructure/repository/postgresql/profile"
	"sadislands/internal/infrastructure/repository/redis/session"
	"sadislands/internal/usecase"

	_ "sadislands/docs"

	"github.com/gomodule/redigo/redis"
	"github.com/jackc/pgx"
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

	psqlConfig := pgx.ConnConfig{
		Host:     "localhost",
		Port:     5432,
		Database: "sadislands",
		User:     "recoveryteam",
		Password: "123456",
	}

	psqlConn, err := pgx.Connect(psqlConfig)
	if err != nil {
		log.Fatal("Postgresql connection refused")
	}
	defer psqlConn.Close()

	if err := postgresql.MakeMigrations(psqlConn, "build/schema/0_initial.sql"); err != nil {
		log.Fatal("Database migrations failed", err)
	}

	redisConn, err := redis.DialURL("redis://:@localhost:6379")
	if err != nil {
		log.Fatal("Redis connection refused")
	}
	defer redisConn.Close()

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal("Logger creation error")
	}
	defer logger.Sync()

	profileInteractor := usecase.NewProfileInteractor(profile.NewProfileRepo(psqlConn))
	sessionInterctor := usecase.NewSessionInteractor(session.NewSessionRepo(&redisConn))
	gameInteractor := usecase.NewGameInteractor(game.NewGameRepo(logger))
	go gameInteractor.Run()

	api := rest.NewRestApi(profileInteractor, sessionInterctor, gameInteractor, logger)

	api.Router.Handler("GET", "/swagger/:file", httpSwagger.WrapHandler)

	api.Router.ServeFiles("/upload/*filepath", http.Dir("../../upload"))

	log.Print(http.ListenAndServe(":"+port, api.Router))
}
