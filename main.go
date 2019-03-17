package main

import (
	"api/database"
	_ "api/docs"
	"api/environment"
	"api/filesystem"
	"api/router"
	"api/session"
	"github.com/swaggo/http-swagger"
	"log"
	"net/http"
	"os"
	"time"
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

	dbm, err := database.InitDatabaseManager("migrations/0_initial.sql")
	if err != nil {
		panic(err)
	}
	defer dbm.Close()

	sm, err := session.InitSessionManager("", "", "redis:6379")
	if err != nil {
		panic(err)
	}
	defer sm.Close()

	env := &environment.Env{
		Dbm: dbm,
		Sm:  sm,
	}

	mainRouter := router.InitRouter(env)

	fs := filesystem.InitFileserverHandler("/upload/", "./upload")
	mainRouter.PathPrefix("/upload/").Handler(fs)

	mainRouter.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	log.Print(http.ListenAndServe(":"+port, mainRouter))
}
