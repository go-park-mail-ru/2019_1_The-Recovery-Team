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
		"GetProfiles",
		"GET",
		"/profiles/",
		handlers.GetProfiles,
	},
	{
		"CreateProfile",
		"POST",
		"/profiles/",
		handlers.PostProfile,
	},
	{
		"GetProfileById",
		"GET",
		"/profiles/{id:[0-9]+}",
		handlers.GetProfile,
	},
	{
		"CheckProfileEmail",
		"GET",
		"/profiles/email/{email}",
		handlers.CheckProfileEmail,
	},
	{
		"CheckProfileNickname",
		"GET",
		"/profiles/nickname/{nickname}",
		handlers.CheckProfileNickname,
	},
	{
		"UpdateProfile",
		"PUT",
		"/profiles/{id:[0-9]+}",
		handlers.PutProfile,
	},
	{
		"UpdateProfileAvatar",
		"PUT",
		"/profiles/{id:[0-9]+}/avatar",
		handlers.PutAvatar,
	},
	{
		"CreateSession",
		"POST",
		"/sessions",
		handlers.PostSession,
	},
}
