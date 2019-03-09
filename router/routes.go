package router

import (
	"api/middleware"
	"net/http"

	"api/handlers"
	"api/models"
)

// Route contains data about route
type Route struct {
	Name        string
	Method      string
	Pattern     string
	Middlewares []middleware.MiddlewareWithEnv
	HandlerFunc func(*models.Env) http.HandlerFunc
}

// Routes contains all routes
type Routes []Route

var routes = Routes{
	{
		"GetProfiles",
		"GET",
		"/profiles",
		[]middleware.MiddlewareWithEnv{middleware.Authentication},
		handlers.GetProfiles,
	},
	{
		"CreateProfile",
		"POST",
		"/profiles",
		[]middleware.MiddlewareWithEnv{},
		handlers.PostProfile,
	},
	{
		"GetProfileById",
		"GET",
		"/profiles/{id:[0-9]+}",
		[]middleware.MiddlewareWithEnv{middleware.Authentication},
		handlers.GetProfile,
	},
	{
		"UpdateProfile",
		"PUT",
		"/profiles/{id:[0-9]+}",
		[]middleware.MiddlewareWithEnv{middleware.Authentication},
		handlers.PutProfile,
	},
	{
		"UpdateProfileAvatar",
		"PUT",
		"/avatars",
		[]middleware.MiddlewareWithEnv{middleware.Authentication},
		handlers.PutAvatar,
	},
	{
		"GetSession",
		"GET",
		"/sessions",
		[]middleware.MiddlewareWithEnv{middleware.Authentication},
		handlers.GetSession,
	},
	{
		"CreateSession",
		"POST",
		"/sessions",
		[]middleware.MiddlewareWithEnv{},
		handlers.PostSession,
	},
	{
		"DeleteSession",
		"DELETE",
		"/sessions",
		[]middleware.MiddlewareWithEnv{middleware.Authentication},
		handlers.DeleteSession,
	},
	{
		"GetScoreboard",
		"GET",
		"/scores",
		[]middleware.MiddlewareWithEnv{},
		handlers.GetScoreboard,
	},
}
