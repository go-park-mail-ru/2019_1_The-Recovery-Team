package main

import (
	"api/database"
	_ "api/docs"
	"api/filesystem"
	"api/models"
	"api/router"
	"api/session"
	"net/http"
	"os"
	"time"

	"github.com/swaggo/http-swagger"
)

// @title Sad Islands API
// @version 1.0
// @description This is a super game.

// @host 104.248.28.45
// @BasePath /api/v1

func main() {
	// Crutch for docker-compose (((
	time.Sleep(15 * time.Second)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbm, err := database.InitDatabaseManager("recoveryteam", "123456", "db:5432", "sadislands", "migrations", false)
	if err != nil {
		panic(err)
	}
	defer dbm.Close()

	sm, err := session.InitSessionManager("", "", "redis:6379")
	if err != nil {
		panic(err)
	}
	defer sm.Close()

	env := &models.Env{
		Dbm: dbm,
		Sm:  sm,
	}

	mainRouter := router.InitRouter(env)

	fs := filesystem.InitFileserverHandler("/upload/", "./upload")
	mainRouter.PathPrefix("/upload/").Handler(fs)

	mainRouter.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	panic(http.ListenAndServe(":"+port, mainRouter))
}
