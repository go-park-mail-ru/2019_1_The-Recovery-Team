package router

import (
	"api/environment"
	"github.com/gorilla/mux"
)

// InitRouter returns router with inintialized handlers
func InitRouter(env *environment.Env) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		handler := route.HandlerFunc(env)
		for _, middlewareWrapper := range route.Middlewares {
			handler = middlewareWrapper(env, handler)
		}
		router.
			PathPrefix("/api/v1/").
			Methods(route.Method, "OPTIONS").
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}
	return router
}
