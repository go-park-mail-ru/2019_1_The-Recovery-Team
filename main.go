package main

import (
	"api/database"
	_ "api/docs"
	"api/filesystem"
	"api/middleware"
	"api/models"
	"api/router"
	"api/session"
	"github.com/swaggo/http-swagger"
	"net/http"
)

// @host 127.0.0.1:8080

func main() {
	dbm, err := database.InitDatabaseManager("recoveryteam", "123456", "db:5432", "sadislands")
	if err != nil {
		panic(err)
	}

	sm, err := session.InitSessionManager("", "", "redis:6379")
	if err != nil {
		panic(err)
	}

	env := &models.Env{
		Dbm: dbm,
		Sm:  sm,
	}

	mainRouter := router.InitRouter(env)

	fs := filesystem.InitFileserverHandler("/upload/", "./upload")
	mainRouter.PathPrefix("/upload/").Handler(middleware.Authentication(env, fs))

	mainRouter.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	panic(http.ListenAndServe(":8080", mainRouter))
}
