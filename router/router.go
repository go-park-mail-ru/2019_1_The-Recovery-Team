package router

import (
	"api/models"

	"github.com/gorilla/mux"
)

// func setupHandler(next http.HandlerFunc) http.HandlerFunc {
// 	next = middleware.CORSMiddleware(next)
// 	next = middleware.AccessLogMiddleware(next)
// 	next = middleware.RecoverMiddleware(next)
// 	return next
// }

// InitRouter returns router with inintialized handlers
func InitRouter(env *models.Env) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		handler := route.HandlerFunc(env)
		// handler = setupHandler(handler)
		for _, middlewareWrapper := range route.Middlewares {
			handler = middlewareWrapper(env, handler)
		}
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}
	return router
}
