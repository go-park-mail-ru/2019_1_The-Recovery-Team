package router

import (
	"api/models"

	"github.com/gorilla/mux"
)

// InitRouter returns router with inintialized handlers
func InitRouter(env *models.Env) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		handler := route.Handler
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler(env))
	}
	return router
}
