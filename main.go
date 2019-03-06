package main

import (
	"api/filesystem"
	"api/middleware"
	"api/session"
	"net/http"

	"api/database"
	"api/models"
	"api/router"
)

func main() {
	dbm, err := database.InitDatabaseManager("recoveryteam", "123456", "localhost", "sadislands")
	if err != nil {
		panic(err)
	}

	sm, err := session.InitSessionManager("", "", "localhost")
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

	panic(http.ListenAndServe(":9090", router))
}
