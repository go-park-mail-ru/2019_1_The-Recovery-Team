package handlers

import (
	"api/environment"
	"api/middleware"
	"api/models"
	"github.com/asaskevich/govalidator"
	"github.com/jackc/pgx"
	"go.uber.org/zap"
	"net/http"
	"time"
)

// GetSession returns handler with environment whick gets profile id of current session
// @Summary Get session
// @Description Get profile id of authorized client
// @ID get-session
// @Produce json
// @Success 200 {object} models.ProfileID "Profile ID found successfully"
// @Failure 403 "Not authorized"
// @Router /sessions [GET]
func GetSession(env *environment.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.Context().Value(middleware.ProfileID)
		result := &models.ProfileID{
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
// @Success 200 {object} models.Profile "Session is created successfully"
// @Failure 400 "Incorrect request data"
// @Failure 422 {object} models.HandlerError "Invalid request data"
// @Failure 403 "Not authorized"
// @Router /sessions [POST]
func PostSession(env *environment.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := r.Context().Value("logger").(*zap.Logger)

		login := &models.ProfileLogin{}
		err := unmarshalJSONBodyToStruct(r, login)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if isValid, err := govalidator.ValidateStruct(login); !isValid && err != nil {
			message := models.HandlerError{
				Description: err.Error(),
			}
			writeResponseJSON(w, http.StatusUnprocessableEntity, message)
			return
		}

		profile, err := env.Dbm.GetProfileByEmailWithPassword(login)
		if err != nil {
			if err == pgx.ErrNoRows {
				message := models.HandlerError{
					Description: "Invalid email or password",
				}
				writeResponseJSON(w, http.StatusUnprocessableEntity, message)
				return
			}
			log.Error(err.Error(),
				zap.String("email", login.Email))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		token, err := env.Sm.Set(profile.ID, 24*time.Hour)
		if err != nil {
			log.Error(err.Error(),
				zap.Uint64("profile_id", profile.ID))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "session_id",
			Value:    token,
			Path:     "/",
			Expires:  time.Now().Add(24*time.Hour - 10*time.Minute),
			HttpOnly: true,
		})

		writeResponseJSON(w, http.StatusOK, profile)
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
func DeleteSession(env *environment.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionID := r.Context().Value(middleware.SessionID)
		err := env.Sm.Delete(sessionID.(string))
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		cookie := http.Cookie{
			Name:     "session_id",
			Path:     "/",
			MaxAge:   -1,
			HttpOnly: true,
		}
		http.SetCookie(w, &cookie)
		w.WriteHeader(http.StatusOK)
	}
}
