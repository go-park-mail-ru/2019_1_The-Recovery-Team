package main

import (
	"api/database"
	_ "api/docs"
	"api/filesystem"
	"api/middleware"
	"api/models"
	"api/router"
	"api/session"
	"net/http"
	"github.com/swaggo/http-swagger"
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

	router := router.InitRouter(env)

	fs := filesystem.InitFileserverHandler("/upload/", "./upload")
	router.PathPrefix("/upload/").Handler(middleware.Authentication(env, fs))

	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	panic(http.ListenAndServe(":8080", router))
}
