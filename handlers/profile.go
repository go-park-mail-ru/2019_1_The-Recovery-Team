package handlers

import (
	"net/http"

	"github.com/mailru/easyjson"

	"api/models"

	"github.com/gorilla/mux"
)

// GetProfile returns handler with environment which processes request for getting profile by id
func GetProfile(env *models.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, _ := vars["id"]

		profile := &models.Profile{}
		err := env.Dbm.Find(profile, QueryProfileById, id)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			message := models.HandlerError{
				Description: "user doesn't exist",
			}
			easyjson.MarshalToHTTPResponseWriter(message, w)
			return
		}

		w.WriteHeader(http.StatusOK)
		easyjson.MarshalToHTTPResponseWriter(profile, w)
	}
}
