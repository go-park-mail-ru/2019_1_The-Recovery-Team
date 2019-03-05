package main

import (
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

	env := &models.Env{
		Dbm: dbm,
	}

	router := router.InitRouter(env)
	panic(http.ListenAndServe(":9090", router))
}
