package router

import (
	"api/environment"
	"api/middleware"
	"net/http"

	"api/handlers"
)

// Route contains data about route
type Route struct {
	Name        string
	Method      string
	Pattern     string
	Middlewares []middleware.MiddlewareWithEnv
	HandlerFunc func(*environment.Env) http.HandlerFunc
}

// Routes contains all routes
type Routes []Route

var routes = Routes{
	{
		"GetProfiles",
		"GET",
		"/profiles",
		[]middleware.MiddlewareWithEnv{middleware.CORSMiddleware, middleware.LoggerMiddleware, middleware.RecoverMiddleware},
		handlers.GetProfiles,
	},
	{
		"CreateProfile",
		"POST",
		"/profiles",
		[]middleware.MiddlewareWithEnv{middleware.CORSMiddleware, middleware.LoggerMiddleware, middleware.RecoverMiddleware},
		handlers.PostProfile,
	},
	{
		"GetProfileById",
		"GET",
		"/profiles/{id:[0-9]+}",
		[]middleware.MiddlewareWithEnv{middleware.Authentication, middleware.CORSMiddleware, middleware.LoggerMiddleware, middleware.RecoverMiddleware},
		handlers.GetProfile,
	},
	{
		"UpdateProfile",
		"PUT",
		"/profiles/{id:[0-9]+}",
		[]middleware.MiddlewareWithEnv{middleware.Authentication, middleware.CORSMiddleware, middleware.LoggerMiddleware, middleware.RecoverMiddleware},
		handlers.PutProfile,
	},
	{
		"UpdateProfilePassword",
		"PUT",
		"/profiles/{id:[0-9]+}/password",
		[]middleware.MiddlewareWithEnv{middleware.Authentication, middleware.CORSMiddleware, middleware.LoggerMiddleware, middleware.RecoverMiddleware},
		handlers.PutProfilePassword,
	},
	{
		"UpdateProfileAvatar",
		"PUT",
		"/avatars",
		[]middleware.MiddlewareWithEnv{middleware.Authentication, middleware.CORSMiddleware, middleware.LoggerMiddleware, middleware.RecoverMiddleware},
		handlers.PutAvatar,
	},
	{
		"GetSession",
		"GET",
		"/sessions",
		[]middleware.MiddlewareWithEnv{middleware.Authentication, middleware.CORSMiddleware, middleware.LoggerMiddleware, middleware.RecoverMiddleware},
		handlers.GetSession,
	},
	{
		"CreateSession",
		"POST",
		"/sessions",
		[]middleware.MiddlewareWithEnv{middleware.CORSMiddleware, middleware.LoggerMiddleware, middleware.RecoverMiddleware},
		handlers.PostSession,
	},
	{
		"DeleteSession",
		"DELETE",
		"/sessions",
		[]middleware.MiddlewareWithEnv{middleware.Authentication, middleware.CORSMiddleware, middleware.LoggerMiddleware, middleware.RecoverMiddleware},
		handlers.DeleteSession,
	},
	{
		"GetScoreboard",
		"GET",
		"/scores",
		[]middleware.MiddlewareWithEnv{middleware.CORSMiddleware, middleware.LoggerMiddleware, middleware.RecoverMiddleware},
		handlers.GetScoreboard,
	},
}
