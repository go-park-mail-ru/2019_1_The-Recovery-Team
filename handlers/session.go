package handlers

import (
	"api/middleware"
	"api/models"
	"database/sql"
	"net/http"
	"time"
)

// GetSession returns handler with environment whick gets profile id of current session
func GetSession(env *models.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.Context().Value(middleware.ProfileID)
		result := &models.Profile{
			ID: id.(uint64),
		}
		writeResponseJSON(w, http.StatusOK, result)
		return
	}
}

// PostSession returns handler with environment which creates session
func PostSession(env *models.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		profile := &models.ProfileLogin{}
		err := unmarshalJSONBodyToStruct(r, profile)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		result := &models.Profile{}
		err = env.Dbm.Find(result, QueryProfileByEmailWithPassword, profile.Email)
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}

		if matches, err := verifyPassword(profile.Password, result.Password); !matches || err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}

		token, err := env.Sm.Set(result.ID, 24*time.Hour)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "session_id",
			Value:    token,
			Expires:  time.Now().Add(24*time.Hour - 10*time.Minute),
			HttpOnly: true,
		})
		w.WriteHeader(http.StatusOK)
	}
}

// DeleteSession returns handler with environment which deletes session
func DeleteSession(env *models.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionID := r.Context().Value(middleware.SessionID)
		err := env.Sm.Delete(sessionID.(string))
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		cookie, _ := r.Cookie("session_id")
		cookie.Expires = time.Now().AddDate(0, 0, -1)
		http.SetCookie(w, cookie)
		w.WriteHeader(http.StatusOK)
	}
}
