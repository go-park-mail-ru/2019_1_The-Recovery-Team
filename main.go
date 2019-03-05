package main

import (
	"net/http"

	"api/database"
	"api/router"
)

func main() {
	dbm, err := database.InitDatabaseManager("recoveryteam", "123456", "localhost", "sadislands")
	if err != nil {
		panic(err)
	}

	router := router.InitRouter(dbm)
	if err := http.ListenAndServe(":9090", router); err != nil {
		panic(err)
	}
}
