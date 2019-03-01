package router

import (
	"net/http"

	"api/handlers"
	"api/models"
)

// Route contains data about route
type Route struct {
	Name    string
	Method  string
	Pattern string
	Handler func(*models.Env) http.HandlerFunc
}

// Routes contains all routes
type Routes []Route

var routes = Routes{
	{
		"GetProfileById",
		"GET",
		"/profiles/{id:[0-9]+}",
		handlers.GetProfile,
	},
}
