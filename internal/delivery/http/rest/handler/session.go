package handler

import (
	"net/http"
	handler "sadislands/internal/delivery/http/rest/handler/error"
	"sadislands/internal/delivery/http/rest/handler/unmarshaler"
	"sadislands/internal/delivery/http/rest/handler/writer"
	"sadislands/internal/delivery/http/rest/middleware"
	"sadislands/internal/domain/profile"
	"sadislands/internal/usecase"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/jackc/pgx"
	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

// GetSession returns handler with environment whick gets profile id of current session
// @Summary Get session
// @Description Get profile id of authorized client
// @ID get-session
// @Produce json
// @Success 200 {object} profile.ID "Profile ID found successfully"
// @Failure 403 "Not authorized"
// @Router /sessions [GET]
func GetSession() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		id := r.Context().Value(middleware.ProfileID)
		result := &profile.ID{
			Id: id.(uint64),
		}
		writer.WriteResponseJSON(w, http.StatusOK, result)
		return
	}
}

// PostSession returns handler with environment which creates session
// @Summary Post session
// @Description Creates client session
// @ID post-session
// @Accept json
// @Param profile_login body profile.Login true "Email, password"
// @Success 200 {object} profile.Profile "Session is created successfully"
// @Failure 400 "Incorrect request data"
// @Failure 422 {object} handler.Error "Invalid request data"
// @Failure 403 "Not authorized"
// @Router /sessions [POST]
func PostSession(profileInteractor *usecase.ProfileInteractor, sessionInteractor *usecase.SessionInteractor) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		log := r.Context().Value("logger").(*zap.Logger)

		login := &profile.Login{}
		err := unmarshaler.UnmarshalJSONBodyToStruct(r, login)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if isValid, err := govalidator.ValidateStruct(login); !isValid && err != nil {
			message := handler.Error{
				Description: err.Error(),
			}
			writer.WriteResponseJSON(w, http.StatusUnprocessableEntity, message)
			return
		}

		// Get profile by email and password
		prof, err := profileInteractor.GetProfileByEmailWithPassword(login)
		if err != nil {
			if err == pgx.ErrNoRows {
				message := handler.Error{
					Description: "Invalid email or password",
				}
				writer.WriteResponseJSON(w, http.StatusUnprocessableEntity, message)
				return
			}
			log.Error(err.Error(),
				zap.String("email", login.Email))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Create session
		token, err := sessionInteractor.Set(prof.ID, 24*time.Hour)
		if err != nil {
			log.Error(err.Error(),
				zap.Uint64("profile_id", prof.ID))
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

		writer.WriteResponseJSON(w, http.StatusOK, prof)
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
func DeleteSession(sessionInteractor *usecase.SessionInteractor) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		sessionID := r.Context().Value(middleware.SessionID)
		err := sessionInteractor.Delete(sessionID.(string))
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
