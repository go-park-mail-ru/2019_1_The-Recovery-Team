package handlers

import (
	"api/middleware"
	"api/models"
	"database/sql"
	"net/http"
	"time"
)

// GetSession returns handler with environment whick gets profile id of current session
// @Summary Get session
// @Description Get profile id of authorized client
// @ID get-session
// @Produce json
// @Success 200 int models.Profile.ID "Profile ID found successfully"
// @Failure 403 "Not authorized"
// @Router /sessions [GET]
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
// @Summary Post session
// @Description Creates client session
// @ID post-session
// @Accept json
// @Param profile_login body models.ProfileLogin true "Email, password"
// @Success 200 {object} models.Profile Session is created successfully"
// @Failure 400 "Incorrect request data"
// @Failure 422 "Unprocessable request data"
// @Failure 403 "Not authorized"
// @Router /sessions [POST]
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

		writeResponseJSON(w, http.StatusOK, result)
	}
}

// DeleteSession returns handler with environment which deletes session
// @Summary Delete session
// @Description Deletes client session
// @ID delete-session
// @Success 200 "Session is deleted successfully"
// @Failure 404 "Session not found"
// @Failure 403 "Not authorized"
// @Router /sessions [DELETE]
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
